package classroom

import (
	"classroom-service/internal/assign"
	"classroom-service/internal/user"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomService interface {
	CreateClassroom(ctx context.Context, req *CreateClassroomRequest, userID string) (string, error)
	UpdateClassroom(ctx context.Context, req *UpdateClassroomRequest, id string) error
	GetClassroomsByUserID(ctx context.Context, userID string) ([]string, error)
}

type classroomService struct {
	ClassroomRepository ClassroomRepository
	AssignRepository    assign.AssignRepository
	UserService         user.UserService
}

func NewClassroomService(classroomRepository ClassroomRepository,
	assignRepository assign.AssignRepository,
	userService user.UserService) ClassroomService {
	return &classroomService{
		ClassroomRepository: classroomRepository,
		AssignRepository:    assignRepository,
		UserService:         userService,
	}
}

func (s *classroomService) CreateClassroom(ctx context.Context, req *CreateClassroomRequest, userID string) (string, error) {

	var locationID *primitive.ObjectID

	if req.LocationID != nil {
		obj, err := primitive.ObjectIDFromHex(*req.LocationID)
		if err != nil {
			return "", err
		}
		locationID = &obj
	} else {
		locationID = nil
	}

	var regionID *primitive.ObjectID

	if req.RegionID != nil {
		obj, err := primitive.ObjectIDFromHex(*req.RegionID)
		if err != nil {
			return "", err
		}
		regionID = &obj
	} else {
		regionID = nil
	}

	if req.Name == "" {
		return "", errors.New("name is required")
	}

	if userID == "" {
		return "", errors.New("user id is required")
	}

	ClassroomID := primitive.NewObjectID()

	data := &ClassRoom{
		ID:          ClassroomID,
		Name:        req.Name,
		Description: req.Description,
		Note:        req.Note,
		Icon:        req.Icon,
		LocationID:  locationID,
		RegionID:    regionID,
		CreatedBy:   userID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.ClassroomRepository.CreateClassroom(ctx, data)
	if err != nil {
		return "", err
	}

	return ClassroomID.Hex(), nil

}

func (s *classroomService) UpdateClassroom(ctx context.Context, req *UpdateClassroomRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid classroom id: %v", err)
	}

	classroom, err := s.ClassroomRepository.GetClassroomByID(ctx, objectID)
	if err != nil {
		return err
	}

	if classroom == nil {
		return fmt.Errorf("classroom not found")
	}

	if req.Name != "" {
		classroom.Name = req.Name
	}

	if req.Description != nil {
		classroom.Description = req.Description
	}

	if req.Icon != nil {
		classroom.Icon = req.Icon
	}

	if req.Note != nil {
		classroom.Note = req.Note
	}

	if req.RegionID != nil {
		regionObjID, err := primitive.ObjectIDFromHex(*req.RegionID)
		if err != nil {
			return fmt.Errorf("invalid region id: %v", err)
		}
		classroom.RegionID = &regionObjID
	}

	if req.LocationID != nil {
		locationObjID, err := primitive.ObjectIDFromHex(*req.LocationID)
		if err != nil {
			return fmt.Errorf("invalid location id: %v", err)
		}
		classroom.LocationID = &locationObjID
	}

	err = s.ClassroomRepository.UpdateClassroom(ctx, objectID, classroom)
	if err != nil {
		return err
	}

	return nil
}


func (s *classroomService) GetClassroomsByUserID(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}
