// Package task -----------------------------
// @file      : task.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/4 22:04
// -------------------------------------------
package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/database/redis"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
	redisdriver "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义任务数据访问接口
type Repository interface {
	Count(ctx *gin.Context, filter interface{}) (int64, error)
	Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]models.Task, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error)
	Insert(ctx *gin.Context, task *models.Task) (string, error) // 新增方法
	ClearTaskCache(ctx context.Context, id string) error
	PushTaskInfoList(ctx context.Context, id string, targetList []string) error
	FindTemplateByID(ctx context.Context, id string) (*models.ScanTemplate, error)
	GetCollectionCursor(ctx *gin.Context, collection string, filter bson.M) (*mongo.Cursor, error)
	RPushNodeTask(ctx context.Context, name string, data any) error
	CheckTaskNameExists(ctx context.Context, name string) (bool, error)
	DeleteMany(ctx *gin.Context, ids []primitive.ObjectID) (int64, error)
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) error
	GetProgressFromRedis(ctx *gin.Context, key string) (map[string]string, error)
	// 新增Redis操作方法
	Exists(ctx context.Context, key string) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

// repository 实现Repository接口
type repository struct {
	db          *mongo.Database
	redisClient *redisdriver.Client
}

// NewRepository 创建任务数据访问实例
func NewRepository() Repository {
	return &repository{
		db:          mongodb.DB,
		redisClient: redis.Client,
	}
}

// Count 统计文档数量
func (r *repository) Count(ctx *gin.Context, filter interface{}) (int64, error) {
	return r.db.Collection("task").CountDocuments(ctx.Request.Context(), filter)
}

// Find 查询多个文档
func (r *repository) Find(ctx context.Context, filter interface{}, opts *options.FindOptions) ([]models.Task, error) {
	cursor, err := r.db.Collection("task").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindByID 根据ID查询单个任务
func (r *repository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Task, error) {
	var task models.Task
	err := r.db.Collection("task").FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没找到返回 nil
		}
		return nil, fmt.Errorf("failed to find task by ID: %w", err)
	}
	return &task, nil
}

// Insert 插入一个新任务文档
func (r *repository) Insert(ctx *gin.Context, task *models.Task) (string, error) {
	result, err := r.db.Collection("task").InsertOne(context.Background(), task)
	if err != nil {
		return "", err
	}
	// 尝试将 InsertedID 转换为 ObjectID 并返回其 Hex 字符串
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to convert InsertedID to ObjectID")
}

func (r *repository) ClearTaskCache(ctx context.Context, id string) error {
	keysToDelete := []string{
		fmt.Sprintf("TaskInfo:tmp:%s", id),
		fmt.Sprintf("TaskInfo:%s", id),
		fmt.Sprintf("TaskInfo:time:%s", id),
	}

	// 查找 progress 相关 key
	progressKeys, err := r.redisClient.Keys(ctx, fmt.Sprintf("TaskInfo:progress:%s:*", id)).Result()
	if err != nil {
		return fmt.Errorf("failed to get progress keys: %w", err)
	}
	keysToDelete = append(keysToDelete, progressKeys...)

	// 查找 duplicates 相关 key
	duplicateKeys, err := r.redisClient.Keys(ctx, fmt.Sprintf("duplicates:%s:*", id)).Result()
	if err != nil {
		return fmt.Errorf("failed to get duplicate keys: %w", err)
	}
	keysToDelete = append(keysToDelete, duplicateKeys...)

	// 删除所有 key
	if len(keysToDelete) > 0 {
		if err := r.redisClient.Del(ctx, keysToDelete...).Err(); err != nil {
			return fmt.Errorf("failed to delete redis keys: %w", err)
		}
	}
	return nil
}

func (r *repository) PushTaskInfoList(ctx context.Context, id string, targetList []string) error {
	if len(targetList) == 0 {
		return nil
	}

	// 将字符串切片转为 interface{} 类型的切片
	values := make([]interface{}, len(targetList))
	for i, v := range targetList {
		values[i] = v
	}

	// 推入 Redis 列表
	if err := r.redisClient.LPush(ctx, fmt.Sprintf("TaskInfo:%s", id), values...).Err(); err != nil {
		return fmt.Errorf("failed to lpush TaskInfo list: %w", err)
	}

	return nil
}

func (r *repository) FindTemplateByID(ctx context.Context, id string) (*models.ScanTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	var result models.ScanTemplate
	err = r.db.Collection("ScanTemplates").FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没找到返回 nil
		}
		return nil, fmt.Errorf("failed to find template: %w", err)
	}

	return &result, nil
}

// GetCollectionCursor 获取集合游标
func (r *repository) GetCollectionCursor(ctx *gin.Context, collection string, filter bson.M) (*mongo.Cursor, error) {
	return r.db.Collection(collection).Find(ctx.Request.Context(), filter)
}

func (r *repository) RPushNodeTask(ctx context.Context, name string, data any) error {
	key := fmt.Sprintf("NodeTask:%s", name)

	// 将 data 序列化为 JSON
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 使用 Redis RPush 插入数据
	return r.redisClient.RPush(ctx, key, value).Err()
}

// CheckTaskNameExists 检查任务名是否存在
func (r *repository) CheckTaskNameExists(ctx context.Context, name string) (bool, error) {
	filter := bson.M{"name": name}
	count, err := r.db.Collection("task").CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check task name: %w", err)
	}
	return count > 0, nil
}

func (r *repository) DeleteMany(ctx *gin.Context, ids []primitive.ObjectID) (int64, error) {
	result, err := r.db.Collection("task").DeleteMany(ctx.Request.Context(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *repository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.db.Collection("task").UpdateOne(ctx, filter, update)
	return err
}

// GetProgressFromRedis 从Redis获取进度信息
func (r *repository) GetProgressFromRedis(ctx *gin.Context, key string) (map[string]string, error) {
	result, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get progress from redis: %w", err)
	}
	return result, nil
}

// Exists 检查Redis键是否存在
func (r *repository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return result > 0, nil
}

// SCard 获取集合中元素的数量
func (r *repository) SCard(ctx context.Context, key string) (int64, error) {
	result, err := r.redisClient.SCard(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get set cardinality: %w", err)
	}
	return result, nil
}

// Get 获取Redis键的值
func (r *repository) Get(ctx context.Context, key string) (string, error) {
	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redisdriver.Nil {
			return "", nil // 键不存在返回空字符串
		}
		return "", fmt.Errorf("failed to get key value: %w", err)
	}
	return result, nil
}

// Del 删除Redis键
func (r *repository) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	_, err := r.redisClient.Del(ctx, keys...).Result()
	if err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}
	return nil
}
