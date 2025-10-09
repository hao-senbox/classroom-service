package region

import (
	"classroom-service/internal/assign"
	"classroom-service/internal/classroom"
	"classroom-service/internal/language"
	"classroom-service/internal/leader"
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegionService interface {
	CreateRegion(ctx context.Context, req *CreateRegionRequest, userID string) (string, error)
	GetAllRegions(ctx context.Context, organizationID string, date string) ([]*RegionResponse, error)
	GetRegion(ctx context.Context, id string, date string) (*RegionResponse, error)
	UpdateRegion(ctx context.Context, id string, req *UpdateRegionRequest) error
	DeleteRegion(ctx context.Context, id string) error
}

type regionService struct {
	RegionRepository    RegionRepository
	ClassroomRepository classroom.ClassroomRepository
	AssignRepository    assign.AssignRepository
	UserService         user.UserService
	RoomService         room.RoomService
	LeaderRepository    leader.LeaderRepository
	LanguageService     language.MessageLanguageGateway
}

func NewRegionService(regionRepository RegionRepository,
	classroomRepository classroom.ClassroomRepository,
	assignRepository assign.AssignRepository,
	userService user.UserService,
	roomService room.RoomService,
	leaderRepository leader.LeaderRepository,
	languageService language.MessageLanguageGateway) RegionService {
	return &regionService{
		RegionRepository:    regionRepository,
		ClassroomRepository: classroomRepository,
		AssignRepository:    assignRepository,
		UserService:         userService,
		RoomService:         roomService,
		LeaderRepository:    leaderRepository,
		LanguageService:     languageService,
	}
}

func (r *regionService) CreateRegion(ctx context.Context, req *CreateRegionRequest, userID string) (string, error) {

	if req.Name == "" {
		return "", errors.New("name is required")
	}

	if userID == "" {
		return "", errors.New("user id is required")
	}

	if req.OrganizationID == "" {
		return "", errors.New("organization id is required")
	}

	ID := primitive.NewObjectID()

	data := &Region{
		ID:             ID,
		Name:           req.Name,
		OrganizationID: req.OrganizationID,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := r.RegionRepository.CreateRegion(ctx, data)
	if err != nil {
		return "", err
	}

	return ID.Hex(), nil

}

func (r *regionService) GetAllRegions(ctx context.Context, organizationID string, date string) ([]*RegionResponse, error) {

	if date == "" {
		return nil, errors.New("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	regions, err := r.RegionRepository.GetRegions(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	var responses []*RegionResponse

	for _, region := range regions {
		classrooms, err := r.ClassroomRepository.GetClassroomByRegion(ctx, region.ID)
		if err != nil {
			return nil, err
		}

		classroomResponses := make([]*ClassRoomResponse, 0)
		for _, classroom := range classrooms {
			var roomInfor *room.RoomInfor
			if classroom.LocationID != nil {
				roomData, err := r.RoomService.GetRoomByID(ctx, classroom.LocationID.Hex())
				if err == nil && roomData != nil {
					roomInfor = roomData
				} else {
					roomInfor = &room.RoomInfor{
						ID:   classroom.LocationID.Hex(),
						Name: "Deleted",
					}
				}
			}

			leader, err := r.LeaderRepository.GetLeaderByClassIDAndDate(ctx, classroom.ID, &dateParse)
			if err != nil {
				return nil, err
			}

			var leaderInfor *user.UserInfor
			if leader != nil {
				if leader.Owner.OwnerRole == "teacher" {
					leaderInforData, err := r.UserService.GetTeacherInfor(ctx, leader.Owner.OwnerID)
					if err != nil {
						return nil, err
					}
					leaderInfor = &user.UserInfor{
						UserID:   leaderInforData.UserID,
						UserName: leaderInforData.UserName,
						Avartar:  leaderInforData.Avartar,
					}
				} else if leader.Owner.OwnerRole == "staff" {
					leaderInforData, err := r.UserService.GetStaffInfor(ctx, leader.Owner.OwnerID)
					if err != nil {
						return nil, err
					}
					leaderInfor = &user.UserInfor{
						UserID:   leaderInforData.UserID,
						UserName: leaderInforData.UserName,
						Avartar:  leaderInforData.Avartar,
					}
				}
			}

			allAssignments, err := r.AssignRepository.GetAssignmentsByClassroomAndDate(ctx, classroom.ID, &dateParse)
			if err != nil {
				return nil, err
			}

			var messageLanguage = make([]language.MessageLanguageResponse, 0)
			messageLanguageData, _ := r.LanguageService.GetMessageLanguages(ctx, classroom.ID.Hex())

			if messageLanguageData != nil {
				messageLanguage = messageLanguageData
			}
			
			assignmentResponses := make([]*SlotAssignmentResponse, 0)
			for _, assignment := range allAssignments {
				assignmentID := assignment.ID.Hex()

				assignmentResp := &SlotAssignmentResponse{
					SlotNumber:     assignment.SlotNumber,
					AssignmentID:   &assignmentID,
					AssignmentDate: &assignment.AssignDate,
					IsAssigned:     true,
					CreatedAt:      &assignment.CreatedAt,
					UpdatedAt:      &assignment.UpdatedAt,
				}

				if assignment.TeacherID != nil && *assignment.TeacherID != "" {
					teacherInfo, err := r.UserService.GetTeacherInfor(ctx, *assignment.TeacherID)
					if err == nil && teacherInfo != nil {
						assignmentResp.Teacher = teacherInfo
					} else {
						assignmentResp.Teacher = &user.UserInfor{
							UserID:   *assignment.TeacherID,
							UserName: "Deleted",
						}
					}
				}

				if assignment.StudentID != nil && *assignment.StudentID != "" {
					studentInfo, err := r.UserService.GetStudentInfor(ctx, *assignment.StudentID)
					if err == nil && studentInfo != nil {
						assignmentResp.Student = studentInfo
					} else {
						assignmentResp.Student = &user.UserInfor{
							UserID:   *assignment.StudentID,
							UserName: "Deleted",
						}
					}
				}

				assignmentResponses = append(assignmentResponses, assignmentResp)
			}

			classroomResp := &ClassRoomResponse{
				ID:                classroom.ID,
				RegionID:          classroom.RegionID,
				Name:              classroom.Name,
				Description:       getStringValue(classroom.Description),
				Icon:              getStringValue(classroom.Icon),
				Note:              getStringValue(classroom.Note),
				Room:              roomInfor,
				Leader:            leaderInfor,
				IsActive:          classroom.IsActive,
				CreatedBy:         classroom.CreatedBy,
				CreatedAt:         classroom.CreatedAt,
				UpdatedAt:         classroom.UpdatedAt,
				MessageLanguages:  messageLanguage,
				TotalSlots:        15,
				AssignedSlots:     len(assignmentResponses),
				AvailableSlots:    15 - len(assignmentResponses),
				RecentAssignments: assignmentResponses,
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

func (r *regionService) GetRegion(ctx context.Context, id string, date string) (*RegionResponse, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}
	if date == "" {
		return nil, errors.New("date is required")
	}

	dateParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
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

	classroomResponses := make([]*ClassRoomResponse, 0)

	for _, classroom := range classrooms {

		var roomInfor *room.RoomInfor
		if classroom.LocationID != nil {
			roomData, err := r.RoomService.GetRoomByID(ctx, classroom.LocationID.Hex())
			if err == nil && roomData != nil {
				roomInfor = roomData
			} else {
				roomInfor = &room.RoomInfor{
					ID:   classroom.LocationID.Hex(),
					Name: "Deleted",
				}
			}
		}

		leader, err := r.LeaderRepository.GetLeaderByClassIDAndDate(ctx, classroom.ID, &dateParse)
		if err != nil {
			return nil, err
		}

		var leaderInfor *user.UserInfor

		if leader != nil {
			if leader.Owner.OwnerRole == "teacher" {
				leaderInforData, err := r.UserService.GetTeacherInfor(ctx, leader.Owner.OwnerID)
				if err != nil {
					return nil, err
				}
				leaderInfor = &user.UserInfor{
					UserID:   leaderInforData.UserID,
					UserName: leaderInforData.UserName,
					Avartar:  leaderInforData.Avartar,
				}
			} else if leader.Owner.OwnerRole == "staff" {
				leaderInforData, err := r.UserService.GetStaffInfor(ctx, leader.Owner.OwnerID)
				if err != nil {
					return nil, err
				}
				leaderInfor = &user.UserInfor{
					UserID:   leaderInforData.UserID,
					UserName: leaderInforData.UserName,
					Avartar:  leaderInforData.Avartar,
				}
			}
		}

		allAssignments, err := r.AssignRepository.GetAssignmentsByClassroomAndDate(ctx, classroom.ID, &dateParse)
		if err != nil {
			return nil, err
		}

		assignmentResponses := make([]*SlotAssignmentResponse, 0)
		for _, assignment := range allAssignments {
			assignmentID := assignment.ID.Hex()

			assignmentResp := &SlotAssignmentResponse{
				SlotNumber:     assignment.SlotNumber,
				AssignmentID:   &assignmentID,
				AssignmentDate: &assignment.AssignDate,
				IsAssigned:     true,
				CreatedAt:      &assignment.CreatedAt,
				UpdatedAt:      &assignment.UpdatedAt,
			}

			if assignment.TeacherID != nil && *assignment.TeacherID != "" {
				teacherInfo, err := r.UserService.GetTeacherInfor(ctx, *assignment.TeacherID)
				if err == nil && teacherInfo != nil {
					assignmentResp.Teacher = teacherInfo
				} else {
					assignmentResp.Teacher = &user.UserInfor{
						UserID:   *assignment.TeacherID,
						UserName: "Deleted",
					}
				}
			}

			if assignment.StudentID != nil && *assignment.StudentID != "" {
				studentInfo, err := r.UserService.GetStudentInfor(ctx, *assignment.StudentID)
				if err == nil && studentInfo != nil {
					assignmentResp.Student = studentInfo
				} else {
					assignmentResp.Student = &user.UserInfor{
						UserID:   *assignment.StudentID,
						UserName: "Deleted",
					}
				}
			}

			assignmentResponses = append(assignmentResponses, assignmentResp)
		}

		classroomResp := &ClassRoomResponse{
			ID:                classroom.ID,
			RegionID:          classroom.RegionID,
			Name:              classroom.Name,
			Description:       getStringValue(classroom.Description),
			Icon:              getStringValue(classroom.Icon),
			Note:              getStringValue(classroom.Note),
			Room:              roomInfor,
			Leader:            leaderInfor,
			IsActive:          classroom.IsActive,
			CreatedBy:         classroom.CreatedBy,
			CreatedAt:         classroom.CreatedAt,
			UpdatedAt:         classroom.UpdatedAt,
			TotalSlots:        15,
			AssignedSlots:     len(assignmentResponses),
			AvailableSlots:    15 - len(assignmentResponses),
			RecentAssignments: assignmentResponses,
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

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
