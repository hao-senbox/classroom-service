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
