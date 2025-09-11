package assign

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherStudentAssignment struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	ClassRoomID    primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	SlotNumber     int                `json:"slot_number" bson:"slot_number"`
	AssignDate     time.Time          `json:"assign_date" bson:"assign_date"`
	TeacherID      *string            `json:"teacher_id" bson:"teacher_id"`
	StudentID      *string            `json:"student_id" bson:"student_id"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	IsNotification bool               `json:"is_notification" bson:"is_notification"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
