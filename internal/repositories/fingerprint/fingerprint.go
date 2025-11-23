package fingerprint

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Count(ctx context.Context, filter bson.M) (int64, error)
	List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.FingerprintRule, error)
	Insert(ctx context.Context, rule *models.FingerprintRule) (string, error)
	Update(ctx context.Context, id string, update bson.M) error
	DeleteMany(ctx context.Context, ids []string) (int64, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{collection: mongodb.DB.Collection("FingerprintRules")}
}

func (r *repository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *repository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.FingerprintRule, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find error: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.FingerprintRule
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return results, nil
}

func (r *repository) Insert(ctx context.Context, rule *models.FingerprintRule) (string, error) {
	res, err := r.collection.InsertOne(ctx, rule)
	if err != nil {
		return "", fmt.Errorf("insert error: %w", err)
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *repository) Update(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}
	return nil
}

func (r *repository) DeleteMany(ctx context.Context, ids []string) (int64, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, fmt.Errorf("invalid id: %s", id)
		}
		objIDs = append(objIDs, oid)
	}
	res, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	if err != nil {
		return 0, fmt.Errorf("delete error: %w", err)
	}
	return res.DeletedCount, nil
}
