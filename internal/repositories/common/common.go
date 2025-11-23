// common-------------------------------------
// @file      : common.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/9/25 23:41
// -------------------------------------------

package common

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/database/redis"
	redisdriver "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	// MongoDB helpers
	Find(ctx context.Context, coll string, filter interface{}, opts *options.FindOptions, result interface{}) error
	FindOne(ctx context.Context, coll string, filter interface{}, projection bson.M) (bson.M, error)
	FindMany(ctx context.Context, coll string, filter interface{}, opts *options.FindOptions) ([]bson.M, error)
	Count(ctx context.Context, coll string, filter interface{}) (int64, error)
	Aggregate(ctx context.Context, coll string, pipeline mongo.Pipeline) (*mongo.Cursor, error)
	InsertOne(ctx context.Context, coll string, doc interface{}) (primitive.ObjectID, error)
	UpdateOne(ctx context.Context, coll string, filter interface{}, update interface{}) error
	UpdateMany(ctx context.Context, coll string, filter interface{}, update interface{}) (int64, error)
	Upsert(ctx context.Context, coll string, filter interface{}, doc interface{}) error
	DeleteOne(ctx context.Context, coll string, filter interface{}) error
	DeleteMany(ctx context.Context, coll string, filter interface{}) (int64, error)

	// Redis helpers
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Del(ctx context.Context, keys ...string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, fields map[string]interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	SCard(ctx context.Context, key string) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SAdd(ctx context.Context, key string, members ...interface{}) error
	LPush(ctx context.Context, key string, values interface{}) error
	RPush(ctx context.Context, key string, values ...interface{}) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}

type repository struct {
	db          *mongo.Database
	redisClient *redisdriver.Client
}

func NewRepository() Repository {
	return &repository{
		db:          mongodb.DB,
		redisClient: redis.Client,
	}
}

func (r *repository) Find(ctx context.Context, coll string, filter interface{}, opts *options.FindOptions, result interface{}) error {
	collection := r.db.Collection(coll)

	// 检查是否是单条（limit=1）
	if opts != nil && opts.Limit != nil && *opts.Limit == 1 {
		single := collection.FindOne(ctx, filter, options.FindOne().SetProjection(opts.Projection))
		if single.Err() != nil {
			return single.Err()
		}
		return single.Decode(result)
	}

	// 多条查询
	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, result)
}

// MongoDB helpers
func (r *repository) FindOne(ctx context.Context, coll string, filter interface{}, projection bson.M) (bson.M, error) {
	var out bson.M
	err := r.db.Collection(coll).FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&out)
	return out, err
}

func (r *repository) FindMany(ctx context.Context, coll string, filter interface{}, opts *options.FindOptions) ([]bson.M, error) {
	cursor, err := r.db.Collection(coll).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []bson.M
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) Count(ctx context.Context, coll string, filter interface{}) (int64, error) {
	return r.db.Collection(coll).CountDocuments(ctx, filter)
}

func (r *repository) Aggregate(ctx context.Context, coll string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return r.db.Collection(coll).Aggregate(ctx, pipeline)
}

func (r *repository) InsertOne(ctx context.Context, coll string, doc interface{}) (primitive.ObjectID, error) {
	res, err := r.db.Collection(coll).InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (r *repository) UpdateOne(ctx context.Context, coll string, filter interface{}, update interface{}) error {
	_, err := r.db.Collection(coll).UpdateOne(ctx, filter, update)
	return err
}

func (r *repository) UpdateMany(ctx context.Context, coll string, filter interface{}, update interface{}) (int64, error) {
	res, err := r.db.Collection(coll).UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *repository) Upsert(ctx context.Context, coll string, filter interface{}, doc interface{}) error {
	_, err := r.db.Collection(coll).UpdateOne(ctx, filter, bson.M{"$set": doc}, options.Update().SetUpsert(true))
	return err
}

func (r *repository) DeleteOne(ctx context.Context, coll string, filter interface{}) error {
	_, err := r.db.Collection(coll).DeleteOne(ctx, filter)
	return err
}

func (r *repository) DeleteMany(ctx context.Context, coll string, filter interface{}) (int64, error) {
	res, err := r.db.Collection(coll).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// Redis helpers
func (r *repository) Get(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}

func (r *repository) Set(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	return r.redisClient.Set(ctx, key, val, ttl).Err()
}

func (r *repository) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (r *repository) Del(ctx context.Context, keys ...string) error {
	return r.redisClient.Del(ctx, keys...).Err()
}

func (r *repository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.redisClient.HGetAll(ctx, key).Result()
}

func (r *repository) HSet(ctx context.Context, key string, fields map[string]interface{}) error {
	return r.redisClient.HSet(ctx, key, fields).Err()
}

func (r *repository) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.redisClient.HGet(ctx, key, field).Result()
}

func (r *repository) SCard(ctx context.Context, key string) (int64, error) {
	return r.redisClient.SCard(ctx, key).Result()
}

func (r *repository) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.redisClient.SMembers(ctx, key).Result()
}

func (r *repository) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.redisClient.SAdd(ctx, key, members...).Err()
}

func (r *repository) LPush(ctx context.Context, key string, values interface{}) error {
	var elems []interface{}

	v := reflect.ValueOf(values)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		// 如果是切片或数组，遍历转换为 []interface{}
		for i := 0; i < v.Len(); i++ {
			elems = append(elems, v.Index(i).Interface())
		}
	default:
		// 否则视为单个值
		elems = append(elems, values)
	}

	if len(elems) == 0 {
		return errors.New("no values to push")
	}

	return r.redisClient.LPush(ctx, key, elems...).Err()
}

func (r *repository) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.redisClient.RPush(ctx, key, values...).Err()
}

func (r *repository) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.redisClient.LRange(ctx, key, start, stop).Result()
}

func (r *repository) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.redisClient.Keys(ctx, pattern).Result()
}
