package class

import (
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassService interface {
	GetClasses(ctx context.Context, date string) ([]*ClassRoomResponse, error)
	AddLeader(ctx context.Context, request *AddLeaderRequest) error
	GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error)
	GetAssgin(ctx context.Context, id, index, date string) (*TeacherStudentAssignmentResponse, error)
	CreateAssgin(ctx context.Context, request *UpdateAssginRequest) error
	DeleteAssgin(ctx context.Context, request *UpdateAssginRequest) error
	CreateSystemNotification(ctx context.Context, request *CreateSystemNotificationRequest) error
	GetSystemNotification(ctx context.Context) (*SystemConfig, error)
	UpdateSystemNotification(ctx context.Context, id string, request *UpdateSystemNotificationRequest) error
	CronNotifications(ctx context.Context) error
	GetNotifications(ctx context.Context) (*NotificationResponse, error)
	ReadNotification(ctx context.Context, id string) error
}

type classService struct {
	repo        ClassRepository
	roomService room.RoomService
	userService user.UserService
}

func NewClassService(repo ClassRepository,
	service room.RoomService,
	userService user.UserService) ClassService {
	return &classService{
		repo:        repo,
		roomService: service,
		userService: userService,
	}
}

func (r *classService) GetClasses(ctx context.Context, date string) ([]*ClassRoomResponse, error) {

	var data []*ClassRoomResponse

	var timeParse *time.Time
	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		timeParse = &t
	}

	rooms, err := r.roomService.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}

	for _, room := range rooms {
		objID, err := primitive.ObjectIDFromHex(room.ID)
		if err != nil {
			return nil, err
		}

		exists, err := r.repo.GetAssignmentsByClassID(ctx, objID, timeParse)
		if err != nil {
			return nil, err
		}

		var assignsFromDB []*TeacherStudentAssignment
		if !exists {
			for i := 1; i <= 15; i++ {
				assign := &TeacherStudentAssignment{
					ID:             primitive.NewObjectID(),
					Index:          i,
					ClassRoomID:    objID,
					TeacherID:      nil,
					StudentID:      nil,
					CreatedBy:      "system",
					IsNotification: false,
					CreatedAt:      timeParse,
					UpdatedAt:      timeParse,
				}
				assignsFromDB = append(assignsFromDB, assign)
			}
			if err := r.repo.CreateManyAssignments(ctx, assignsFromDB); err != nil {
				return nil, err
			}
		} else {
			assignsFromDB, err = r.repo.GetAssignmentsByClass(ctx, objID, timeParse)
			if err != nil {
				return nil, err
			}
		}

		var assignsResp []*TeacherStudentAssignmentResponse
		for _, assign := range assignsFromDB {
			var teacherInfo *user.UserInfor
			var studentInfo *user.UserInfor

			if assign.TeacherID != nil && *assign.TeacherID != "" {
				u, err := r.userService.GetUserInfor(ctx, *assign.TeacherID)
				if err != nil || u == nil {
					log.Printf("User not found: %s", *assign.TeacherID)
				} else {
					teacherInfo = u
				}
			}

			if assign.StudentID != nil && *assign.StudentID != "" {
				u, err := r.userService.GetUserInfor(ctx, *assign.StudentID)
				if err != nil || u == nil {
					log.Printf("User not found: %s", *assign.TeacherID)
				} else {
					studentInfo = u
				}
			}
			assignsResp = append(assignsResp, &TeacherStudentAssignmentResponse{
				ID:             assign.ID,
				Index:          assign.Index,
				ClassRoomID:    assign.ClassRoomID,
				Teacher:        teacherInfo,
				Student:        studentInfo,
				CreatedBy:      assign.CreatedBy,
				IsNotification: assign.IsNotification,
				CreatedAt:      assign.CreatedAt,
				UpdatedAt:      assign.UpdatedAt,
			})
		}

		// Lấy leader
		var leaderResponse *user.UserInfor
		existingLeader, err := r.repo.GetLeaderByClassID(ctx, objID)
		if err != nil {
			return nil, err
		}
		if existingLeader == nil {
			newLeader := &Leader{
				ID:          primitive.NewObjectID(),
				LeaderID:    "",
				ClassRoomID: objID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := r.repo.CreateLeader(ctx, newLeader); err != nil {
				return nil, err
			}
		} else if existingLeader.LeaderID != "" {
			if u, err := r.userService.GetUserInfor(ctx, existingLeader.LeaderID); err == nil {
				leaderResponse = u
			}
		}

		data = append(data, &ClassRoomResponse{
			Room:           *room,
			Assigns:        assignsResp,
			LeaderResponse: leaderResponse,
		})
	}

	return data, nil
}

