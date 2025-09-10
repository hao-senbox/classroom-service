package assign

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssignService interface {
	UpdateAssgin(ctx context.Context, request *UpdateAssginRequest, userID string, id string) error
}

type assignService struct {
	AssignRepository AssignRepository
}

func NewAssignService(repo AssignRepository) AssignService {
	return &assignService{
		AssignRepository: repo,
	}
}

func (s *assignService) UpdateAssgin(ctx context.Context, request *UpdateAssginRequest, userID string, id string) error {

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	assign, err := s.AssignRepository.FindAssign(ctx, obj)
	if err != nil {
		return err
	}

	if assign == nil {
		return errors.New("assign not found")
	}

	assignID := assign.ID

	if request.TeacherID != nil {
		if assign.StudentID != nil {
			exists, err := s.AssignRepository.FindDuplicate(ctx, assign.ClassRoomID, *assign.StudentID, *request.TeacherID)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("student already assigned to teacher")
			}
		}
		assign.TeacherID = request.TeacherID
	}

	if request.StudentID != nil {
		if assign.TeacherID != nil {
			exists, err := s.AssignRepository.FindDuplicate(ctx, assign.ClassRoomID, *request.StudentID, *assign.TeacherID)
			if err != nil {
				return err
			}
			if exists {
				return errors.New("teacher already assigned to student")
			}
		}
		assign.StudentID = request.StudentID
	}

	now := time.Now()
	assign.UpdatedAt = now

	return s.AssignRepository.UpdateAssgin(ctx, assignID, assign)
}
