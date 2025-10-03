package leader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Leader struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Owner       Owner              `json:"owner" bson:"owner"`
	Date        time.Time          `json:"date" bson:"date"`
	ClassRoomID primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Owner struct {
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	OwnerRole string `json:"owner_role" bson:"owner_role"`
}

type LeaderTemplate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Owner       *Owner              `json:"owner" bson:"owner"`
	ClassRoomID primitive.ObjectID `json:"class_room_id" bson:"class_room_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
