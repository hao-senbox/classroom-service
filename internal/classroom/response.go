package classroom

import (
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomResponseData struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id"`
	RegionID       *primitive.ObjectID `json:"region_id" bson:"region_id"`
	OrganizationID string              `json:"organization_id" bson:"organization_id"`
	Name           string              `json:"name" bson:"name"`
	Description    *string             `json:"description" bson:"description"`
	Icon           *string             `json:"icon" bson:"icon"`
	Note           *string             `json:"note" bson:"note"`
	Room           *room.RoomInfor     `json:"location"`
	IsActive       bool                `json:"is_active" bson:"is_active"`
	CreatedBy      string              `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
}

type ClassroomTemplateResponse struct {
	ClassroomID    string                    `json:"classroom_id,omitempty"`
	Leader         *user.UserInfor           `json:"leader"`
	SlotAssignment []*SlotAssignmentResponse `json:"slot_assignment"`
}

type SlotAssignmentResponse struct {
	AssignmentID *string         `json:"assignment_id,omitempty"`
	SlotNumber   int             `json:"slot_number"`
	Teacher      *user.UserInfor `json:"teacher"`
	Student      *user.UserInfor `json:"student"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
}

type TeacherAssignmentResponse struct {
	Teacher      user.UserInfor  `json:"teacher"`
	Assignments  []Assignment    `json:"assignments"`
	SeenStudents map[string]bool `json:"-"`
}

type ClassroomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Assignment struct {
	ID         string         `json:"id"`
	AssignDate string         `json:"assign_date"`
	Student    user.UserInfor `json:"student"`
}

type Term struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ClassroomScheduleResponse struct {
	ClassroomID string           `json:"classroom_id"`
	ClassName   string           `json:"class_name"`
	Schedule    []*DailySchedule `json:"schedule"`
	Pagination  Pagination       `json:"pagination"`
}

type DailySchedule struct {
	Date        string                    `json:"date"`
	Leader      *user.UserInfor           `json:"leader,omitempty"`
	Assignments []*SlotAssignmentResponse `json:"assignments"`
}

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
}
