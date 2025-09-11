package region

import (
	"classroom-service/internal/room"
	"classroom-service/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SlotAssignmentResponse struct {
	SlotNumber     int             `json:"slot_number"`
	AssignmentID   *string         `json:"assignment_id,omitempty"`
	AssignmentDate *time.Time      `json:"assignment_date,omitempty"`
	Teacher        *user.UserInfor `json:"teacher,omitempty"`
	Student        *user.UserInfor `json:"student,omitempty"`
	IsAssigned     bool            `json:"is_assigned"`
	CreatedAt      *time.Time      `json:"created_at,omitempty"`
	UpdatedAt      *time.Time      `json:"updated_at,omitempty"`
}

type ClassRoomResponse struct {
	ID          primitive.ObjectID  `json:"id"`
	RegionID    *primitive.ObjectID `json:"region_id"`
	Name        string              `json:"classroom_name"`
	Description string              `json:"description"`
	Icon        string              `json:"icon"`
	Note        string              `json:"note"`
	Room        *room.RoomInfor     `json:"location,omitempty"`
	Leader      *user.UserInfor     `json:"leader,omitempty"`
	IsActive    bool                `json:"is_active"`
	CreatedBy   string              `json:"created_by"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`

	TotalSlots     int `json:"total_slots"`
	AssignedSlots  int `json:"assigned_slots"`
	AvailableSlots int `json:"available_slots"`

	RecentAssignments []*SlotAssignmentResponse `json:"recent_assignments"`
}

type RegionResponse struct {
	ID         primitive.ObjectID   `json:"id"`
	Name       string               `json:"region_name"`
	Classrooms []*ClassRoomResponse `json:"classrooms"`
	CreatedBy  string               `json:"created_by"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}
