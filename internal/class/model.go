package class

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRoom struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	Name        string              `json:"name" bson:"name"`
	Description string              `json:"description" bson:"description"`
	LocationID  *primitive.ObjectID `json:"location_id" bson:"location_id"`
	IsActive    bool                `json:"is_active" bson:"is_active"`
	CreatedBy   string              `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
}

type TeacherStudentAssignment struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Index          int                `json:"index" bson:"index"`
	ClassRoomID    primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	TeacherID      *string            `json:"teacher_id" bson:"teacher_id"`
	StudentID      *string            `json:"student_id" bson:"student_id"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	IsNotification bool               `json:"is_notification" bson:"is_notification"`
	CreatedAt      *time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at" bson:"updated_at"`
}

type SystemConfig struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	NotificationDelay int                `json:"notification_delay" bson:"notification_delay"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

type Notification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	AssignID    primitive.ObjectID `json:"assign_id" bson:"assign_id"`
	ClassRoomID primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	TeacherID   *string            `json:"teacher_id" bson:"teacher_id"`
	StudentID   *string            `json:"student_id" bson:"student_id"`
	Message     string             `json:"message" bson:"message"`
	NotifyAt    time.Time          `json:"notify_at" bson:"notify_at"`
	IsProcessed bool               `json:"is_processed" bson:"is_processed"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Leader struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	LeaderID    string             `json:"leader_id" bson:"leader_id"`
	ClassRoomID primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}