package app

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AppCollection = "app"
)

type Repository interface {
	CountDocuments(ctx context.Context, filter bson.M) (int64, error)
	Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.AppResult, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection(AppCollection),
	}
}

func (r *repository) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *repository) Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.AppResult, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.AppResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
} 