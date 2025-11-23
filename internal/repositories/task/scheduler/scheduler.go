package scheduler

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Count(ctx *gin.Context, filter interface{}) (int64, error)
	Find(ctx *gin.Context, filter interface{}, opts *options.FindOptions) ([]models.Task, error)
	FindOne(ctx *gin.Context, filter interface{}) (*models.Task, error)
	Insert(ctx *gin.Context, task *models.Task) error
	UpdateOne(ctx *gin.Context, filter interface{}, update interface{}) error
	DeleteMany(ctx *gin.Context, filter interface{}) error
	CountPageMonit(ctx *gin.Context, filter interface{}) (int64, error)
	FindPageMonit(ctx *gin.Context, filter interface{}, opts *options.FindOptions) ([]models.PageMonitoringTask, error)
	InsertPageMonit(ctx *gin.Context, task *models.PageMonitoringTask) error
	DeletePageMonit(ctx *gin.Context, filter interface{}) error
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection("ScheduledTasks"),
	}
}

// Count 统计文档数量
func (r *repository) Count(ctx *gin.Context, filter interface{}) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Find 查询文档列表
func (r *repository) Find(ctx *gin.Context, filter interface{}, opts *options.FindOptions) ([]models.Task, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Task
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

// FindOne 查询单个文档
func (r *repository) FindOne(ctx *gin.Context, filter interface{}) (*models.Task, error) {
	var result models.Task
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return &result, nil
}

// Insert 插入文档
func (r *repository) Insert(ctx *gin.Context, task *models.Task) error {
	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}
	return nil
}

// UpdateOne 更新单个文档
func (r *repository) UpdateOne(ctx *gin.Context, filter interface{}, update interface{}) error {
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}

// DeleteMany 删除多个文档
func (r *repository) DeleteMany(ctx *gin.Context, filter interface{}) error {
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	return nil
}

// CountPageMonit 统计PageMonitoring文档数量
func (r *repository) CountPageMonit(ctx *gin.Context, filter interface{}) (int64, error) {
	collection := mongodb.DB.Collection("PageMonitoring")
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// FindPageMonit 查询PageMonitoring文档列表
func (r *repository) FindPageMonit(ctx *gin.Context, filter interface{}, opts *options.FindOptions) ([]models.PageMonitoringTask, error) {
	collection := mongodb.DB.Collection("PageMonitoring")
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.PageMonitoringTask
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

// InsertPageMonit 插入PageMonitoring文档
func (r *repository) InsertPageMonit(ctx *gin.Context, task *models.PageMonitoringTask) error {
	collection := mongodb.DB.Collection("PageMonitoring")
	_, err := collection.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}
	return nil
}

// DeletePageMonit 删除PageMonitoring文档
func (r *repository) DeletePageMonit(ctx *gin.Context, filter interface{}) error {
	collection := mongodb.DB.Collection("PageMonitoring")
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	return nil
}
