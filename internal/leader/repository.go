package leader

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaderRepository interface {
	CreateLeader(ctx context.Context, leader *Leader) error
	GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID) (*Leader, error)
	DeleteLeader(ctx context.Context, classroomID primitive.ObjectID) error
}

type leaderRepository struct {
	leaderCollection *mongo.Collection
}

func NewLeaderRepository(leaderCollection *mongo.Collection) LeaderRepository {
	return &leaderRepository{
		leaderCollection: leaderCollection,
	}
}

func (r *leaderRepository) CreateLeader(ctx context.Context, leader *Leader) error {

	filter := bson.M{
		"class_room_id": leader.ClassRoomID,
	}

	_, err := r.leaderCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	_, err = r.leaderCollection.InsertOne(ctx, leader)
	if err != nil {
		return err
	}

	return nil
}

func (r *leaderRepository) GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID) (*Leader, error) {

	filter := bson.M{
		"class_room_id": classroomID,
	}

	var leader Leader
	err := r.leaderCollection.FindOne(ctx, filter).Decode(&leader)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &leader, nil
	
}

func (r *leaderRepository) DeleteLeader(ctx context.Context, classroomID primitive.ObjectID) error {

	filter := bson.M{
		"class_room_id": classroomID,
	}

	_, err := r.leaderCollection.DeleteOne(ctx, filter)
	return err

}