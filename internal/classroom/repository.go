package classroom

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomRepository interface {
	CreateClassroom(ctx context.Context, data *ClassRoom) error
	UpdateClassroom(ctx context.Context, classroomID primitive.ObjectID, data *ClassRoom) error
	GetClassroomByRegion(ctx context.Context, regionID primitive.ObjectID) ([]*ClassRoom, error)
	GetClassroomByID(ctx context.Context, classroomID primitive.ObjectID) (*ClassRoom, error)
	GetClassroomsByOrgID(ctx context.Context, orgID string) ([]*ClassRoom, error)
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

func (c *classroomRepository) UpdateClassroom(ctx context.Context, classroomID primitive.ObjectID, data *ClassRoom) error {

	_, err := c.classroomCollection.UpdateOne(ctx, bson.M{"_id": classroomID}, bson.M{"$set": data})
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

func (c *classroomRepository) GetClassroomsByOrgID(ctx context.Context, orgID string) ([]*ClassRoom, error) {

	var classrooms []*ClassRoom

	cursor, err := c.classroomCollection.Find(ctx, bson.M{"organization_id": orgID})
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