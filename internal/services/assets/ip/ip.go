package ip

import (
	"strings"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	iprepo "github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/ip"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/random"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义 IP 资产服务接口
type Service interface {
	GetIPAssets(ctx *gin.Context, query models.SearchRequest) (*models.IPAssetResponse, error)
}

type service struct {
	repo iprepo.Repository
}

// NewService 创建 IP 资产服务实例
func NewService() Service {
	return &service{
		repo: iprepo.NewRepository(),
	}
}

// GetIPAssets 获取 IP 资产列表
func (s *service) GetIPAssets(ctx *gin.Context, query models.SearchRequest) (*models.IPAssetResponse, error) {
	if query.PageIndex <= 0 {
		query.PageIndex = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	query.Index = "IPAsset"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	reqCtx := ctx.Request.Context()

	// 检测是否涉及嵌套字段查询
	hasNestedField := hasNestedFieldQuery(filter)

	var assets []models.IPAsset
	var total int64

	if hasNestedField {
		// 使用聚合管道过滤嵌套数组
		pipeline := buildAggregationPipeline(filter, query.PageIndex, query.PageSize)

		// 构建计数管道（不包含分页）
		countPipeline := buildCountPipeline(filter)
		total, err = s.repo.CountWithAggregation(reqCtx, countPipeline)
		if err != nil {
			return nil, err
		}

		assets, err = s.repo.FindWithAggregation(reqCtx, pipeline)
		if err != nil {
			return nil, err
		}
	} else {
		// 使用普通查询（只查询顶级字段，如 ip）
		total, err = s.repo.Count(reqCtx, filter)
		if err != nil {
			return nil, err
		}

		projection := bson.M{
			"ip":    1,
			"ports": 1,
			"time":  1,
		}

		opts := options.Find().
			SetProjection(projection).
			SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
			SetLimit(int64(query.PageSize)).
			SetSort(bson.D{{Key: "time", Value: -1}})

		assets, err = s.repo.Find(reqCtx, filter, opts)
		if err != nil {
			return nil, err
		}
	}

	rows := flattenIPAssets(assets)

	return &models.IPAssetResponse{
		List:  rows,
		Total: total,
	}, nil
}

// hasNestedFieldQuery 检测查询条件是否涉及嵌套字段
func hasNestedFieldQuery(filter bson.M) bool {
	nestedFields := map[string]bool{
		"ports.port":                true,
		"ports.server.domain":       true,
		"ports.server.service":      true,
		"ports.server.webServer":    true,
		"ports.server.technologies": true,
	}

	return checkNestedFieldInQuery(filter, nestedFields)
}

// checkNestedFieldInQuery 递归检查查询条件中是否包含嵌套字段
func checkNestedFieldInQuery(query interface{}, nestedFields map[string]bool) bool {
	switch v := query.(type) {
	case bson.M:
		for key, value := range v {
			// 检查 key 本身是否是嵌套字段
			if nestedFields[key] {
				return true
			}
			// 检查 key 是否以 ports. 开头
			if strings.HasPrefix(key, "ports.") {
				return true
			}
			// 递归检查值
			if checkNestedFieldInQuery(value, nestedFields) {
				return true
			}
		}
	case bson.A:
		for _, item := range v {
			if checkNestedFieldInQuery(item, nestedFields) {
				return true
			}
		}
	case bson.D:
		for _, elem := range v {
			if nestedFields[elem.Key] {
				return true
			}
			if strings.HasPrefix(elem.Key, "ports.") {
				return true
			}
			if checkNestedFieldInQuery(elem.Value, nestedFields) {
				return true
			}
		}
	case []interface{}:
		for _, item := range v {
			if checkNestedFieldInQuery(item, nestedFields) {
				return true
			}
		}
	case map[string]interface{}:
		for key, value := range v {
			if nestedFields[key] {
				return true
			}
			if strings.HasPrefix(key, "ports.") {
				return true
			}
			if checkNestedFieldInQuery(value, nestedFields) {
				return true
			}
		}
	}
	return false
}

// buildAggregationPipeline 构建聚合管道，过滤嵌套数组
func buildAggregationPipeline(filter bson.M, pageIndex, pageSize int) mongo.Pipeline {
	// 分离顶级字段和嵌套字段的查询条件
	topLevelFilter := bson.M{}
	nestedFilter := bson.M{}

	topLevelFields := map[string]bool{
		"ip":         true,
		"time":       true,
		"project":    true,
		"taskName":   true,
		"rootDomain": true,
		"_id":        true,
	}

	separateFilters(filter, topLevelFilter, nestedFilter, topLevelFields)

	pipeline := mongo.Pipeline{}

	// 1. 首先匹配顶级字段（如 ip）
	if len(topLevelFilter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: topLevelFilter}})
	}

	// 2. 展开 ports 数组
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$ports",
		"preserveNullAndEmptyArrays": true,
	}}})

	// 3. 展开 ports.server 数组
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$ports.server",
		"preserveNullAndEmptyArrays": true,
	}}})

	// 4. 匹配嵌套字段条件
	if len(nestedFilter) > 0 {
		// 将嵌套字段路径转换为展开后的路径
		expandedFilter := expandNestedFilter(nestedFilter)
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: expandedFilter}})
	}

	// 5. 按 IP 和端口分组，重新组合 server 数组
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id": bson.M{
			"ip":   "$ip",
			"port": "$ports.port",
		},
		"ip":     bson.M{"$first": "$ip"},
		"port":   bson.M{"$first": "$ports.port"},
		"time":   bson.M{"$first": "$time"},
		"origId": bson.M{"$first": "$_id"},
		"server": bson.M{"$push": "$ports.server"},
	}}})

	// 6. 按 IP 分组，重新组合 ports 数组
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id":    "$ip",
		"ip":     bson.M{"$first": "$ip"},
		"time":   bson.M{"$first": "$time"},
		"origId": bson.M{"$first": "$origId"},
		"ports": bson.M{"$push": bson.M{
			"port":   "$port",
			"server": "$server",
		}},
	}}})

	// 7. 恢复 _id 字段
	pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.M{
		"_id": "$origId",
	}}})

	// 移除临时字段
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{
		"origId": 0,
	}}})

	// 8. 排序
	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{
		{Key: "time", Value: -1},
	}}})

	// 9. 分页
	skip := (pageIndex - 1) * pageSize
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: pageSize}})

	return pipeline
}

