package classroom

import (
	"classroom-service/internal/user"
	"time"
)

type ClassroomTemplateResponse struct {
	ClassroomID    string                    `json:"classroom_id"`
	Leader         *user.UserInfor           `json:"leader"`
	SlotAssignment []*SlotAssignmentResponse `json:"slot_assignment"`
}

type SlotAssignmentResponse struct {
	SlotNumber   int             `json:"slot_number"`
	AssignmentID *string         `json:"assignment_id,omitempty"`
	Teacher      *user.UserInfor `json:"teacher"`
	Student      *user.UserInfor `json:"student"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
}

type TeacherAssignmentResponse struct {
	Classroom   ClassroomResponse   `json:"classroom"`
	Teacher     user.UserInfor      `json:"teacher"`
	Assignments []Assignment `json:"assignments"`
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
