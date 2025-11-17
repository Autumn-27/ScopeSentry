package subdomain

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TakerCollection = "SubdomainTakerResult"
)

type TakerRepository interface {
	CountDocuments(ctx context.Context, filter bson.M) (int64, error)
	Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.SubdomainTakerResult, error)
}

type takerRepository struct {
	collection *mongo.Collection
}

func NewTakerRepository() TakerRepository {
	return &takerRepository{
		collection: mongodb.DB.Collection(TakerCollection),
	}
}

func (r *takerRepository) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *takerRepository) Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.SubdomainTakerResult, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.SubdomainTakerResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}
