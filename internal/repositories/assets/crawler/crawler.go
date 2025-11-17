package crawler

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.CrawlerResult, error)
}

const (
	CrawlerCollection = "crawler"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection(CrawlerCollection),
	}
}

func (r *repository) Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.CrawlerResult, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.CrawlerResult
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
