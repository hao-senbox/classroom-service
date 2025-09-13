package leader

type CreateLeaderRequest struct {
	ClassroomID string `json:"classroom_id" bson:"classroom_id"`
	Date        string `json:"date" bson:"date"`
	Owner       Owner  `json:"owner" bson:"owner"`
}

type DeleteLeaderRequest struct {
	ClassroomID string `json:"classroom_id" bson:"classroom_id"`
	Date        string `json:"date" bson:"date"`
}
