package asset

import (
	"context"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository 定义资产仓库接口
type Repository interface {
	GetAssetByID(ctx context.Context, id string) (*models.Asset, error)
	GetAssetByIDs(ctx context.Context, ids []string) ([]models.Asset, error)
	CreateAsset(ctx context.Context, asset *models.Asset) error
	UpdateAsset(ctx context.Context, id string, asset *models.Asset) error
	DeleteAsset(ctx context.Context, id string) error
	GetScreenshot(ctx context.Context, id string) (string, error)
	GetChangeLog(ctx context.Context, id string) ([]models.AssetChangeLog, error)
	DeduplicateAssets(ctx context.Context, assetType string, filter bson.M, groupFields []string) error
	CountDocuments(ctx context.Context, filter bson.M) (int64, error)
	FindWithOptions(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.Asset, error)
	FindBodyByHash(ctx context.Context, hash string) (models.HttpBody, error)
	FindBodyWithOptions(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.HttpBody, error)
	GetDistinctHostsByRootDomain(ctx context.Context, rootDomain string, skip, limit int64) ([]string, int64, error)
	GetAssetsByHost(ctx context.Context, host string, skip, limit int64) ([]models.Asset, int64, error)
}

type repository struct {
	collection *mongo.Collection
	bodyColl   *mongo.Collection
}

// NewRepository 创建资产仓库实例
func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection("asset"),
		bodyColl:   mongodb.DB.Collection("HttpBody"),
	}
}

