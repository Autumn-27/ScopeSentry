package subdomain

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义子域名仓库接口
type Repository interface {
	BuildSearchQuery(query map[string]interface{}) (bson.M, error)
	CountDocuments(ctx context.Context, query bson.M) (int64, error)
	FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Subdomain, error)
}

type repository struct {
	collection *mongo.Collection
}

// NewRepository 创建子域名仓库实例
func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection("subdomain"),
	}
}

// BuildSearchQuery 构建搜索查询条件
func (r *repository) BuildSearchQuery(query map[string]interface{}) (bson.M, error) {
	if len(query) == 0 {
		return bson.M{"_id": bson.M{"$exists": true}}, nil
	}

	// 构建查询条件
	searchQuery := bson.M{}
	for key, value := range query {
		searchQuery[key] = value
	}

	return searchQuery, nil
}

// CountDocuments 统计文档数量
func (r *repository) CountDocuments(ctx context.Context, query bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, query)
}

// FindWithPagination 分页查询文档
func (r *repository) FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Subdomain, error) {
	// 执行查询
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var results []models.Subdomain
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}
