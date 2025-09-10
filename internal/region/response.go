package region

import (
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

	type RegionResponse struct {
		ID         primitive.ObjectID `json:"id" bson:"_id"`
		Name       string             `json:"name" bson:"name"`
		Classrooms []ClassRoom        `json:"classrooms" bson:"classrooms"`
		CreatedBy  string             `json:"created_by" bson:"created_by"`
		CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	}

	type ClassRoom struct {
		ID          primitive.ObjectID          `json:"id" bson:"_id"`
		RegionID    *primitive.ObjectID         `json:"region_id" bson:"region_id"`
		Name        string                      `json:"name" bson:"name"`
		Description *string                     `json:"description" bson:"description"`
		Icon        *string                     `json:"icon" bson:"icon"`
		Note        *string                     `json:"note" bson:"note"`
		Assignments []*TeacherStudentAssignment `json:"assignments" bson:"assignments"`
		Room        room.RoomInfor              `json:"room" bson:"room"`
		IsActive    bool                        `json:"is_active" bson:"is_active"`
		CreatedBy   string                      `json:"created_by" bson:"created_by"`
		CreatedAt   time.Time                   `json:"created_at" bson:"created_at"`
		UpdatedAt   time.Time                   `json:"updated_at" bson:"updated_at"`
	}

	type TeacherStudentAssignment struct {
		ID             primitive.ObjectID `json:"id" bson:"_id"`
		ClassRoomID    primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
		Teacher        user.UserInfor     `json:"teacher" bson:"teacher_id"`
		Student        user.UserInfor     `json:"student" bson:"student_id"`
		CreatedBy      string             `json:"created_by" bson:"created_by"`
		IsNotification bool               `json:"is_notification" bson:"is_notification"`
		CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	}
