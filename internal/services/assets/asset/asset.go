package asset

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/assets/asset"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义资产服务接口
type Service interface {
	GetAssets(ctx *gin.Context, query models.SearchRequest) ([]models.Asset, error)
	GetAssetByIDs(ctx *gin.Context, ids []string) ([]models.Asset, error)
	GetAssetByID(ctx *gin.Context, id string) (*models.Asset, error)
	CreateAsset(ctx *gin.Context, asset *models.Asset) error
	UpdateAsset(ctx *gin.Context, id string, asset *models.Asset) error
	DeleteAsset(ctx *gin.Context, id string) error
	GetScreenshot(ctx *gin.Context, id string) (string, error)
	GetChangeLog(ctx *gin.Context, id string) ([]models.AssetChangeLog, error)
	DeduplicateAssets(ctx *gin.Context, assetType string, filter bson.M, groupFields []string) error
	GetAssetCardData(ctx *gin.Context, query models.SearchRequest) ([]models.Asset, error)
	GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error)
	GetTaskTargetByIDs(ctx *gin.Context, ids []string) (string, error)
	BuildSearchQuery(query models.SearchRequest) (map[string]interface{}, error)
}

// 全局 bodyhash 缓存实例，TTL 为 30 秒（可通过 NewCache 参数自定义）
var globalBodyHashCache = helper.NewCache(30 * time.Second)

type service struct {
	assetRepo asset.Repository
}

// NewService 创建资产服务实例
func NewService() Service {
	return &service{
		assetRepo: asset.NewRepository(),
	}
}

func (s *service) ReplaceBodyQuery(ctx context.Context, query map[string]interface{}) map[string]interface{} {
	return s.replaceBodyRecursively(ctx, query).(map[string]interface{})
}

// convertToInterfaceSlice 将各种类型的数组转换为 []interface{}
func convertToInterfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	// 使用反射处理各种数组类型
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		result := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}
		return result
	}

	// 如果已经是 []interface{}，直接返回
	if arr, ok := v.([]interface{}); ok {
		return arr
	}

	return nil
}

// replaceBodyRecursively 内部递归处理函数
func (s *service) replaceBodyRecursively(ctx context.Context, query any) any {
	switch q := query.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range q {
			// 逻辑操作符递归处理
			if k == "$and" || k == "$or" || k == "$nor" {
				// 处理多种可能的数组类型
				arr := convertToInterfaceSlice(v)
				if arr != nil {
					newArr := []interface{}{}
					for _, sub := range arr {
						newArr = append(newArr, s.replaceBodyRecursively(ctx, sub))
					}
					newMap[k] = newArr
					continue
				}
			}

			// ✅ 特殊处理 body 字段
			if k == "body" {
				cond := s.convertBodyToBodyhash(ctx, v)
				// 将转换后的条件合并到 newMap 中
				// 如果 cond 为空（否定查询且无匹配），不添加任何条件，表示不限制
				for hk, hv := range cond {
					newMap[hk] = hv
				}
				continue
			}

			// 普通字段递归
			newMap[k] = s.replaceBodyRecursively(ctx, v)
		}
		return newMap

	case []interface{}:
		newArr := []interface{}{}
		for _, sub := range q {
			newArr = append(newArr, s.replaceBodyRecursively(ctx, sub))
		}
		return newArr

	default:
		return q
	}
}

