package page_monitoring

import (
	"context"
	"errors"
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义页面监控仓库接口
type Repository interface {
	FindHistory(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]*models.PageMonitoring, error)
	FindContent(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (*models.PageMonitoringBody, error)
	FindDiff(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (*models.PageMonitoringBody, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
	Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]models.PageMonitoringResult, error)
}

// repository 实现页面监控仓库接口
type repository struct {
	pageCollection    *mongo.Collection
	contentCollection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		pageCollection:    mongodb.DB.Collection("PageMonitoring"),
		contentCollection: mongodb.DB.Collection("PageMonitoringBody"),
	}
}

// FindHistory 获取页面监控历史记录
func (r *repository) FindHistory(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]*models.PageMonitoring, error) {
	cursor, err := r.pageCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pages []*models.PageMonitoring
	if err = cursor.All(ctx, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}

// FindContent 获取页面监控内容
func (r *repository) FindContent(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (*models.PageMonitoringBody, error) {
	var content models.PageMonitoringBody
	err := r.contentCollection.FindOne(ctx, filter, opts).Decode(&content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// FindDiff 获取页面监控差异
func (r *repository) FindDiff(ctx context.Context, filter interface{}, opts *options.FindOneOptions) (*models.PageMonitoringBody, error) {
	var content models.PageMonitoringBody
	err := r.contentCollection.FindOne(ctx, filter, opts).Decode(&content)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("content not found for the provided ID")
		}
		return nil, err
	}

	if len(content.Content) < 2 {
		return nil, fmt.Errorf("invalid content format")
	}

	return &content, nil
}

// Count 获取文档数量
func (r *repository) Count(ctx context.Context, filter interface{}) (int64, error) {
	return r.pageCollection.CountDocuments(ctx, filter)
}

// Find 查询文档
func (r *repository) Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]models.PageMonitoringResult, error) {
	cursor, err := r.pageCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.PageMonitoringResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}
