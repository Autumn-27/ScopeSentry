package dictionary

import (
	"context"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PortRepository interface {
	Count(ctx context.Context, search string) (int64, error)
	Find(ctx context.Context, search string, pageIndex, pageSize int) ([]models.PortDoc, error)
	Insert(ctx context.Context, name, value string) error
	Update(ctx context.Context, id, name, value string) error
	Delete(ctx context.Context, ids []string) error
}

type portRepository struct {
	collection *mongo.Collection
}

func NewPortRepository() PortRepository {
	return &portRepository{collection: mongodb.DB.Collection("PortDict")}
}

func (r *portRepository) buildSearch(search string) bson.M {
	if search == "" {
		return bson.M{}
	}
	regex := bson.M{"$regex": search, "$options": "i"}
	return bson.M{"$or": []bson.M{{"name": regex}, {"value": regex}}}
}

func (r *portRepository) Count(ctx context.Context, search string) (int64, error) {
	return r.collection.CountDocuments(ctx, r.buildSearch(search))
}

func (r *portRepository) Find(ctx context.Context, search string, pageIndex, pageSize int) ([]models.PortDoc, error) {
	filter := r.buildSearch(search)
	opts := options.Find().SetSkip(int64((pageIndex - 1) * pageSize)).SetLimit(int64(pageSize))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []models.PortDoc
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *portRepository) Insert(ctx context.Context, name, value string) error {
	_, err := r.collection.InsertOne(ctx, bson.M{"name": name, "value": value})
	return err
}

func (r *portRepository) Update(ctx context.Context, id, name, value string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"name": name, "value": value}})
	return err
}

func (r *portRepository) Delete(ctx context.Context, ids []string) error {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, oid)
	}
	_, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	return err
}
