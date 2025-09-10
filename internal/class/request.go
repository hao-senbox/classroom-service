package class

type CreateClassRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LocationID  *string `json:"location_id"`
	RegionID    string  `json:"region_id"`
	Note        *string `json:"note"`
	Icon        *string `json:"icon"`
}

type UpdateClassRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LocationID  *string `json:"location_id"`
}

type CreateAssginRequest struct {
	ClassRoomID string `json:"class_room_id"`
	Index       int    `json:"index"`
}

type UpdateAssginRequest struct {
	ClassroomID string  `json:"class_room_id"`
	Index       int     `json:"index"`
	Date        string  `json:"date"`
	TeacherID   *string `json:"teacher_id"`
	StudentID   *string `json:"student_id"`
}

type AddLeaderRequest struct {
	ClassroomID string `json:"class_room_id"`
	LeaderID    string `json:"leader_id"`
}

type CreateSystemNotificationRequest struct {
	Delay int `json:"delay"`
}

type UpdateSystemNotificationRequest struct {
	Delay int `json:"delay"`
}
