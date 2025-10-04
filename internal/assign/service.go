package assign

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssignService interface {
	AssignSlot(ctx context.Context, request *UpdateAssginRequest, userID string) error
	UnAssignSlot(ctx context.Context, request *UpdateAssginRequest, userID string) error
	CreateAssignmentTemplate(ctx context.Context, request *UpdateAssginRequest, userID string) error
	DeleteAssignmentTemplate(ctx context.Context, request *UpdateAssginRequest, userID string) error
}

type assignService struct {
	AssignRepository AssignRepository
}

func NewAssignService(repo AssignRepository) AssignService {
	return &assignService{
		AssignRepository: repo,
	}
}

func (s *assignService) AssignSlot(ctx context.Context, request *UpdateAssginRequest, userID string) error {

	if request.SlotNumber < -1 || request.SlotNumber > 15 {
		return errors.New("slot number must be between 1 and 15")
	}

	classroomObjID, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	dateParse, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}

	existingAssignment, err := s.AssignRepository.GetAssignmentBySlotAndDate(ctx, classroomObjID, request.SlotNumber, &dateParse)
	if err != nil {
		return err
	}

	if existingAssignment == nil {
		newAssignment := &TeacherStudentAssignment{
			ID:             primitive.NewObjectID(),
			ClassRoomID:    classroomObjID,
			SlotNumber:     request.SlotNumber,
			AssignDate:     dateParse,
			TeacherID:      request.TeacherID,
			StudentID:      request.StudentID,
			CreatedBy:      userID,
			IsNotification: false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		return s.AssignRepository.CreateAssignment(ctx, newAssignment)
	} else {
		if request.TeacherID != nil {
			if existingAssignment.StudentID != nil {
				exists, err := s.AssignRepository.CheckDuplicateAssignmentForDate(ctx, classroomObjID, dateParse, *existingAssignment.StudentID, *request.TeacherID)
				if err != nil {
					return err
				}
				if exists {
					return errors.New("student already assigned to teacher")
				}
			}
			existingAssignment.TeacherID = request.TeacherID
		}
		if request.StudentID != nil {
			if existingAssignment.TeacherID != nil {
				exists, err := s.AssignRepository.CheckDuplicateAssignmentForDate(ctx, classroomObjID, dateParse, *request.StudentID, *existingAssignment.TeacherID)
				if err != nil {
					return err
				}
				if exists {
					return errors.New("teacher already assigned to student")
				}
			}
			existingAssignment.StudentID = request.StudentID
		}

		assign := &TeacherStudentAssignment{
			ID:             existingAssignment.ID,
			ClassRoomID:    existingAssignment.ClassRoomID,
			SlotNumber:     existingAssignment.SlotNumber,
			AssignDate:     existingAssignment.AssignDate,
			TeacherID:      existingAssignment.TeacherID,
			StudentID:      existingAssignment.StudentID,
			CreatedBy:      existingAssignment.CreatedBy,
			IsNotification: existingAssignment.IsNotification,
			CreatedAt:      existingAssignment.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		return s.AssignRepository.UpdateAssgin(ctx, assign.ID, assign)
	}
}

func (s *assignService) UnAssignSlot(ctx context.Context, request *UpdateAssginRequest, userID string) error {

	if request.SlotNumber < -1 || request.SlotNumber > 15 {
		return errors.New("slot number must be between 1 and 15")
	}

	classroomObjID, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	if request.Date == "" {
		return errors.New("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return err
	}

	assign, err := s.AssignRepository.GetAssignmentBySlotAndDate(ctx, classroomObjID, request.SlotNumber, &dateParse)
	if err != nil {
		return err
	}
	if assign == nil {
		return errors.New("assign not found")
	}

	if request.TeacherID != nil {
		assign.TeacherID = nil
	}
	if request.StudentID != nil {
		assign.StudentID = nil
	}

	return s.AssignRepository.UpdateAssgin(ctx, assign.ID, assign)
}

func (s *assignService) CreateAssignmentTemplate(ctx context.Context, request *UpdateAssginRequest, userID string) error {

	if request.SlotNumber < -1 || request.SlotNumber > 15 {
		return errors.New("slot number must be between 1 and 15")
	}

	classroomObjID, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	existingAssignment, err := s.AssignRepository.GetAssignmentTemplateBySlot(ctx, classroomObjID, request.SlotNumber)
	if err != nil {
		return err
	}

	if existingAssignment == nil {
		newAssignment := &ClassRoomTemplateAssignment{
			ID:          primitive.NewObjectID(),
			ClassRoomID: classroomObjID,
			SlotNumber:  request.SlotNumber,
			TeacherID:   request.TeacherID,
			StudentID:   request.StudentID,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		return s.AssignRepository.CreateAssignmentTemplate(ctx, newAssignment)
	} else {
		if request.TeacherID != nil {
			if existingAssignment.StudentID != nil {
				exists, err := s.AssignRepository.CheckDuplicateAssignmentTemplate(
					ctx,
					classroomObjID,
					*existingAssignment.StudentID,
					*request.TeacherID,
				)
				if err != nil {
					return err
				}
				if exists {
					return errors.New("teacher already assigned to this student")
				}
			}
			existingAssignment.TeacherID = request.TeacherID
		}

		if request.StudentID != nil {
			if existingAssignment.TeacherID != nil {
				exists, err := s.AssignRepository.CheckDuplicateAssignmentTemplate(
					ctx,
					classroomObjID,
					*request.StudentID,
					*existingAssignment.TeacherID,
				)
				if err != nil {
					return err
				}
				if exists {
					return errors.New("student already assigned to this teacher")
				}
			}
			existingAssignment.StudentID = request.StudentID
		}

		return s.AssignRepository.UpdateAssginTemplate(ctx, existingAssignment.ID, existingAssignment)
	}

}

func (s *assignService) DeleteAssignmentTemplate(ctx context.Context, request *UpdateAssginRequest, userID string) error {

	if request.SlotNumber < -1 || request.SlotNumber > 15 {
		return errors.New("slot number must be between 1 and 15")
	}

	classroomObjID, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	if request.Date == "" {
		return errors.New("date is required")
	}

	assign, err := s.AssignRepository.GetAssignmentTemplateBySlot(ctx, classroomObjID, request.SlotNumber)
	if err != nil {
		return err
	}

	if assign == nil {
		return errors.New("assign not found")
	}

	if request.TeacherID != nil {
		assign.TeacherID = nil
	}
	if request.StudentID != nil {
		assign.StudentID = nil
	}

	return s.AssignRepository.UpdateAssginTemplate(ctx, assign.ID, assign)

}