func (r *classService) AddLeader(ctx context.Context, request *AddLeaderRequest) error {

	if request.ClassroomID == "" {
		return errors.New("classroom id is required")
	}

	if request.LeaderID == "" {
		return errors.New("leader id is required")
	}

	obj, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	leader := &Leader{
		ID:          primitive.NewObjectID(),
		LeaderID:    request.LeaderID,
		ClassRoomID: obj,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return r.repo.CreateLeader(ctx, leader)

}
func (r *classService) GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error) {
	return r.repo.GetAssgins(ctx)
}

func (r *classService) GetAssgin(ctx context.Context, id, index, date string) (*TeacherStudentAssignmentResponse, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if index == "" {
		return nil, errors.New("index is required")
	}

	parseIndex, err := strconv.Atoi(index)
	if err != nil {
		return nil, err
	}

	if date == "" {
		return nil, errors.New("date is required")
	}

	timeParse, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	assign, err := r.repo.FindAssign(ctx, obj, parseIndex, &timeParse)
	if err != nil {
		return nil, err
	}
	if assign == nil {
		return nil, errors.New("assign not found")
	}

	var teacherInfo *user.UserInfor
	if assign.TeacherID != nil && *assign.TeacherID != "" {
		u, err := r.userService.GetUserInfor(ctx, *assign.TeacherID)
		if err != nil || u == nil {
			log.Printf("User not found: %s", *assign.TeacherID)
		} else {
			teacherInfo = u
		}
	}

	var studentInfo *user.UserInfor
	if assign.StudentID != nil && *assign.StudentID != "" {
		u, err := r.userService.GetUserInfor(ctx, *assign.StudentID)
		if err != nil || u == nil {
			log.Printf("User not found: %s", *assign.StudentID)
		} else {
			studentInfo = u
		}
	}

	return &TeacherStudentAssignmentResponse{
		ID:             assign.ID,
		Index:          assign.Index,
		ClassRoomID:    assign.ClassRoomID,
		Teacher:        teacherInfo,
		Student:        studentInfo,
		CreatedBy:      assign.CreatedBy,
		IsNotification: assign.IsNotification,
		CreatedAt:      assign.CreatedAt,
		UpdatedAt:      assign.UpdatedAt,
	}, nil

}

