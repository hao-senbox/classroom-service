package region

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegionRepository interface {
	GetRegions(ctx context.Context) ([]*Region, error)
	GetRegion(ctx context.Context, id primitive.ObjectID) (*Region, error)
	CreateRegion(ctx context.Context, data *Region) error
	UpdateRegion(ctx context.Context, id primitive.ObjectID, data *Region) error
	DeleteRegion(ctx context.Context, id primitive.ObjectID) error
}

type regionRepository struct {
	regionCollection *mongo.Collection
}

func NewRegionRepository(collection *mongo.Collection) RegionRepository {
	return &regionRepository{
		regionCollection: collection,
	}
}

func (r *regionRepository) GetRegions(ctx context.Context) ([]*Region, error) {

	var regions []*Region

	cursor, err := r.regionCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var region Region
		if err := cursor.Decode(&region); err != nil {
			return nil, err
		}
		regions = append(regions, &region)
	}

	return regions, nil

}

func (r *regionRepository) GetRegion(ctx context.Context, id primitive.ObjectID) (*Region, error) {

	var region Region

	if err := r.regionCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&region); err != nil {
		return nil, err
	}

	return &region, nil

}

func (r *regionRepository) CreateRegion(ctx context.Context, data *Region) error {

	_, err := r.regionCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil

}

func (r *regionRepository) UpdateRegion(ctx context.Context, id primitive.ObjectID, data *Region) error {

	_, err := r.regionCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
	if err != nil {
		return err
	}

	return nil
	
}

func (r *regionRepository) DeleteRegion(ctx context.Context, id primitive.ObjectID) error {

	_, err := r.regionCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil

}