package export

import (
	"context"
	"fmt"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"

	"github.com/Autumn-27/ScopeSentry/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
)

// Repository 定义导出仓库接口
type Repository interface {
	// CreateExportTask 创建导出任务
	CreateExportTask(ctx context.Context, task *models.ExportTask) error
	// UpdateExportTaskStatus 更新导出任务状态
	UpdateExportTaskStatus(ctx context.Context, fileName string, state int, fileSize float64) error
	// GetExportTasks 获取导出任务列表
	GetExportTasks(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]map[string]interface{}, error)
	// DeleteExportTasks 删除导出任务
	DeleteExportTasks(ctx context.Context, fileNames []string) error
	// GetCollectionCursor 获取数据集合游标
	GetCollectionCursor(ctx context.Context, collection string, query bson.M, opts *options.FindOptions) (*mongo.Cursor, error)
	// GetExportRecords 获取导出记录
	GetExportRecords(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.ExportTask, error)
}

// ExportTask 导出任务结构体
type ExportTask struct {
	FileName   string    `bson:"file_name"`
	CreateTime time.Time `bson:"create_time"`
	Quantity   int       `bson:"quantity"`
	DataType   string    `bson:"data_type"`
	FileType   string    `bson:"file_type"`
	State      int       `bson:"state"`
	EndTime    string    `bson:"end_time,omitempty"`
	FileSize   string    `bson:"file_size,omitempty"`
}

type repository struct {
	db *mongo.Database
}

// NewRepository 创建导出仓库实例
func NewRepository() Repository {
	return &repository{
		db: mongodb.DB,
	}
}

// CreateExportTask 创建导出任务
func (r *repository) CreateExportTask(ctx context.Context, task *models.ExportTask) error {
	_, err := r.db.Collection("export").InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to create export task: %w", err)
	}
	return nil
}

// UpdateExportTaskStatus 更新导出任务状态
func (r *repository) UpdateExportTaskStatus(ctx context.Context, fileName string, state int, fileSize float64) error {
	update := bson.M{
		"$set": bson.M{
			"state":     state,
			"end_time":  helper.GetNowTimeString(),
			"file_size": fmt.Sprintf("%.2f", fileSize),
		},
	}

	_, err := r.db.Collection("export").UpdateOne(
		ctx,
		bson.M{"file_name": fileName},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update export task status: %w", err)
	}
	return nil
}

// GetExportTasks 获取导出任务列表
func (r *repository) GetExportTasks(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]map[string]interface{}, error) {
	cursor, err := r.db.Collection("export").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find export tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode export tasks: %w", err)
	}

	return results, nil
}

// DeleteExportTasks 删除导出任务
func (r *repository) DeleteExportTasks(ctx context.Context, fileNames []string) error {
	_, err := r.db.Collection("export").DeleteMany(
		ctx,
		bson.M{"file_name": bson.M{"$in": fileNames}},
	)
	if err != nil {
		return fmt.Errorf("failed to delete export tasks: %w", err)
	}
	return nil
}

// GetCollectionCursor 获取数据集合游标
func (r *repository) GetCollectionCursor(ctx context.Context, collection string, query bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {
	// 设置游标选项
	if opts == nil {
		opts = options.Find()
	}

	// 设置批处理大小
	opts.SetBatchSize(1000) // 每次从MongoDB获取1000条数据

	// 获取游标
	cursor, err := r.db.Collection(collection).Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get cursor: %w", err)
	}

	return cursor, nil
}

// GetExportRecords 获取导出记录
func (r *repository) GetExportRecords(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.ExportTask, error) {
	cursor, err := r.db.Collection("export").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find export records: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.ExportTask
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode export records: %w", err)
	}

	return results, nil
}
