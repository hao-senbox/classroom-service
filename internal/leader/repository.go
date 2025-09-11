package leader

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaderRepository interface {
	CreateLeader(ctx context.Context, leader *Leader) error
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