// convertBodyToBodyhash 将 body 条件转成 bodyhash 查询（含全局缓存与 != 支持）
func (s *service) convertBodyToBodyhash(ctx context.Context, v interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	m, ok := v.(map[string]interface{})
	if !ok {
		return result
	}

	var (
		regexVal   string
		isNegation bool
		cacheKey   string
		bodyHashes []string
	)

	// ✅ body = "xxx" 或 body: { "$regex": "xxx", "$options": "i" }
	if regex, ok := m["$regex"].(string); ok {
		regexVal = regex
		cacheKey = fmt.Sprintf("eq|%s", regexVal)
		isNegation = false
	} else if not, ok := m["$not"].(map[string]interface{}); ok {
		// ✅ body != "xxx" 或 body: { "$not": { "$regex": "xxx", "$options": "i" } }
		if regex, ok := not["$regex"].(string); ok {
			regexVal = regex
			isNegation = true
			cacheKey = fmt.Sprintf("ne|%s", regexVal)
		}
	}

	if regexVal == "" {
		// 无法识别的 body 查询格式，返回空结果
		return result
	}

	// 先查全局缓存
	if hashes, ok := globalBodyHashCache.GetStringSlice(cacheKey); ok {
		bodyHashes = hashes
	} else {
		// 没缓存才去数据库查
		filter := bson.M{"content": bson.M{"$regex": regexVal, "$options": "i"}}
		opts := options.Find().SetProjection(bson.M{"hash": 1}).SetLimit(1000)

		bodies, err := s.assetRepo.FindBodyWithOptions(ctx, filter, opts)
		if err != nil {
			// 查询出错时，对于否定查询返回空（表示不限制），对于正向查询返回不匹配的条件
			if isNegation {
				return map[string]interface{}{}
			}
			// 正向查询出错时，返回一个不可能匹配的条件
			return map[string]interface{}{"bodyhash": bson.M{"$in": []string{}}}
		}

		for _, b := range bodies {
			if b.BodyHash != "" {
				bodyHashes = append(bodyHashes, b.BodyHash)
			}
		}

		// 存入全局缓存
		globalBodyHashCache.SetStringSlice(cacheKey, bodyHashes)
	}

	// 构建返回条件
	if len(bodyHashes) == 0 {
		if isNegation {
			// 没匹配到内容，否定查询表示 bodyhash 可以是任意值（包括不存在的）
			// 返回空 map，不添加任何限制条件
			return map[string]interface{}{}
		}
		// 没匹配到正向内容，返回一个不可能匹配的条件（空数组 $in）
		return map[string]interface{}{"bodyhash": bson.M{"$in": []string{}}}
	}

	if isNegation {
		result["bodyhash"] = bson.M{"$nin": bodyHashes}
	} else {
		result["bodyhash"] = bson.M{"$in": bodyHashes}
	}

	return result
}

func (s *service) BuildSearchQuery(query models.SearchRequest) (map[string]interface{}, error) {
	query.Index = "asset"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	return s.ReplaceBodyQuery(context.Background(), searchQuery), nil

}

func (s *service) GetAssetByIDs(ctx *gin.Context, ids []string) ([]models.Asset, error) {
	return s.assetRepo.GetAssetByIDs(ctx, ids)
}