// GetAssetByID 根据ID获取资产
func (r *repository) GetAssetByID(ctx context.Context, id string) (*models.Asset, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var asset models.Asset
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&asset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

// GetAssetByIDs 根据ID数组获取资产列表
func (r *repository) GetAssetByIDs(ctx context.Context, ids []string) ([]models.Asset, error) {
	if len(ids) == 0 {
		return []models.Asset{}, nil
	}

	// 将字符串ID转换为ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	// 使用$in操作符查询多个ID
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assets []models.Asset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	return assets, nil
}

// CreateAsset 创建资产
func (r *repository) CreateAsset(ctx context.Context, asset *models.Asset) error {
	asset.ID = primitive.NewObjectID()

	result, err := r.collection.InsertOne(ctx, asset)
	if err != nil {
		return err
	}

	asset.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateAsset 更新资产
func (r *repository) UpdateAsset(ctx context.Context, id string, asset *models.Asset) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": asset}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeleteAsset 删除资产
func (r *repository) DeleteAsset(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetScreenshot 获取资产截图
func (r *repository) GetScreenshot(ctx context.Context, id string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	var asset models.Asset
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&asset)
	if err != nil {
		return "", err
	}

	return asset.Screenshot, nil
}

// GetChangeLog 获取资产变更日志
func (r *repository) GetChangeLog(ctx context.Context, id string) ([]models.AssetChangeLog, error) {
	var logs []models.AssetChangeLog

	filter := bson.M{"assetid": id}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}) // 按时间倒序

	cursor, err := r.collection.Database().Collection("AssetChangeLog").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// DeduplicateAssets 资产去重
func (r *repository) DeduplicateAssets(ctx context.Context, assetType string, filter bson.M, groupFields []string) error {
	// 构建聚合管道
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"group": "$" + groupFields[0],
			},
			"count": bson.M{"$sum": 1},
			"docs":  bson.M{"$push": "$$ROOT"},
		}}},
		{{Key: "$match", Value: bson.M{"count": bson.M{"$gt": 1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// 处理重复数据
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, result := range results {
		docs := result["docs"].([]interface{})
		// 保留第一个文档，删除其他文档
		for i := 1; i < len(docs); i++ {
			doc := docs[i].(bson.M)
			_, err := r.collection.DeleteOne(ctx, bson.M{"_id": doc["_id"]})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CountDocuments 统计文档数量
func (r *repository) CountDocuments(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

// FindWithOptions 使用选项进行查询
func (r *repository) FindWithOptions(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.Asset, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assets []models.Asset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *repository) FindBodyByHash(ctx context.Context, hash string) (models.HttpBody, error) {
	filter := bson.M{"hash": hash}

	var httpBody models.HttpBody
	err := r.bodyColl.FindOne(ctx, filter).Decode(&httpBody)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.HttpBody{}, nil // 没有匹配记录时返回空对象，不算错误
		}
		return models.HttpBody{}, err
	}
	return httpBody, nil
}

func (r *repository) FindBodyWithOptions(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.HttpBody, error) {
	cursor, err := r.bodyColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assets []models.HttpBody
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, err
	}

	return assets, nil
}

// GetDistinctHostsByRootDomain 根据 rootDomain 获取去重后的 host 列表
// 使用聚合管道优化性能：先按 time 排序，再按 host 分组去重，使用 $facet 同时获取总数和数据
func (r *repository) GetDistinctHostsByRootDomain(ctx context.Context, rootDomain string, skip, limit int64) ([]string, int64, error) {
	// 构建基础聚合管道
	basePipeline := mongo.Pipeline{
		// 第一步：匹配 rootDomain（使用索引 {rootDomain: 1} 或 {rootDomain: 1, time: -1}）
		{{Key: "$match", Value: bson.M{"rootDomain": rootDomain, "host": bson.M{"$ne": ""}}}},
		// 第二步：按 time 倒序排序（如果有复合索引 {rootDomain: 1, time: -1} 会更快）
		{{Key: "$sort", Value: bson.D{{Key: "time", Value: -1}}}},
		// 第三步：按 host 分组去重，保留每个 host 的第一个（最新的）记录
		{{Key: "$group", Value: bson.M{
			"_id":  "$host",
			"host": bson.M{"$first": "$host"},
		}}},
		// 第四步：按 host 字母顺序排序
		{{Key: "$sort", Value: bson.D{{Key: "host", Value: 1}}}},
	}

	// 使用 $facet 在一个聚合管道中同时获取总数和数据，避免两次查询
	pipeline := make(mongo.Pipeline, len(basePipeline))
	copy(pipeline, basePipeline)
	pipeline = append(pipeline, primitive.D{
		{Key: "$facet", Value: bson.M{
			// 统计总数（去重后的数量）
			"total": mongo.Pipeline{
				{{Key: "$count", Value: "count"}},
			},
			// 获取分页数据
			"data": mongo.Pipeline{
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
				{{Key: "$project", Value: bson.M{"host": 1, "_id": 0}}},
			},
		}},
	})

	// 执行聚合查询（只执行一次）
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	if len(results) == 0 {
		return []string{}, 0, nil
	}

	facetResult := results[0]

	// 提取总数
	var total int64
	if totalArr, ok := facetResult["total"].(bson.A); ok && len(totalArr) > 0 {
		if totalDoc, ok := totalArr[0].(bson.M); ok {
			if countVal, ok := totalDoc["count"].(int32); ok {
				total = int64(countVal)
			} else if countVal, ok := totalDoc["count"].(int64); ok {
				total = countVal
			}
		}
	}

	// 提取 host 列表
	var hosts []string
	if dataArr, ok := facetResult["data"].(bson.A); ok {
		for _, item := range dataArr {
			if doc, ok := item.(bson.M); ok {
				if host, ok := doc["host"].(string); ok && host != "" {
					hosts = append(hosts, host)
				}
			}
		}
	}

	return hosts, total, nil
}

// GetAssetsByHost 根据 host 查询资产列表
func (r *repository) GetAssetsByHost(ctx context.Context, host string, skip, limit int64) ([]models.Asset, int64, error) {
	// 构建查询条件
	filter := bson.M{"host": host}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 构建投影，只返回需要的字段
	projection := bson.M{
		"ip":          1,
		"port":        1,
		"service":     1,
		"statuscode":  1,
		"faviconmmh3": 1,
		"bodyhash":    1,
	}

	// 构建查询选项
	opts := options.Find().
		SetProjection(projection).
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "time", Value: -1}})

	// 执行查询
	assets, err := r.FindWithOptions(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}
