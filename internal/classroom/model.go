package classroom

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRoom struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	RegionID    *primitive.ObjectID  `json:"region_id" bson:"region_id"`
	Name        string              `json:"name" bson:"name"`
	Description *string             `json:"description" bson:"description"`
	Icon        *string             `json:"icon" bson:"icon"`
	Note        *string             `json:"note" bson:"note"`
	LocationID  *primitive.ObjectID `json:"location_id" bson:"location_id"`
	IsActive    bool                `json:"is_active" bson:"is_active"`
	CreatedBy   string              `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" bson:"updated_at"`
}
