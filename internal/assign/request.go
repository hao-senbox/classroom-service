package assign

type UpdateAssginRequest struct {
	TeacherID   *string `json:"teacher_id"`
	StudentID   *string `json:"student_id"`
	ClassroomID string  `json:"class_room_id"`
	SlotNumber  int     `json:"slot_number"`
	Date        string  `json:"date"`
}
