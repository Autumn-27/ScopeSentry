package plugin

import (
	"context"
	"fmt"
	"sort"
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
	FindByHash(ctx context.Context, hash string) (*models.Plugin, error)
	Create(ctx context.Context, plugin *models.Plugin) (string, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) error
	UpdateByHash(ctx context.Context, hash string, update bson.M) error
	UpsertByHash(ctx context.Context, plugin *models.Plugin) (string, error)
	DeleteByHash(ctx context.Context, hashes []string) error
	Count(ctx context.Context, query bson.M) (int64, error)
	GetLogs(ctx context.Context, key string) (string, error)
	CleanLogs(ctx context.Context, key string) error
	CleanAllLogs(ctx context.Context) error
}

// pluginLogRedis 定义读取/清理插件日志所需的最小 Redis 能力。
// 依赖接口而不是具体的 *redis.Client，单元测试可以注入内存实现，不需要真实 Redis。
type pluginLogRedis interface {
	Type(ctx context.Context, key string) *redisdriver.StatusCmd
	LRange(ctx context.Context, key string, start, stop int64) *redisdriver.StringSliceCmd
	SMembers(ctx context.Context, key string) *redisdriver.StringSliceCmd
	Del(ctx context.Context, keys ...string) *redisdriver.IntCmd
	Keys(ctx context.Context, pattern string) *redisdriver.StringSliceCmd
}

// repository 实现插件仓库接口
type repository struct {
	collection  *mongo.Collection
	redisClient pluginLogRedis
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

	opts := options.Find().SetProjection(bson.M{
		"source": 0,
	})

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

// FindByHash 根据hash查询插件
func (r *repository) FindByHash(ctx context.Context, hash string) (*models.Plugin, error) {
	var result models.Plugin
	err := r.collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&result)
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

// UpdateByHash 根据hash更新插件
func (r *repository) UpdateByHash(ctx context.Context, hash string, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"hash": hash}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	return nil
}

// UpsertByHash 根据hash插入或更新插件（如果存在则更新，不存在则插入）
func (r *repository) UpsertByHash(ctx context.Context, plugin *models.Plugin) (string, error) {
	opts := options.FindOneAndReplace().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result models.Plugin
	err := r.collection.FindOneAndReplace(ctx, bson.M{"hash": plugin.Hash}, plugin, opts).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to upsert document: %w", err)
	}

	return result.ID.Hex(), nil
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

// GetLogs 读取插件日志。
//
// 新版把日志写成 Redis list（RPush），但历史版本的 key 可能是 set（SAdd）。
// 之前无条件用 LRange 读取，遇到老数据会直接报
// "WRONGTYPE Operation against a key holding the wrong kind of value"（issue #290）。
// 现在先看 key 的实际类型，list / set 都能读，其它类型返回明确的错误。
func (r *repository) GetLogs(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", nil
	}

	keyType, err := r.redisClient.Type(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get type of plugin log key %q: %w", key, err)
	}

	switch keyType {
	case "none":
		// key 不存在：没有日志，而不是错误
		return "", nil
	case "list":
		logs, err := r.redisClient.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return "", fmt.Errorf("failed to get logs from Redis: %w", err)
		}
		return strings.Join(logs, "\n"), nil
	case "set":
		// 兼容历史版本：日志以 set 存储，成员本身无顺序，这里按行排序保证输出稳定
		logs, err := r.redisClient.SMembers(ctx, key).Result()
		if err != nil {
			return "", fmt.Errorf("failed to get logs from Redis: %w", err)
		}
		sort.Strings(logs)
		return strings.Join(logs, "\n"), nil
	default:
		return "", fmt.Errorf("plugin log key %q holds unsupported redis type %q", key, keyType)
	}
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
