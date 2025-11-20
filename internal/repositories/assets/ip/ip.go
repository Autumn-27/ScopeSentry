package ip

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义 IP 资产仓库接口
type Repository interface {
	Count(ctx context.Context, filter bson.M) (int64, error)
	Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.IPAsset, error)
	FindWithAggregation(ctx context.Context, pipeline mongo.Pipeline) ([]models.IPAsset, error)
	CountWithAggregation(ctx context.Context, pipeline mongo.Pipeline) (int64, error)
}

type repository struct {
	collection *mongo.Collection
}

// NewRepository 创建 IP 资产仓库实例
func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection("IPAsset"),
	}
}

// Count 统计文档数量
func (r *repository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count ip assets: %w", err)
	}
	return count, nil
}

// Find 分页查询文档
func (r *repository) Find(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.IPAsset, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find ip assets: %w", err)
	}
	defer cursor.Close(ctx)

	var assets []models.IPAsset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, fmt.Errorf("failed to decode ip assets: %w", err)
	}

	return assets, nil
}

// FindWithAggregation 使用聚合管道查询文档
func (r *repository) FindWithAggregation(ctx context.Context, pipeline mongo.Pipeline) ([]models.IPAsset, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate ip assets: %w", err)
	}
	defer cursor.Close(ctx)

	var assets []models.IPAsset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, fmt.Errorf("failed to decode ip assets: %w", err)
	}

	return assets, nil
}

// CountWithAggregation 使用聚合管道统计文档数量
func (r *repository) CountWithAggregation(ctx context.Context, pipeline mongo.Pipeline) (int64, error) {
	// 添加 $count 阶段
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "total"}})

	cursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count ip assets: %w", err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("failed to decode count result: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	if total, ok := result[0]["total"].(int32); ok {
		return int64(total), nil
	}
	if total, ok := result[0]["total"].(int64); ok {
		return total, nil
	}

	return 0, fmt.Errorf("unexpected count result type")
}