// GetAssets 获取资产列表
func (s *service) GetAssets(ctx *gin.Context, query models.SearchRequest) ([]models.Asset, error) {
	// 构建查询条件
	searchQuery, err := s.BuildSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	// 分页查询
	skip := (query.PageIndex - 1) * query.PageSize
	limit := query.PageSize

	// 构建排序
	sort := bson.D{
		{
			Key: "time", Value: -1,
		},
	}

	// 构建投影
	projection := bson.M{
		"_id":          1,
		"host":         1,
		"url":          1,
		"ip":           1,
		"port":         1,
		"service":      1,
		"type":         1,
		"title":        1,
		"statuscode":   1,
		"rawheaders":   1,
		"technologies": 1,
		"metadata":     1,
		"time":         1,
		"faviconmmh3":  1,
		"tags":         1,
		"bodyhash":     1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(sort)

	return s.assetRepo.FindWithOptions(ctx.Request.Context(), filter, opts)
}

// GetAssetByID 根据ID获取资产
func (s *service) GetAssetByID(ctx *gin.Context, id string) (*models.Asset, error) {
	assetByID, err := s.assetRepo.GetAssetByID(ctx.Request.Context(), id)
	if err != nil {
		return nil, err
	}
	httpBody, err := s.assetRepo.FindBodyByHash(context.Background(), assetByID.ResponseBodyHash)
	if err != nil {
		return nil, err
	}
	assetByID.Body = httpBody.Content
	return assetByID, nil
}

// CreateAsset 创建资产
func (s *service) CreateAsset(ctx *gin.Context, asset *models.Asset) error {
	return s.assetRepo.CreateAsset(ctx.Request.Context(), asset)
}

// UpdateAsset 更新资产
func (s *service) UpdateAsset(ctx *gin.Context, id string, asset *models.Asset) error {
	return s.assetRepo.UpdateAsset(ctx.Request.Context(), id, asset)
}

// DeleteAsset 删除资产
func (s *service) DeleteAsset(ctx *gin.Context, id string) error {
	return s.assetRepo.DeleteAsset(ctx.Request.Context(), id)
}

// GetScreenshot 获取资产截图
func (s *service) GetScreenshot(ctx *gin.Context, id string) (string, error) {
	return s.assetRepo.GetScreenshot(ctx.Request.Context(), id)
}

// GetChangeLog 获取资产变更日志
func (s *service) GetChangeLog(ctx *gin.Context, id string) ([]models.AssetChangeLog, error) {
	return s.assetRepo.GetChangeLog(ctx.Request.Context(), id)
}

// DeduplicateAssets 资产去重
func (s *service) DeduplicateAssets(ctx *gin.Context, assetType string, filter bson.M, groupFields []string) error {
	return s.assetRepo.DeduplicateAssets(ctx.Request.Context(), assetType, filter, groupFields)
}

// GetAssetCardData 获取资产卡片数据
func (s *service) GetAssetCardData(ctx *gin.Context, query models.SearchRequest) ([]models.Asset, error) {
	// 构建查询条件
	query.Index = "asset"
	searchQuery, err := s.BuildSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	// 分页查询
	skip := (query.PageIndex - 1) * query.PageSize
	limit := query.PageSize

	// 构建排序
	sort := bson.D{
		{
			Key: "time", Value: -1,
		},
	}

	// 构建投影
	projection := bson.M{
		"_id":        1,
		"host":       1,
		"url":        1,
		"port":       1,
		"service":    1,
		"type":       1,
		"title":      1,
		"statuscode": 1,
		"bodyhash":   1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(sort)

	assets, err := s.assetRepo.FindWithOptions(ctx.Request.Context(), filter, opts)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *service) GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error) {
	query.Index = "asset"
	searchQuery, err := s.BuildSearchQuery(query)
	if err != nil {
		return "", err
	}
	filter := bson.M(searchQuery)
	projection := bson.M{
		"url":     1,
		"host":    1,
		"port":    1,
		"service": 1,
		"type":    1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection)

	if query.PageSize != 0 {
		opts = opts.SetLimit(int64(query.PageSize))
	}

	assets, err := s.assetRepo.FindWithOptions(ctx.Request.Context(), filter, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i, a := range assets {
		if a.Type == "http" {
			builder.WriteString(strings.TrimRightFunc(a.URL, unicode.IsSpace))
		} else {
			if a.Service != "" {
				builder.WriteString(a.Service + "://")
			} else {
				builder.WriteString("http" + "://")
			}
			builder.WriteString(a.Host + ":")
			builder.WriteString(a.Port)
		}
		if i != len(assets)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}

func (s *service) GetTaskTargetByIDs(ctx *gin.Context, ids []string) (string, error) {
	assets, err := s.assetRepo.GetAssetByIDs(ctx, ids)
	if err != nil {
		return "", err
	}
	targets := ""
	for _, asset := range assets {
		if asset.Type == "http" {
			targets += asset.URL + "\n"
		} else {
			sv := "http"
			if asset.Service != "" {
				sv = asset.Service
			}
			targets += sv + "://" + asset.Host + ":" + asset.Port + "\n"
		}
	}
	return targets, nil
}