func (r *classService) CreateAssgin(ctx context.Context, request *UpdateAssginRequest) error {

	if request.ClassroomID == "" {
		return errors.New("classroom id is required")
	}

	if request.Index == 0 {
		return errors.New("index is required")
	}

	obj, err := primitive.ObjectIDFromHex(request.ClassroomID)
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

	assign, err := r.repo.FindAssign(ctx, obj, request.Index, &dateParse)
	if err != nil {
		return err
	}
	if assign == nil {
		return errors.New("assign not found")
	}

	assignID := assign.ID

	if request.TeacherID != nil {
		if assign.StudentID != nil {
			exists, err := r.repo.FindDuplicate(ctx, obj, *assign.StudentID, *request.TeacherID)
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
			exists, err := r.repo.FindDuplicate(ctx, obj, *request.StudentID, *assign.TeacherID)
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
	assign.UpdatedAt = &now

	return r.repo.UpdateAssgin(ctx, assignID, assign)
}

func (r *classService) DeleteAssgin(ctx context.Context, request *UpdateAssginRequest) error {

	if request.ClassroomID == "" {
		return errors.New("classroom id is required")
	}

	if request.Index == 0 {
		return errors.New("index is required")
	}

	obj, err := primitive.ObjectIDFromHex(request.ClassroomID)
	if err != nil {
		return err
	}

	return r.repo.DeleteAssgin(ctx, obj, request.Index)

}

func (r *classService) CreateSystemNotification(ctx context.Context, request *CreateSystemNotificationRequest) error {

	if request.Delay == 0 {
		return errors.New("delay is required")
	}

	system := &SystemConfig{
		ID:                primitive.NewObjectID(),
		NotificationDelay: request.Delay,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	return r.repo.CreateSystemNotification(ctx, system)
}

func (r *classService) GetSystemNotification(ctx context.Context) (*SystemConfig, error) {

	result, err := r.repo.GetFirstSystemNotification(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		system := &SystemConfig{
			ID:                primitive.NewObjectID(),
			NotificationDelay: 168,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		return system, r.repo.CreateSystemNotification(ctx, system)
	} else {
		return result, nil
	}

}

func (r *classService) UpdateSystemNotification(ctx context.Context, id string, request *UpdateSystemNotificationRequest) error {

	if request.Delay == 0 {
		return errors.New("delay is required")
	}

	system, err := r.repo.GetFirstSystemNotification(ctx)
	if err != nil {
		return err
	}

	system.NotificationDelay = request.Delay
	system.UpdatedAt = time.Now()

	return r.repo.UpdateSystemNotification(ctx, system)
}

func (r *classService) CronNotifications(ctx context.Context) error {

	result, err := r.repo.FindAssignNotTeacher(ctx)
	if err != nil {
		return err
	}

	delay, err := r.GetSystemNotification(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	fmt.Printf("👉 Now: %v\n", now)
	fmt.Printf("👉 NotificationDelay (hours): %d\n", delay.NotificationDelay)

	for _, assign := range result {
		now := time.Now().UTC()
		notifyAt := assign.CreatedAt.Add(time.Hour * time.Duration(delay.NotificationDelay))

		fmt.Printf("------------------------------\n")
		fmt.Printf("Assign ID: %v\n", assign.ID)
		fmt.Printf("CreatedAt: %v\n", assign.CreatedAt)
		fmt.Printf("NotifyAt:  %v\n", notifyAt)
		fmt.Printf("Now:       %v\n", now)

		if now.After(*assign.CreatedAt) && now.Before(notifyAt) {
			fmt.Println("✅ Bây giờ nằm TRONG khoảng: Tạo notification")

			err := r.repo.CreateNotification(ctx, &Notification{
				ID:          primitive.NewObjectID(),
				AssignID:    assign.ID,
				ClassRoomID: assign.ClassRoomID,
				TeacherID:   assign.TeacherID,
				StudentID:   assign.StudentID,
				Message:     "Please assign teacher to student",
				NotifyAt:    now,
				IsProcessed: false,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return err
			}

			err = r.repo.UpdateAssgin(ctx, assign.ID, &TeacherStudentAssignment{
				ID:             assign.ID,
				ClassRoomID:    assign.ClassRoomID,
				TeacherID:      assign.TeacherID,
				StudentID:      assign.StudentID,
				CreatedBy:      assign.CreatedBy,
				IsNotification: true,
				CreatedAt:      assign.CreatedAt,
				UpdatedAt:      &now,
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println("❌ Không nằm TRONG khoảng: Bỏ qua")
		}
	}

	return nil
}

func (r *classService) GetNotifications(ctx context.Context) (*NotificationResponse, error) {
	result, err := r.repo.GetNotifications(ctx)
	if err != nil {
		return nil, err
	}

	unreadCount := 0
	for _, n := range result {
		if !n.IsProcessed {
			unreadCount++
		}
	}

	data := &NotificationResponse{
		Notifications: result,
		Unread:        unreadCount,
	}

	return data, nil
}

func (r *classService) ReadNotification(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.repo.ReadNotification(ctx, obj)

}
