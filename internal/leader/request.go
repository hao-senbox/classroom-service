package leader

type CreateLeaderRequest struct {
	ClassroomID string `json:"classroom_id" bson:"classroom_id"`
	UserID      string `json:"user_id" bson:"user_id"`
}
