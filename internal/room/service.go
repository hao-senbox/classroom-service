package room

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService interface {
	CreateRoom(ctx context.Context, request *CreateRoomRequest, userID string) (string, error)
	GetRooms(ctx context.Context) ([]*ClassRoom, error)
	GetRoom(ctx context.Context, id string) (*ClassRoom, error)
	UpdateRoom(ctx context.Context, request *UpdateRoomRequest, id string) error
	DeleteRoom(ctx context.Context, id string) error
	CreateAssgin(ctx context.Context, request *CreateAssginRequest, userID string) (string, error)
	GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error)
	GetAssgin(ctx context.Context, id string) (*TeacherStudentAssignment, error)
	UpdateAssgin(ctx context.Context, request *UpdateAssginRequest, id string) error
	DeleteAssgin(ctx context.Context, id string) error
	CreateSystemNotification(ctx context.Context, request *CreateSystemNotificationRequest) error
	GetSystemNotification(ctx context.Context) (*SystemConfig, error)
	UpdateSystemNotification(ctx context.Context, id string, request *UpdateSystemNotificationRequest) error
	CronNotifications(ctx context.Context) error
	GetNotifications(ctx context.Context) (*NotificationResponse, error)
	ReadNotification(ctx context.Context, id string) error
}

type roomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) RoomService {
	return &roomService{
		repo: repo,
	}
}

func (r *roomService) CreateRoom(ctx context.Context, request *CreateRoomRequest, userID string) (string, error) {

	if request.Name == "" {
		return "", errors.New("name is required")
	}

	if userID == "" {
		return "", errors.New("user_id is required")
	}

	var description *string
	if request.Description != nil {
		description = request.Description
	} else {
		description = nil
	}

	var locationID *primitive.ObjectID
	if request.LocationID != nil {
		obj, err := primitive.ObjectIDFromHex(*request.LocationID)
		if err != nil {
			return "", err
		}
		locationID = &obj
	} else {
		locationID = nil
	}

	room := &ClassRoom{
		ID:          primitive.NewObjectID(),
		Name:        request.Name,
		Description: *description,
		LocationID:  locationID,
		CreatedBy:   userID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	response, err := r.repo.CreateRoom(ctx, room)
	if err != nil {
		return "", err
	}

	return response, nil

}

func (r *roomService) GetRooms(ctx context.Context) ([]*ClassRoom, error) {
	return r.repo.GetRooms(ctx)
}

func (r *roomService) GetRoom(ctx context.Context, id string) (*ClassRoom, error) {

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.repo.GetRoom(ctx, obj)
}

func (r *roomService) UpdateRoom(ctx context.Context, request *UpdateRoomRequest, id string) error {

	if request.Name == "" {
		return errors.New("name is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	room, err := r.repo.GetRoom(ctx, obj)
	if err != nil {
		return err
	}

	if request.Description != nil {
		room.Description = *request.Description
	}

	if request.LocationID != nil {
		obj, err := primitive.ObjectIDFromHex(*request.LocationID)
		if err != nil {
			return err
		}
		room.LocationID = &obj
	}

	room.UpdatedAt = time.Now()
	room.Name = request.Name

	err = r.repo.UpdateRoom(ctx, obj, room)
	if err != nil {
		return err
	}

	return nil
}

func (r *roomService) DeleteRoom(ctx context.Context, id string) error {

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.repo.DeleteRoom(ctx, obj)

}

func (r *roomService) CreateAssgin(ctx context.Context, request *CreateAssginRequest, userID string) (string, error) {

	var assignID string

	obj, err := primitive.ObjectIDFromHex(request.ClassRoomID)
	if err != nil {
		return "", err
	}

	class, err := r.repo.GetRoom(ctx, obj)
	if err != nil {
		return "", err
	}

	if class == nil {
		return "", errors.New("class room not found")
	}

	if request.StudentID != nil {
		assign := &TeacherStudentAssignment{
			ID:             primitive.NewObjectID(),
			ClassRoomID:    obj,
			TeacherID:      nil,
			StudentID:      request.StudentID,
			CreatedBy:      userID,
			IsNotification: false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		assignID, err = r.repo.CreateAssgin(ctx, assign)
		if err != nil {
			return "", err
		}
	} else if request.TeacherID != nil {
		assign := &TeacherStudentAssignment{
			ID:          primitive.NewObjectID(),
			ClassRoomID: obj,
			TeacherID:   request.TeacherID,
			StudentID:   nil,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		assignID, err = r.repo.CreateAssgin(ctx, assign)
		if err != nil {
			return "", err
		}
	}

	return assignID, nil
}

func (r *roomService) GetAssgins(ctx context.Context) ([]*TeacherStudentAssignment, error) {
	return r.repo.GetAssgins(ctx)
}

func (r *roomService) GetAssgin(ctx context.Context, id string) (*TeacherStudentAssignment, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.repo.GetAssgin(ctx, obj)

}

func (r *roomService) UpdateAssgin(ctx context.Context, request *UpdateAssginRequest, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	assign, err := r.repo.GetAssgin(ctx, obj)
	if err != nil {
		return err
	}

	if request.TeacherID != nil {
		if assign.StudentID != nil {
			exists, err := r.repo.FindDuplicate(ctx, assign.ClassRoomID, *assign.StudentID, *request.TeacherID)
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
			exists, err := r.repo.FindDuplicate(ctx, assign.ClassRoomID, *request.StudentID, *assign.TeacherID)
			if err != nil {
				return err
			}

			if exists {
				return errors.New("teacher already assigned to student")
			}
		}
		assign.StudentID = request.StudentID
	}

	assign.UpdatedAt = time.Now()

	return r.repo.UpdateAssgin(ctx, obj, assign)

}

func (r *roomService) DeleteAssgin(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.repo.DeleteAssgin(ctx, obj)

}

func (r *roomService) CreateSystemNotification(ctx context.Context, request *CreateSystemNotificationRequest) error {

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

func (r *roomService) GetSystemNotification(ctx context.Context) (*SystemConfig, error) {

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

func (r *roomService) UpdateSystemNotification(ctx context.Context, id string, request *UpdateSystemNotificationRequest) error {

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

func (r *roomService) CronNotifications(ctx context.Context) error {
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

		if now.After(assign.CreatedAt) && now.Before(notifyAt) {
			fmt.Println("✅ Bây giờ nằm TRONG khoảng: Tạo notification")

			err := r.repo.CreateNotification(ctx, &Notification{
				ID:          primitive.NewObjectID(),
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
				UpdatedAt:      now,
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

func (r *roomService) GetNotifications(ctx context.Context) (*NotificationResponse, error) {
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

func (r *roomService) ReadNotification(ctx context.Context, id string) error {

	if id == "" {
		return errors.New("id is required")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return r.repo.ReadNotification(ctx, obj)

}
