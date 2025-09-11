package region

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Region struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