// buildCountPipeline 构建计数聚合管道
func buildCountPipeline(filter bson.M) mongo.Pipeline {
	topLevelFilter := bson.M{}
	nestedFilter := bson.M{}

	topLevelFields := map[string]bool{
		"ip":         true,
		"time":       true,
		"project":    true,
		"taskName":   true,
		"rootDomain": true,
		"_id":        true,
	}

	separateFilters(filter, topLevelFilter, nestedFilter, topLevelFields)

	pipeline := mongo.Pipeline{}

	if len(topLevelFilter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: topLevelFilter}})
	}

	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$ports",
		"preserveNullAndEmptyArrays": true,
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$ports.server",
		"preserveNullAndEmptyArrays": true,
	}}})

	if len(nestedFilter) > 0 {
		expandedFilter := expandNestedFilter(nestedFilter)
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: expandedFilter}})
	}

	// 按 IP 去重计数
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id": "$ip",
	}}})

	return pipeline
}

// separateFilters 分离顶级字段和嵌套字段的查询条件
func separateFilters(filter bson.M, topLevelFilter, nestedFilter bson.M, topLevelFields map[string]bool) {
	for key, value := range filter {
		if topLevelFields[key] || key == "$and" || key == "$or" {
			// 处理逻辑运算符
			if key == "$and" || key == "$or" {
				topLevelArray := []interface{}{}
				nestedArray := []interface{}{}

				if arr, ok := value.([]interface{}); ok {
					for _, item := range arr {
						if itemMap, ok := item.(bson.M); ok {
							tempTop := bson.M{}
							tempNested := bson.M{}
							separateFilters(itemMap, tempTop, tempNested, topLevelFields)

							if len(tempTop) > 0 {
								topLevelArray = append(topLevelArray, tempTop)
							}
							if len(tempNested) > 0 {
								nestedArray = append(nestedArray, tempNested)
							}
						} else if itemMap, ok := item.(map[string]interface{}); ok {
							tempTop := bson.M{}
							tempNested := bson.M{}
							separateFilters(bson.M(itemMap), tempTop, tempNested, topLevelFields)

							if len(tempTop) > 0 {
								topLevelArray = append(topLevelArray, tempTop)
							}
							if len(tempNested) > 0 {
								nestedArray = append(nestedArray, tempNested)
							}
						}
					}
				}

				if len(topLevelArray) > 0 {
					topLevelFilter[key] = topLevelArray
				}
				if len(nestedArray) > 0 {
					nestedFilter[key] = nestedArray
				}
			} else {
				topLevelFilter[key] = value
			}
		} else if strings.HasPrefix(key, "ports.") {
			// 嵌套字段
			nestedFilter[key] = value
		} else {
			// 其他字段，默认放到顶级
			topLevelFilter[key] = value
		}
	}
}

