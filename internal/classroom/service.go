package classroom

import (
	"classroom-service/internal/assign"
	"classroom-service/internal/language"
	"classroom-service/internal/leader"
	"classroom-service/internal/term"
	"classroom-service/internal/user"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomService interface {
	CreateClassroom(ctx context.Context, req *CreateClassroomRequest, userID string) (string, error)
	UpdateClassroom(ctx context.Context, req *UpdateClassroomRequest, id string) error
	GetClassroomsByUserID(ctx context.Context, userID string) ([]string, error)
	GetClassroomByID(ctx context.Context, id, start, end string) (*ClassroomScheduleResponse, error)
	//Classroom Template
	GetClassroomByIDTemplate(ctx context.Context, id string) (*ClassroomTemplateResponse, error)
	CreateAssignmentByTemplate(ctx context.Context, req *CreateAssignmentByTemplateRequest) error

	//Assignment
	GetTeacherAssignments(ctx context.Context, userID, organizationID string, termID string) ([]TeacherAssignmentResponse, error)
}

type classroomService struct {
	ClassroomRepository ClassroomRepository
	AssignRepository    assign.AssignRepository
	UserService         user.UserService
	LeaderRopitory      leader.LeaderRepository
	LanguageService     language.MessageLanguageGateway
	TermService         term.TermService
}

func NewClassroomService(classroomRepository ClassroomRepository,
	assignRepository assign.AssignRepository,
	userService user.UserService,
	leaderRepository leader.LeaderRepository,
	languageService language.MessageLanguageGateway,
	termService term.TermService) ClassroomService {
	return &classroomService{
		ClassroomRepository: classroomRepository,
		AssignRepository:    assignRepository,
		UserService:         userService,
		LeaderRopitory:      leaderRepository,
		LanguageService:     languageService,
		TermService:         termService,
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

	languageReq := BuildDepartmentMessagesUpdate(ClassroomID.Hex(), *req)

	err = s.LanguageService.UploadMessages(ctx, languageReq)
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

	if req.Icon != nil {
		classroom.Icon = req.Icon
	}

	note := ""
	if req.Note != nil {
		note = *req.Note
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
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

	reqLanguage := &CreateClassroomRequest{
		Name:        req.Name,
		LanguageID:  req.LanguageID,
		RegionID:    req.RegionID,
		LocationID:  req.LocationID,
		Description: &desc,
		Note:        &note,
		Icon:        req.Icon,
	}

	languageReq := BuildDepartmentMessagesUpdate(classroom.ID.Hex(), *reqLanguage)

	err = s.LanguageService.UploadMessages(ctx, languageReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *classroomService) GetClassroomsByUserID(ctx context.Context, userID string) ([]string, error) {
	return []string{}, nil
}

func (s *classroomService) GetClassroomByIDTemplate(ctx context.Context, id string) (*ClassroomTemplateResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	leader, err := s.LeaderRopitory.GetLeaderTemplateByClassID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	var leaderInfor *user.UserInfor

	if leader != nil && leader.Owner != nil {
		switch leader.Owner.OwnerRole {
		case "teacher":
			leaderInforData, err := s.UserService.GetTeacherInfor(ctx, leader.Owner.OwnerID)
			if err != nil {
				return nil, err
			}
			if leaderInforData != nil {
				leaderInfor = &user.UserInfor{
					UserID:   leaderInforData.UserID,
					UserName: leaderInforData.UserName,
					Avartar:  leaderInforData.Avartar,
				}
			} else {
				leaderInfor = &user.UserInfor{
					UserID:   "",
					UserName: "",
					Avartar:  user.Avatar{},
				}
			}

		case "staff":
			leaderInforData, err := s.UserService.GetStaffInfor(ctx, leader.Owner.OwnerID)
			if err != nil {
				return nil, err
			}
			if leaderInforData != nil {
				leaderInfor = &user.UserInfor{
					UserID:   leaderInforData.UserID,
					UserName: leaderInforData.UserName,
					Avartar:  leaderInforData.Avartar,
				}
			} else {
				leaderInfor = &user.UserInfor{
					UserID:   "",
					UserName: "",
					Avartar:  user.Avatar{},
				}
			}
		}
	} else {
		leaderInfor = &user.UserInfor{}
	}

	var assignTemplateResponse []*SlotAssignmentResponse

	for _, assignment := range assignTemplate {

		assignmentID := assignment.ID.Hex()

		assignmentResp := &SlotAssignmentResponse{
			SlotNumber:   assignment.SlotNumber,
			AssignmentID: &assignmentID,
			CreatedAt:    &assignment.CreatedAt,
			UpdatedAt:    &assignment.UpdatedAt,
		}

		if assignment.TeacherID != nil && *assignment.TeacherID != "" {
			teacherInfo, err := s.UserService.GetTeacherInfor(ctx, *assignment.TeacherID)
			if err == nil && teacherInfo != nil {
				assignmentResp.Teacher = teacherInfo
			} else {
				assignmentResp.Teacher = &user.UserInfor{
					UserID:   *assignment.TeacherID,
					UserName: "",
				}
			}
		}

		if assignment.StudentID != nil && *assignment.StudentID != "" {
			studentInfo, err := s.UserService.GetStudentInfor(ctx, *assignment.StudentID)
			if err == nil && studentInfo != nil {
				assignmentResp.Student = studentInfo
			} else {
				assignmentResp.Student = &user.UserInfor{
					UserID:   *assignment.StudentID,
					UserName: "",
				}
			}
		}

		assignTemplateResponse = append(assignTemplateResponse, assignmentResp)

	}

	return &ClassroomTemplateResponse{
		ClassroomID:    id,
		Leader:         leaderInfor,
		SlotAssignment: assignTemplateResponse,
	}, nil

}

func (s *classroomService) CreateAssignmentByTemplate(ctx context.Context, req *CreateAssignmentByTemplateRequest) error {

	if req.ClassroomID == "" {
		return errors.New("classroom id is required")
	}

	if req.StartDate == "" {
		return errors.New("start_date is required")
	}

	if req.EndDate == "" {
		return errors.New("end_date is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.ClassroomID)
	if err != nil {
		return err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID)
	if err != nil {
		return err
	}

	leaderTemplate, err := s.LeaderRopitory.GetLeaderTemplateByClassID(ctx, objectID)
	if err != nil {
		return err
	}

	startParse, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return err
	}

	endParse, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return err
	}

	if assignTemplate != nil && leaderTemplate != nil {
		for d := startParse; d.Before(endParse); d = d.AddDate(0, 0, 1) {
			leaderData := leader.Leader{
				ID:          primitive.NewObjectID(),
				Owner:       leaderTemplate.Owner,
				ClassRoomID: leaderTemplate.ClassRoomID,
				Date:        d,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := s.LeaderRopitory.CreateLeader(ctx, &leaderData)
			if err != nil {
				return err
			}
			for _, assignment := range assignTemplate {
				assignmentData := assign.TeacherStudentAssignment{
					ID:             primitive.NewObjectID(),
					ClassRoomID:    objectID,
					SlotNumber:     assignment.SlotNumber,
					AssignDate:     d,
					TeacherID:      assignment.TeacherID,
					StudentID:      assignment.StudentID,
					CreatedBy:      assignment.CreatedBy,
					IsNotification: false,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				err := s.AssignRepository.CreateAssignment(ctx, &assignmentData)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return errors.New("template not found")
	}

	return nil

}

func (s *classroomService) GetTeacherAssignments(ctx context.Context, userID, organizationID string, termID string) ([]TeacherAssignmentResponse, error) {

	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if termID == "" {
		return nil, errors.New("term id is required")
	}

	if organizationID == "" {
		return nil, errors.New("organization id is required")
	}

	term, err := s.TermService.GetTermByID(ctx, termID)
	if err != nil {
		log.Printf("[ERROR] termService.GetTermByID failed (id=%s): %v", termID, err)
	}

	teacher, err := s.UserService.GetTeacherInforByOrg(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	if teacher.UserID == "" {
		return nil, errors.New("teacher not found")
	}

	startDateParse, err := time.Parse("2006-01-02", term.StartDate)
	if err != nil {
		return nil, err
	}

	endDateParse, err := time.Parse("2006-01-02", term.EndDate)
	if err != nil {
		return nil, err
	}

	assignments, err := s.AssignRepository.GetAssignmentsByStartDateAndEndDate(ctx, &startDateParse, &endDateParse)
	if err != nil {
		return nil, err
	}

	classroomMap := make(map[string]*TeacherAssignmentResponse)

	for _, a := range assignments {

		if a.TeacherID == nil {
			continue
		}

		classroomID := a.ClassRoomID.Hex()

		if _, ok := classroomMap[classroomID]; !ok {

			classroomIDParse, err := primitive.ObjectIDFromHex(classroomID)
			if err != nil {
				return nil, err
			}

			classroom, err := s.ClassroomRepository.GetClassroomByID(ctx, classroomIDParse)
			if err != nil {
				return nil, err
			}

			if classroom == nil {
				return nil, errors.New("classroom not found")
			}

			clasroomRes := ClassroomResponse{
				ID:   classroomID,
				Name: classroom.Name,
			}

			teacherInfor, err := s.UserService.GetTeacherInfor(ctx, teacher.UserID)
			if err != nil {
				return nil, err
			}

			classroomMap[classroomID] = &TeacherAssignmentResponse{
				Classroom:   clasroomRes,
				Teacher:     *teacherInfor,
				Assignments: []Assignment{},
			}
		}

		var studentInfo user.UserInfor
		if a.StudentID != nil {
			st, err := s.UserService.GetStudentInfor(ctx, *a.StudentID)
			if err == nil && st != nil {
				studentInfo = *st
			} else {
				studentInfo = user.UserInfor{
					UserID:   *a.StudentID,
					UserName: "Unknown",
				}
			}
		}

		assignmentItem := Assignment{
			ID:         a.ID.Hex(),
			AssignDate: a.AssignDate.Format("2006-01-02"),
			Student:    studentInfo,
		}

		classroomMap[classroomID].Assignments = append(classroomMap[classroomID].Assignments, assignmentItem)
	}

	var results []TeacherAssignmentResponse
	for _, v := range classroomMap {
		results = append(results, *v)
	}

	return results, nil

}

func (s *classroomService) GetClassroomByID(ctx context.Context, id, start, end string) (*ClassroomScheduleResponse, error) {

	if id == "" {
		return nil, errors.New("classroom id is required")
	}

	if start == "" {
		return nil, errors.New("start date is required")
	}

	if end == "" {
		return nil, errors.New("end date is required")
	}

	startParse, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}

	endParse, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	classroom, err := s.ClassroomRepository.GetClassroomByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if classroom == nil {
		return nil, errors.New("classroom not found")
	}

	assignments, err := s.AssignRepository.GetAssignmentsByClassroomID(ctx, objectID, &startParse, &endParse)
	if err != nil {
		return nil, err
	}

	leaderByClasses, err := s.LeaderRopitory.GetLeaderByClassID(ctx, objectID, &startParse, &endParse)
	if err != nil {
		return nil, err
	}

	scheduleMap := make(map[string]*DailySchedule)

	for _, leader := range leaderByClasses {
		if leader != nil && leader.Owner != nil {
			date := leader.Date.Format("2006-01-02")
			switch leader.Owner.OwnerRole {
			case "teacher":
				leaderInforData, err := s.UserService.GetTeacherInfor(ctx, leader.Owner.OwnerID)
				if err != nil {
					return nil, err
				}

				var leaderInfor *user.UserInfor

				if leaderInforData != nil {
					leaderInfor = &user.UserInfor{
						UserID:   leaderInforData.UserID,
						UserName: leaderInforData.UserName,
						Avartar:  leaderInforData.Avartar,
					}
				} else {
					leaderInfor = &user.UserInfor{
						UserID:   "",
						UserName: "",
						Avartar:  user.Avatar{},
					}
				}

				scheduleMap[date] = &DailySchedule{
					Date:        date,
					Leader:      leaderInfor,
					Assignments: []*SlotAssignmentResponse{},
				}

			case "staff":
				leaderInforData, err := s.UserService.GetStaffInfor(ctx, leader.Owner.OwnerID)
				if err != nil {
					return nil, err
				}

				var leaderInfor *user.UserInfor

				if leaderInforData != nil {
					leaderInfor = &user.UserInfor{
						UserID:   leaderInforData.UserID,
						UserName: leaderInforData.UserName,
						Avartar:  leaderInforData.Avartar,
					}
				} else {
					leaderInfor = &user.UserInfor{
						UserID:   "",
						UserName: "",
						Avartar:  user.Avatar{},
					}
				}
				scheduleMap[date] = &DailySchedule{
					Date:        date,
					Leader:      leaderInfor,
					Assignments: []*SlotAssignmentResponse{},
				}
			}
		}
	}

	for _, a := range assignments {
		date := a.AssignDate.Format("2006-01-02")
		if _, ok := scheduleMap[date]; !ok {
			scheduleMap[date] = &DailySchedule{
				Date:        date,
				Leader:      nil,
				Assignments: []*SlotAssignmentResponse{},
			}
		}

		var teacherInfo *user.UserInfor
		if a.TeacherID != nil && *a.TeacherID != "" {
			info, err := s.UserService.GetTeacherInfor(ctx, *a.TeacherID)
			if err == nil && info != nil {
				teacherInfo = info
			} else {
				teacherInfo = &user.UserInfor{
					UserID:   "",
					UserName: "",
					Avartar:  user.Avatar{},
				}
			}
		}

		var studentInfo *user.UserInfor
		if a.StudentID != nil && *a.StudentID != "" {
			info, err := s.UserService.GetStudentInfor(ctx, *a.StudentID)
			if err == nil && info != nil {
				studentInfo = info
			} else {
				studentInfo = &user.UserInfor{
					UserID:   "",
					UserName: "",
					Avartar:  user.Avatar{},
				}
			}
		}

		id := a.ID.Hex()

		scheduleMap[date].Assignments = append(scheduleMap[date].Assignments, &SlotAssignmentResponse{
			AssignmentID: &id,
			SlotNumber:   a.SlotNumber,
			Teacher:      teacherInfo,
			Student:      studentInfo,
		})
	}

	var schedule []*DailySchedule
	for _, v := range scheduleMap {
		schedule = append(schedule, v)
	}

	sort.Slice(schedule, func(i, j int) bool {
		return schedule[i].Date < schedule[j].Date
	})

	return &ClassroomScheduleResponse{
		ClassroomID: classroom.ID.Hex(),
		ClassName:   classroom.Name,
		Schedule:    schedule,
	}, nil

}
