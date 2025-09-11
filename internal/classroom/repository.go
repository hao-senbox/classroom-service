package classroom

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomRepository interface {
	CreateClassroom(ctx context.Context, data *ClassRoom) error
	GetClassroomByRegion(ctx context.Context, regionID primitive.ObjectID) ([]*ClassRoom, error)
	GetClassroomByID(ctx context.Context, classroomID primitive.ObjectID) (*ClassRoom, error)
}

type classroomRepository struct {
	classroomCollection *mongo.Collection
}

func NewClassroomRepository(collection *mongo.Collection) ClassroomRepository {
	return &classroomRepository{
		classroomCollection: collection,
	}
}

func (c *classroomRepository) CreateClassroom(ctx context.Context, data *ClassRoom) error {

	_, err := c.classroomCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil

}

func (c *classroomRepository) GetClassroomByRegion(ctx context.Context, regionID primitive.ObjectID) ([]*ClassRoom, error) {

	var classrooms []*ClassRoom

	cursor, err := c.classroomCollection.Find(ctx, bson.M{"region_id": regionID})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var classroom ClassRoom
		if err := cursor.Decode(&classroom); err != nil {
			return nil, err
		}
		classrooms = append(classrooms, &classroom)
	}

	return classrooms, nil

}

func (c *classroomRepository) GetClassroomByID(ctx context.Context, classroomID primitive.ObjectID) (*ClassRoom, error) {

	var classroom ClassRoom

	err := c.classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
	if err != nil {
		return nil, err
	}

	return &classroom, nil
	
}