// expandNestedFilter 处理嵌套字段查询条件
// 注意：在 $unwind 之后，字段路径保持不变（ports.port 仍然是 ports.port）
func expandNestedFilter(filter bson.M) bson.M {
	expanded := bson.M{}

	for key, value := range filter {
		if key == "$and" || key == "$or" {
			// 递归处理逻辑运算符
			if arr, ok := value.([]interface{}); ok {
				expandedArr := []interface{}{}
				for _, item := range arr {
					if itemMap, ok := item.(bson.M); ok {
						expandedArr = append(expandedArr, expandNestedFilter(itemMap))
					} else if itemMap, ok := item.(map[string]interface{}); ok {
						expandedArr = append(expandedArr, expandNestedFilter(bson.M(itemMap)))
					} else {
						expandedArr = append(expandedArr, item)
					}
				}
				expanded[key] = expandedArr
			}
		} else {
			// 字段路径保持不变，因为 $unwind 后路径仍然有效
			expanded[key] = value
		}
	}

	return expanded
}

func flattenIPAssets(assets []models.IPAsset) []models.IPAssetFlat {
	var rows []models.IPAssetFlat

	for _, asset := range assets {
		id := asset.ID.Hex()
		totalRows := 0
		for _, port := range asset.Ports {
			serverCount := len(port.Server)
			if serverCount == 0 {
				serverCount = 1
			}
			totalRows += serverCount
		}
		if totalRows == 0 {
			totalRows = 1
		}

		if len(asset.Ports) == 0 {
			rows = append(rows, models.IPAssetFlat{
				IP:          asset.IP,
				IPRowSpan:   totalRows,
				PortRowSpan: 1,
				Time:        asset.Time,
			})
			continue
		}

		ipRowAssigned := false
		for _, port := range asset.Ports {
			portRowSpan := len(port.Server)
			if portRowSpan == 0 {
				portRowSpan = 1
			}
			portRowAssigned := false

			servers := port.Server
			if len(servers) == 0 {
				servers = []models.PortServer{{}}
			}

			for _, server := range servers {
				row := models.IPAssetFlat{
					ID:        id,
					IP:        asset.IP,
					Port:      port.Port,
					Domain:    server.Domain,
					Service:   server.Service,
					WebServer: server.WebServer,
					DataKey:   random.GenerateRandomString(6),
					Products:  append([]string{}, server.Technologies...),
					Time:      asset.Time,
				}

				if !ipRowAssigned {
					row.IPRowSpan = totalRows
					ipRowAssigned = true
				} else {
					row.IPRowSpan = 0
				}

				if !portRowAssigned {
					row.PortRowSpan = portRowSpan
					portRowAssigned = true
				} else {
					row.PortRowSpan = 0
				}

				rows = append(rows, row)
			}
		}
	}

	return rows
}
