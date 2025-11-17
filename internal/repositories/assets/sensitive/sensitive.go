package sensitive

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	CountDocuments(ctx context.Context, filter bson.M) (int64, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]models.SensitiveInfo, error)
	AggregateRaw(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error)
	FindBodyByMD5(ctx context.Context, md5 string) (string, error)
	AggregateSIDStat(ctx context.Context, pipeline mongo.Pipeline) ([]models.SensitiveSIDStat, error)
	AggregateMatchInfo(ctx context.Context, pipeline mongo.Pipeline) ([]string, error)
	FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Sensitive, error)
}

type repository struct {
	collection     *mongo.Collection
	BodyCollection *mongo.Collection
	RuleCollection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection:     mongodb.DB.Collection("SensitiveResult"),
		BodyCollection: mongodb.DB.Collection("SensitiveBody"),
		RuleCollection: mongodb.DB.Collection("SensitiveRule"),
	}
}

func (r *repository) FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Sensitive, error) {
	// 执行查询
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var results []models.Sensitive
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

func (r *repository) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count error: %w", err)
	}
	return count, nil
}

func (r *repository) Aggregate(ctx context.Context, pipeline mongo.Pipeline) ([]models.SensitiveInfo, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation error: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.SensitiveInfo
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return results, nil
}

func (r *repository) AggregateRaw(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) FindBodyByMD5(ctx context.Context, md5 string) (string, error) {
	filter := bson.M{"md5": md5}

	var result struct {
		Body string `bson:"body"`
	}
	err := r.BodyCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", fmt.Errorf("failed to find content: %w", err)
	}
	return result.Body, nil
}

func (r *repository) AggregateSIDStat(ctx context.Context, pipeline mongo.Pipeline) ([]models.SensitiveSIDStat, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.SensitiveSIDStat
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	return results, nil
}

func (r *repository) AggregateMatchInfo(ctx context.Context, pipeline mongo.Pipeline) ([]string, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation error: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Matches []string `bson:"unique_matches"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(result) == 0 {
		return []string{}, nil
	}

	return result[0].Matches, nil
}
