package leader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Leader struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	ClassRoomID primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
