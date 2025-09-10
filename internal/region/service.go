package region

import (
	"classroom-service/internal/assign"
	"classroom-service/internal/classroom"
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegionService interface {
	CreateRegion(ctx context.Context, req *CreateRegionRequest, userID string) (string, error)
	GetAllRegions(ctx context.Context) ([]*RegionResponse, error)
	GetRegion(ctx context.Context, id string) (*RegionResponse, error)
	UpdateRegion(ctx context.Context, id string, req *UpdateRegionRequest) error
	DeleteRegion(ctx context.Context, id string) error
}

type regionService struct {
	RegionRepository    RegionRepository
	ClassroomRepository classroom.ClassroomRepository
	AssignRepository    assign.AssignRepository
	UserService         user.UserService
	RoomService         room.RoomService
}

func NewRegionService(regionRepository RegionRepository,
	classroomRepository classroom.ClassroomRepository,
	assignRepository assign.AssignRepository,
	userService user.UserService,
	roomService room.RoomService) RegionService {
	return &regionService{
		RegionRepository:    regionRepository,
		ClassroomRepository: classroomRepository,
		AssignRepository:    assignRepository,
		UserService:         userService,
		RoomService:         roomService,
	}
}

func (r *regionService) CreateRegion(ctx context.Context, req *CreateRegionRequest, userID string) (string, error) {

	if req.Name == "" {
		return "", errors.New("name is required")
	}

	if userID == "" {
		return "", errors.New("user id is required")
	}

	ID := primitive.NewObjectID()

	data := &Region{
		ID:        ID,
		Name:      req.Name,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.RegionRepository.CreateRegion(ctx, data)
	if err != nil {
		return "", err
	}

	return ID.Hex(), nil

}

func (r *regionService) GetAllRegions(ctx context.Context) ([]*RegionResponse, error) {

	regions, err := r.RegionRepository.GetRegions(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*RegionResponse

	for _, region := range regions {

		classrooms, err := r.ClassroomRepository.GetClassroomByRegion(ctx, region.ID)
		if err != nil {
			return nil, err
		}

		var classroomResponses []ClassRoom

		for _, classroom := range classrooms {

			var roomInfor room.RoomInfor
			if classroom.LocationID != nil {
				roomData, err := r.RoomService.GetRoomByID(ctx, classroom.LocationID.Hex())
				if err == nil && roomData != nil {
					roomInfor = *roomData
				} else {
					roomInfor = room.RoomInfor{
						ID:   classroom.LocationID.Hex(),
						Name: "Deleted",
					}
				}
			}

			assignments, err := r.AssignRepository.GetAssignmentsByClassroomID(ctx, classroom.ID)
			if err != nil {
				return nil, err
			}

			var assignmentResponses []*TeacherStudentAssignment
			for _, a := range assignments {

				var studentInfo user.UserInfor
				if a.StudentID != nil && *a.StudentID != "" {
					stu, err := r.UserService.GetStudentInfor(ctx, *a.StudentID)
					if err == nil && stu != nil {
						studentInfo = *stu
					} else {
						studentInfo = user.UserInfor{
							UserID:   *a.StudentID,
							UserName: "Deleted",
						}
					}
				}

				var teacherInfo user.UserInfor
				if a.TeacherID != nil && *a.TeacherID != "" {
					tea, err := r.UserService.GetTeacherInfor(ctx, *a.TeacherID)
					if err == nil && tea != nil {
						teacherInfo = *tea
					} else {
						teacherInfo = user.UserInfor{
							UserID:   *a.TeacherID,
							UserName: "Deleted",
						}
					}
				}

				assignmentResp := &TeacherStudentAssignment{
					ID:             a.ID,
					ClassRoomID:    a.ClassRoomID,
					Teacher:        teacherInfo,
					Student:        studentInfo,
					CreatedBy:      a.CreatedBy,
					IsNotification: a.IsNotification,
					CreatedAt:      a.CreatedAt,
					UpdatedAt:      a.UpdatedAt,
				}
				assignmentResponses = append(assignmentResponses, assignmentResp)
			}

			classroomResp := ClassRoom{
				ID:          classroom.ID,
				RegionID:    classroom.RegionID,
				Name:        classroom.Name,
				Description: classroom.Description,
				Icon:        classroom.Icon,
				Note:        classroom.Note,
				Assignments: assignmentResponses,
				Room:        roomInfor,
				IsActive:    classroom.IsActive,
				CreatedBy:   classroom.CreatedBy,
				CreatedAt:   classroom.CreatedAt,
				UpdatedAt:   classroom.UpdatedAt,
			}

			classroomResponses = append(classroomResponses, classroomResp)
		}

		res := &RegionResponse{
			ID:         region.ID,
			Name:       region.Name,
			Classrooms: classroomResponses,
			CreatedBy:  region.CreatedBy,
			CreatedAt:  region.CreatedAt,
			UpdatedAt:  region.UpdatedAt,
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func (r *regionService) GetRegion(ctx context.Context, id string) (*RegionResponse, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	region, err := r.RegionRepository.GetRegion(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if region == nil {
		return nil, errors.New("region not found")
	}

	classrooms, err := r.ClassroomRepository.GetClassroomByRegion(ctx, region.ID)
	if err != nil {
		return nil, err
	}

	var classroomResponses []ClassRoom

	for _, classroom := range classrooms {

		var roomInfor room.RoomInfor
		if classroom.LocationID != nil {
			roomData, err := r.RoomService.GetRoomByID(ctx, classroom.LocationID.Hex())
			if err == nil && roomData != nil {
				roomInfor = *roomData
			} else {
				roomInfor = room.RoomInfor{
					ID:   classroom.LocationID.Hex(),
					Name: "Deleted",
				}
			}
		}

		assignments, err := r.AssignRepository.GetAssignmentsByClassroomID(ctx, classroom.ID)
		if err != nil {
			return nil, err
		}

		var assignmentResponses []*TeacherStudentAssignment
		for _, a := range assignments {

			var studentInfo user.UserInfor
			if a.StudentID != nil && *a.StudentID != "" {
				stu, err := r.UserService.GetStudentInfor(ctx, *a.StudentID)
				if err == nil && stu != nil {
					studentInfo = *stu
				} else {
					studentInfo = user.UserInfor{
						UserID:   *a.StudentID,
						UserName: "Deleted",
					}
				}
			}

			var teacherInfo user.UserInfor
			if a.TeacherID != nil && *a.TeacherID != "" {
				tea, err := r.UserService.GetTeacherInfor(ctx, *a.TeacherID)
				if err == nil && tea != nil {
					teacherInfo = *tea
				} else {
					teacherInfo = user.UserInfor{
						UserID:   *a.TeacherID,
						UserName: "Deleted",
					}
				}
			}

			assignmentResp := &TeacherStudentAssignment{
				ID:             a.ID,
				ClassRoomID:    a.ClassRoomID,
				Teacher:        teacherInfo,
				Student:        studentInfo,
				CreatedBy:      a.CreatedBy,
				IsNotification: a.IsNotification,
				CreatedAt:      a.CreatedAt,
				UpdatedAt:      a.UpdatedAt,
			}
			assignmentResponses = append(assignmentResponses, assignmentResp)
		}

		classroomResp := ClassRoom{
			ID:          classroom.ID,
			RegionID:    classroom.RegionID,
			Name:        classroom.Name,
			Description: classroom.Description,
			Icon:        classroom.Icon,
			Note:        classroom.Note,
			Assignments: assignmentResponses,
			Room:        roomInfor,
			IsActive:    classroom.IsActive,
			CreatedBy:   classroom.CreatedBy,
			CreatedAt:   classroom.CreatedAt,
			UpdatedAt:   classroom.UpdatedAt,
		}

		classroomResponses = append(classroomResponses, classroomResp)
	}

	res := &RegionResponse{
		ID:         region.ID,
		Name:       region.Name,
		Classrooms: classroomResponses,
		CreatedBy:  region.CreatedBy,
		CreatedAt:  region.CreatedAt,
		UpdatedAt:  region.UpdatedAt,
	}

	return res, nil

}

func (r *regionService) UpdateRegion(ctx context.Context, id string, req *UpdateRegionRequest) error {

	if id == "" {
		return errors.New("id is required")
	}

	if req.Name == "" {
		return errors.New("name is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	region, err := r.RegionRepository.GetRegion(ctx, objectID)
	if err != nil {
		return err
	}

	if region == nil {
		return errors.New("region not found")
	}

	region.Name = req.Name
	region.UpdatedAt = time.Now()

	return r.RegionRepository.UpdateRegion(ctx, objectID, region)

}

func (r *regionService) DeleteRegion(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.RegionRepository.DeleteRegion(ctx, objectID)

}
