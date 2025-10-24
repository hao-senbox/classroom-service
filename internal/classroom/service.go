package classroom

import (
	"classroom-service/internal/assign"
	"classroom-service/internal/language"
	"classroom-service/internal/leader"
	"classroom-service/internal/room"
	"classroom-service/internal/term"
	"classroom-service/internal/user"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomService interface {
	CreateClassroom(ctx context.Context, req *CreateClassroomRequest, userID string) (string, error)
	UpdateClassroom(ctx context.Context, req *UpdateClassroomRequest, id string) error
	GetClassroomsByOrg(ctx context.Context) ([]*ClassroomResponseData, error)
	GetClassroomByID(ctx context.Context, id, start, end string, page, limit int) (*ClassroomScheduleResponse, error)
	//Classroom Template
	GetClassroomByIDTemplate(ctx context.Context, id, termID string) (*ClassroomTemplateResponse, error)
	CreateAssignmentByTemplate(ctx context.Context, req *CreateAssignmentByTemplateRequest) error
	GetClassroomTemplateByTermIDAndStudentID(ctx context.Context, studentID, termID string) (*ClassroomTemplateByTermIDAndStudentIDResponse, error)
	//Assignment
	GetTeacherAssignments(ctx context.Context, userID, organizationID string, termID string) ([]TeacherAssignmentResponse, error)
	GetTeacherAssignmentsByClassroomID(ctx context.Context, classroomID, teacherID, termID string) ([]*user.UserInfor, error)
	GetStudentsByTermAndClassroomID(ctx context.Context, classroomID, termID string) ([]*user.UserInfor, error)
	GetTeacherTemplateByTermIDAndStudentID(ctx context.Context, studentID, termID string) ([]*user.UserInfor, error)

	//Gateway
	GetStudentsAndTeachersClassroomTemplateByClassroomID(ctx context.Context, classroomID, termID string) (*ClassroomTemplateByTeacherAndStudent, error)
	GetClassroomTemplateByTermID(ctx context.Context, termID string) ([]*ClassroomTemplateGatewayResponse, error)
	GetClassroomTemplateByTermIDAndClassroomID(ctx context.Context, classroomID, termID string) (*ClassroomTemplateGatewayResponse, error)
}

type classroomService struct {
	ClassroomRepository ClassroomRepository
	AssignRepository    assign.AssignRepository
	UserService         user.UserService
	LeaderRopitory      leader.LeaderRepository
	LanguageService     language.MessageLanguageGateway
	TermService         term.TermService
	RoomService         room.RoomService
}

func NewClassroomService(classroomRepository ClassroomRepository,
	assignRepository assign.AssignRepository,
	userService user.UserService,
	leaderRepository leader.LeaderRepository,
	languageService language.MessageLanguageGateway,
	termService term.TermService,
	roomService room.RoomService) ClassroomService {
	return &classroomService{
		ClassroomRepository: classroomRepository,
		AssignRepository:    assignRepository,
		UserService:         userService,
		LeaderRopitory:      leaderRepository,
		LanguageService:     languageService,
		TermService:         termService,
		RoomService:         roomService,
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

	user, err := s.UserService.GetCurrentUser(ctx)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	orgID := user.OrganizationAdmin.ID

	ClassroomID := primitive.NewObjectID()

	data := &ClassRoom{
		ID:             ClassroomID,
		Name:           req.Name,
		OrganizationID: orgID,
		Description:    req.Description,
		Note:           req.Note,
		Icon:           req.Icon,
		LocationID:     locationID,
		RegionID:       regionID,
		CreatedBy:      userID,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.ClassroomRepository.CreateClassroom(ctx, data)
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

func (s *classroomService) GetClassroomsByOrg(ctx context.Context) ([]*ClassroomResponseData, error) {

	user, err := s.UserService.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	orgID := user.OrganizationAdmin.ID

	classrooms, err := s.ClassroomRepository.GetClassroomsByOrgID(ctx, orgID)
	if err != nil {
		return nil, err

	}

	data := make([]*ClassroomResponseData, 0)

	for _, classroom := range classrooms {
		roomData, err := s.RoomService.GetRoomByID(ctx, classroom.LocationID.Hex())
		if err != nil {
			log.Println(err)
		}

		var roomRes *room.RoomInfor

		if roomData == nil {
			roomRes = &room.RoomInfor{
				ID:   "",
				Name: "",
			}
		} else {
			roomRes = &room.RoomInfor{
				ID:   roomData.ID,
				Name: roomData.Name,
			}
		}

		data = append(data, &ClassroomResponseData{
			ID:          classroom.ID,
			Name:        classroom.Name,
			Icon:        classroom.Icon,
			Note:        classroom.Note,
			Room:        roomRes,
			Description: classroom.Description,
			RegionID:    classroom.RegionID,
			IsActive:    classroom.IsActive,
			CreatedBy:   classroom.CreatedBy,
			CreatedAt:   classroom.CreatedAt,
			UpdatedAt:   classroom.UpdatedAt,
		})

	}

	return data, nil

}

func (s *classroomService) GetClassroomByIDTemplate(ctx context.Context, id, termID string) (*ClassroomTemplateResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID, objectIDTerm)
	if err != nil {
		return nil, err
	}

	leader, err := s.LeaderRopitory.GetLeaderTemplateByClassID(ctx, objectID, objectIDTerm)
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

	if req.TermID == "" {
		return errors.New("term_id is required")
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

	objectTermID, err := primitive.ObjectIDFromHex(req.TermID)
	if err != nil {
		return err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID, objectTermID)
	if err != nil {
		return err
	}

	leaderTemplate, err := s.LeaderRopitory.GetLeaderTemplateByClassID(ctx, objectID, objectTermID)
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

	teacher, err := s.UserService.GetTeacherInforByOrg(ctx, userID, organizationID)
	if err != nil {
		return nil, err
	}

	if teacher == nil {
		log.Printf("[ERROR] userService.GetTeacherInforByOrg failed (id=%s): %v", userID, err)
	}

	termObjectID, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	assignments, err := s.AssignRepository.GetAssignmentTemplateByTermID(ctx, termObjectID)
	if err != nil {
		return nil, err
	}

	response := TeacherAssignmentResponse{
		Teacher:      *teacher,
		Assignments:  []Assignment{},
		SeenStudents: make(map[string]bool),
	}

	for _, a := range assignments {

		if a.TeacherID == nil || *a.TeacherID != teacher.UserID {
			continue
		}

		if a.StudentID == nil {
			continue
		}

		if response.SeenStudents[*a.StudentID] {
			continue
		}

		response.SeenStudents[*a.StudentID] = true

		var studentInfo user.UserInfor
		st, err := s.UserService.GetStudentInfor(ctx, *a.StudentID)
		if err == nil && st != nil {
			studentInfo = *st
		} else {
			studentInfo = user.UserInfor{
				UserID:   *a.StudentID,
				UserName: "Unknown",
			}
		}

		assignmentItem := Assignment{
			ID:         a.ID.Hex(),
			AssignDate: a.CreatedAt.Format("2006-01-02"),
			Student:    studentInfo,
		}

		response.Assignments = append(response.Assignments, assignmentItem)
	}

	return []TeacherAssignmentResponse{response}, nil

}

func (s *classroomService) GetClassroomByID(ctx context.Context, id, start, end string, page, limit int) (*ClassroomScheduleResponse, error) {

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

	leaderByClasses, err := s.LeaderRopitory.GetLeaderByClassID(ctx, objectID, &startParse, &endParse, page, limit)
	if err != nil {
		return nil, err
	}

	count, err := s.LeaderRopitory.CountLeaderByClassroomID(ctx, objectID, &startParse, &endParse)
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
			continue
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
		Pagination: Pagination{
			TotalCount: int64(count),
			TotalPages: int64(math.Ceil(float64(count) / float64(limit))),
			Page:       int64(page),
			Limit:      int64(limit),
		},
	}, nil

}

func (s *classroomService) GetTeacherAssignmentsByClassroomID(ctx context.Context, classroomID, teacherID, termID string) ([]*user.UserInfor, error) {

	objectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}

	term, err := s.TermService.GetTermByID(ctx, termID)
	if err != nil {
		log.Printf("[ERROR] termService.GetTermByID failed (id=%s): %v", termID, err)
	}

	if term == nil {
		return nil, fmt.Errorf("term not found")
	}

	start, err := time.Parse("2006-01-02", term.StartDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", term.EndDate)
	if err != nil {
		return nil, err
	}

	assignments, err := s.AssignRepository.GetTeacherAssignmentsByClassroomID(ctx, objectID, teacherID, &start, &end)
	if err != nil {
		return nil, err
	}

	var infor []*user.UserInfor
	seen := make(map[string]bool)

	for _, a := range assignments {

		studentID := *a.StudentID
		if seen[studentID] {
			continue
		}

		seen[studentID] = true

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

		infor = append(infor, studentInfo)
	}

	return infor, nil

}

func (s *classroomService) GetStudentsByTermAndClassroomID(ctx context.Context, classroomID, termID string) ([]*user.UserInfor, error) {

	objectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}

	term, err := s.TermService.GetTermByID(ctx, termID)
	if err != nil {
		log.Printf("[ERROR] termService.GetTermByID failed (id=%s): %v", termID, err)
	}

	if term == nil {
		return nil, fmt.Errorf("term not found")
	}

	start, err := time.Parse("2006-01-02", term.StartDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", term.EndDate)
	if err != nil {
		return nil, err
	}

	assignments, err := s.AssignRepository.GetAssignmentsByClassroomID(ctx, objectID, &start, &end)
	if err != nil {
		return nil, err
	}

	var infor []*user.UserInfor
	seen := make(map[string]bool)

	for _, a := range assignments {

		studentID := *a.StudentID
		if seen[studentID] {
			continue
		}

		seen[studentID] = true

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

		infor = append(infor, studentInfo)
	}

	return infor, nil
}

func (s *classroomService) GetStudentsAndTeachersClassroomTemplateByClassroomID(ctx context.Context, classroomID, termID string) (*ClassroomTemplateByTeacherAndStudent, error) {

	objectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}

	objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID, objectIDTerm)
	if err != nil {
		return nil, err
	}

	if assignTemplate == nil {
		log.Printf("ClassroomTemplateByTeacherAndStudent not found for classroomID=%s", classroomID)
		return &ClassroomTemplateByTeacherAndStudent{
			Teachers: []*user.UserInfor{},
			Students: []*user.UserInfor{},
		}, nil
	}

	var studentArr []*user.UserInfor
	var teacherArr []*user.UserInfor

	studentSeen := make(map[string]bool)
	teacherSeen := make(map[string]bool)

	for _, a := range assignTemplate {

		if a.StudentID != nil && *a.StudentID != "" && a.TeacherID != nil && *a.TeacherID != "" {
			if !studentSeen[*a.StudentID] {
				info, err := s.UserService.GetStudentInfor(ctx, *a.StudentID)
				if err == nil && info != nil {
					studentArr = append(studentArr, info)
					studentSeen[*a.StudentID] = true
				}
			}
		}

		if a.TeacherID != nil && *a.TeacherID != "" && a.StudentID != nil && *a.StudentID != "" {
			if !teacherSeen[*a.TeacherID] {
				info, err := s.UserService.GetTeacherInfor(ctx, *a.TeacherID)
				if err == nil && info != nil {
					teacherArr = append(teacherArr, info)
					teacherSeen[*a.TeacherID] = true
				}
			}
		}

	}

	return &ClassroomTemplateByTeacherAndStudent{
		Teachers: teacherArr,
		Students: studentArr,
	}, nil
}

func (s *classroomService) GetClassroomTemplateByTermID(ctx context.Context, termID string) ([]*ClassroomTemplateGatewayResponse, error) {

	objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByTermID(ctx, objectIDTerm)
	if err != nil {
		return nil, err
	}

	if assignTemplate == nil {
		log.Printf("ClassroomTemplateByTeacherAndStudent not found for classroomID=%s", termID)
		return []*ClassroomTemplateGatewayResponse{}, nil
	}

	classMap := make(map[string]*ClassroomTemplateGatewayResponse)

	for _, a := range assignTemplate {

		classID := a.ClassRoomID.Hex()

		if _, ok := classMap[classID]; !ok {

			class, err := s.ClassroomRepository.GetClassroomByID(ctx, a.ClassRoomID)
			if err != nil {
				log.Printf("Cannot fetch class info for class_room_id=%s: %v", classID, err)
				continue
			}

			if class == nil {
				log.Printf("Classroom not found for class_room_id=%s", classID)
				continue
			}

			classMap[classID] = &ClassroomTemplateGatewayResponse{
				ClassID:         classID,
				ClassName:       class.Name,
				ClassIcon:       *class.Icon,
				AssignTemplates: []*AssignTemplate{},
			}

		}

		if a.TeacherID != nil && *a.TeacherID != "" && a.StudentID != nil && *a.StudentID != "" {
			classMap[classID].AssignTemplates = append(classMap[classID].AssignTemplates, &AssignTemplate{
				TeacherID: a.TeacherID,
				StudentID: a.StudentID,
			})
		}

	}

	var classArr []*ClassroomTemplateGatewayResponse
	for _, v := range classMap {
		classArr = append(classArr, v)
	}

	return classArr, nil

}

func (s *classroomService) GetClassroomTemplateByTermIDAndClassroomID(ctx context.Context, classroomID, termID string) (*ClassroomTemplateGatewayResponse, error) {

	objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(classroomID)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByClassroomID(ctx, objectID, objectIDTerm)
	if err != nil {
		return nil, err
	}

	if assignTemplate == nil {
		log.Printf("ClassroomTemplateByTeacherAndStudent not found for classroomID=%s", classroomID)
		return nil, nil
	}

	classroom, err := s.ClassroomRepository.GetClassroomByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	data := &ClassroomTemplateGatewayResponse{
		ClassID:         classroomID,
		ClassName:       classroom.Name,
		ClassIcon:       *classroom.Icon,
		AssignTemplates: []*AssignTemplate{},
	}

	for _, a := range assignTemplate {
		if a.StudentID != nil && *a.StudentID != "" && a.TeacherID != nil && *a.TeacherID != "" {
			data.AssignTemplates = append(data.AssignTemplates, &AssignTemplate{
				TeacherID: a.TeacherID,
				StudentID: a.StudentID,
			})
		}
	}

	return data, nil

}

func (s *classroomService) GetClassroomTemplateByTermIDAndStudentID(ctx context.Context, studentID, termID string) (*ClassroomTemplateByTermIDAndStudentIDResponse, error) {

	// objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	// if err != nil {
	// 	return nil, err
	// }

	// assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByTermIDAndStudentID(ctx, studentID, objectIDTerm)
	// if err != nil {
	// 	return nil, err
	// }

	// if assignTemplate == nil {
	// 	log.Printf("ClassroomTemplateByTeacherAndStudent not found for classroomID=%s", termID)
	// 	return nil, nil
	// }

	// classroom, err := s.ClassroomRepository.GetClassroomByID(ctx, assignTemplate.ClassRoomID)
	// if err != nil {
	// 	return nil, err
	// }

	// var classroomData ClassRoomTemplateResponse
	// var leaderData *user.UserInfor
	// var studentData *user.UserInfor

	// if classroom == nil {
	// 	classroomData = ClassRoomTemplateResponse{
	// 		ClassID:   "",
	// 		ClassName: "",
	// 	}
	// } else {
	// 	classroomData = ClassRoomTemplateResponse{
	// 		ClassID:   classroom.ID.Hex(),
	// 		ClassName: classroom.Name,
	// 	}
	// }

	// if assignTemplate.StudentID != nil {

	// 	student, err := s.UserService.GetStudentInfor(ctx, *assignTemplate.StudentID)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if student == nil {
	// 		log.Printf("Student not found for student_id=%s", *assignTemplate.StudentID)
	// 		studentData = &user.UserInfor{
	// 			UserID:   "",
	// 			UserName: "",
	// 		}
	// 	} else {
	// 		studentData = student
	// 	}
	// }

	// leader, err := s.LeaderRopitory.GetLeaderTemplateByClassID(ctx, assignTemplate.ClassRoomID, objectIDTerm)
	// if err != nil {
	// 	return nil, err
	// }

	// if leader != nil && leader.Owner != nil {
	// 	switch leader.Owner.OwnerRole {
	// 	case "teacher":
	// 		leaderData, err = s.UserService.GetTeacherInfor(ctx, leader.Owner.OwnerID)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	case "student":
	// 		leaderData, err = s.UserService.GetStudentInfor(ctx, leader.Owner.OwnerID)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	case "staff":
	// 		leaderData, err = s.UserService.GetStaffInfor(ctx, leader.Owner.OwnerID)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
		
	// }

	// return &ClassroomTemplateByTermIDAndStudentIDResponse{
	// 	ClassRoom: &classroomData,
	// 	Leader:    leaderData,
	// 	Student:   studentData,
	// }, nil

	return nil, nil
	
}

func (s *classroomService) GetTeacherTemplateByTermIDAndStudentID(ctx context.Context, studentID, termID string) ([]*user.UserInfor, error) {

	objectIDTerm, err := primitive.ObjectIDFromHex(termID)
	if err != nil {
		return nil, err
	}

	assignTemplate, err := s.AssignRepository.GetAssignmentTemplateByTermIDAndStudentID(ctx, studentID, objectIDTerm)
	if err != nil {
		return nil, err
	}

	if assignTemplate == nil {
		log.Printf("ClassroomTemplateByTeacherAndStudent not found for classroomID=%s", termID)
		return nil, nil
	}

	var teacherArr []*user.UserInfor
	for _, a := range assignTemplate {
		if a.TeacherID != nil && *a.TeacherID != "" {
			teacher, err := s.UserService.GetTeacherInfor(ctx, *a.TeacherID)
			if err != nil {
				log.Printf("Cannot fetch teacher info for teacher_id=%s: %v", *a.TeacherID, err)
				continue
			}
			if teacher != nil {
				teacherArr = append(teacherArr, teacher)
			}
		}
	}

	return teacherArr, nil

}