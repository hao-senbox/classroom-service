package class

import (
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationResponse struct {
	Notifications []*Notification `json:"notifications"`
	Unread        int             `json:"unread"`
}

type ClassRoomResponse struct {
	Room           room.RoomInfor                      `json:"room"`
	Assigns        []*TeacherStudentAssignmentResponse `json:"assigns"`
	LeaderResponse *user.UserInfor                     `json:"leader"`
}

type TeacherStudentAssignmentResponse struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Index          int                `json:"index" bson:"index"`
	ClassRoomID    primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	Teacher        *user.UserInfor    `json:"teacher" bson:"teacher_id"`
	Student        *user.UserInfor    `json:"student" bson:"student_id"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	IsNotification bool               `json:"is_notification" bson:"is_notification"`
	CreatedAt      *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at" bson:"updated_at"`
}
