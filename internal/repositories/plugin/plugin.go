package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/database/redis"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	redisdriver "github.com/redis/go-redis/v9"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义插件仓库接口
type Repository interface {
	FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Plugin, error)
	FindByModule(ctx context.Context, module string) ([]models.Plugin, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Plugin, error)
	Create(ctx context.Context, plugin *models.Plugin) (string, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteByHash(ctx context.Context, hashes []string) error
	Count(ctx context.Context, query bson.M) (int64, error)
	GetLogs(ctx context.Context, key string) (string, error)
	CleanLogs(ctx context.Context, key string) error
	CleanAllLogs(ctx context.Context) error
}

// repository 实现插件仓库接口
type repository struct {
	collection  *mongo.Collection
	redisClient *redisdriver.Client
}

// NewRepository 创建新的插件仓库实例
func NewRepository() Repository {
	return &repository{
		collection:  mongodb.DB.Collection("plugins"),
		redisClient: redis.Client,
	}
}

// FindWithPagination 分页查询插件
func (r *repository) FindWithPagination(ctx context.Context, query bson.M, opts *options.FindOptions) ([]models.Plugin, error) {
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Plugin
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

// FindByModule 根据模块查询插件
func (r *repository) FindByModule(ctx context.Context, module string) ([]models.Plugin, error) {
	query := bson.M{"module": module}
	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.Plugin
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	return results, nil
}

// FindByID 根据ID查询插件
func (r *repository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Plugin, error) {
	var result models.Plugin
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return &result, nil
}

// Create 创建插件
func (r *repository) Create(ctx context.Context, plugin *models.Plugin) (string, error) {
	result, err := r.collection.InsertOne(ctx, plugin)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}
	// 转换 ID 为字符串
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("inserted ID is not an ObjectID")
	}
	return objectID.Hex(), err
}

// Update 更新插件
func (r *repository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}

// DeleteByHash 根据哈希删除插件
func (r *repository) DeleteByHash(ctx context.Context, hashes []string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"hash": bson.M{"$in": hashes}})
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}
	return nil
}

// Count 统计插件数量
func (r *repository) Count(ctx context.Context, query bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

func (r *repository) GetLogs(ctx context.Context, key string) (string, error) {
	// 从 Redis 获取 set 中的所有日志
	logs, err := r.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get logs from Redis: %w", err)
	}

	logData := strings.Join(logs, "\n")

	return logData, nil
}

func (r *repository) CleanLogs(ctx context.Context, key string) error {
	// 从 Redis 获取 set 中的所有日志
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to get logs from Redis: %w", err)
	}
	return nil
}

func (r *repository) CleanAllLogs(ctx context.Context) error {
	// 获取所有以 logs:plugins: 开头的键
	keys, err := r.redisClient.Keys(ctx, "logs:plugins:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get log keys from Redis: %w", err)
	}

	// 如果没有匹配的键，直接返回
	if len(keys) == 0 {
		return nil
	}

	// 删除所有匹配的键
	err = r.redisClient.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete logs from Redis: %w", err)
	}
	return nil
}
