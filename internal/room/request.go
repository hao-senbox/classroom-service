package room

type CreateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LocationID  *string `json:"location_id"`
}

type UpdateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LocationID  *string `json:"location_id"`
}

type CreateAssginRequest struct {
	ClassRoomID string  `json:"class_room_id"`
	TeacherID   *string `json:"teacher_id"`
	StudentID   *string `json:"student_id"`
}

type UpdateAssginRequest struct {
	TeacherID   *string `json:"teacher_id"`
	StudentID   *string `json:"student_id"`
}

type CreateSystemNotificationRequest struct {
	Delay int `json:"delay"`
}

type UpdateSystemNotificationRequest struct {
	Delay int `json:"delay"`
}