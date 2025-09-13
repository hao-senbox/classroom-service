package leader

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaderRepository interface {
	CreateLeader(ctx context.Context, leader *Leader) error
	GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) (*Leader, error)
	DeleteLeader(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) error
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

func (r *leaderRepository) GetLeaderByClassID(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) (*Leader, error) {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	filter := bson.M{
		"class_room_id": classroomID,
		"date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
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

func (r *leaderRepository) DeleteLeader(ctx context.Context, classroomID primitive.ObjectID, date *time.Time) error {

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	filter := bson.M{
		"class_room_id": classroomID,
		"date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}

	_, err := r.leaderCollection.DeleteOne(ctx, filter)
	return err

}
 