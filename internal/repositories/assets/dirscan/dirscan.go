// dirscan-------------------------------------
// @file      : dirscan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/17 15:07
// -------------------------------------------

package dirscan

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	FindWithPagination(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.DirScanResult, error)
}

const (
	DirscanCollection = "DirScanResult"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection(DirscanCollection),
	}
}

func (r *repository) FindWithPagination(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.DirScanResult, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.DirScanResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return results, nil
}
