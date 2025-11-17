package statistics

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/statistics"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/asset"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

type Service interface {
	GetStatisticsData(ctx *gin.Context) (*models.StatisticsData, error)
	GetPortStatistics(ctx *gin.Context, req *models.SearchRequest) (*models.PortStatisticsResponse, error)
	GetTitleStatistics(ctx *gin.Context, req *models.TitleStatisticsRequest) (*models.TitleStatisticsResponse, error)
	GetTypeStatistics(ctx *gin.Context, req *models.SearchRequest) (*models.TypeStatisticsResponse, error)
	GetIconStatistics(ctx *gin.Context, req *models.SearchRequest) (*models.IconStatisticsResponse, error)
	GetAppStatistics(ctx *gin.Context, req *models.SearchRequest) (*models.AppStatisticsResponse, error)
}

type service struct {
	repo         statistics.Repository
	assetService asset.Service
}

func NewService() Service {
	return &service{
		repo:         statistics.NewRepository(),
		assetService: asset.NewService(),
	}
}

func (s *service) GetStatisticsData(ctx *gin.Context) (*models.StatisticsData, error) {
	var wg sync.WaitGroup
	var err error

	// 创建变量来存储每个查询的结果
	var assetCount, subdomainCount, sensitiveCount, urlCount, vulnerabilityCount int64

	// 查询并发执行
	wg.Add(5) // 添加五个任务到WaitGroup

	// 查询资产数量
	go func() {
		defer wg.Done()
		assetCount, err = s.repo.CountDocuments(ctx, "asset", bson.M{})
	}()

	// 查询子域数量
	go func() {
		defer wg.Done()
		subdomainCount, err = s.repo.CountDocuments(ctx, "subdomain", bson.M{})
	}()

	// 查询敏感信息数量
	go func() {
		defer wg.Done()
		sensitiveCount, err = s.repo.CountDocuments(ctx, "SensitiveResult", bson.M{})
	}()

	// 查询URL扫描数量
	go func() {
		defer wg.Done()
		urlCount, err = s.repo.CountDocuments(ctx, "UrlScan", bson.M{})
	}()

	// 查询漏洞数量
	go func() {
		defer wg.Done()
		vulnerabilityCount, err = s.repo.CountDocuments(ctx, "vulnerability", bson.M{})
	}()

	// 等待所有 goroutine 完成
	wg.Wait()

	// 如果有错误，返回错误
	if err != nil {
		return nil, err
	}

	// 返回结果
	return &models.StatisticsData{
		AssetCount:         assetCount,
		SubdomainCount:     subdomainCount,
		SensitiveCount:     sensitiveCount,
		UrlCount:           urlCount,
		VulnerabilityCount: vulnerabilityCount,
	}, nil
}

func (s *service) GetPortStatistics(ctx *gin.Context, query *models.SearchRequest) (*models.PortStatisticsResponse, error) {
	query.Index = "asset"
	searchQuery, err := s.assetService.BuildSearchQuery(*query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$group": bson.M{
				"_id":    "$port",
				"number": bson.M{"$sum": 1},
			},
		},
		{
			"$match": bson.M{
				"_id": bson.M{"$ne": nil},
			},
		},
		{
			"$sort": bson.M{
				"number": -1,
			},
		},
		{
			"$limit": 200,
		},
	}

	// 返回 *mongo.Cursor
	cursor, err := s.repo.Aggregate(ctx, "asset", pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 直接将结果写入结构体切片
	var items []models.PortItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return &models.PortStatisticsResponse{
		Port: items,
	}, nil
}

func (s *service) GetTitleStatistics(ctx *gin.Context, req *models.TitleStatisticsRequest) (*models.TitleStatisticsResponse, error) {
	//req.Filter["type"] = []string{"https", "http"}
	//
	//pipeline := []bson.M{
	//	{
	//		"$match": req.Filter,
	//	},
	//	{
	//		"$facet": bson.M{
	//			"by_title": []bson.M{
	//				{"$group": bson.M{"_id": "$title", "num_tutorial": bson.M{"$sum": 1}}},
	//				{"$match": bson.M{"_id": bson.M{"$ne": ""}}},
	//			},
	//		},
	//	},
	//}
	//
	//result, err := s.repo.Aggregate(ctx, "asset", pipeline)
	//if err != nil {
	//	return nil, err
	//}
	//
	//response := &models.TitleStatisticsResponse{
	//	Title: make([]models.TitleItem, 0),
	//}
	//
	//for _, r := range result {
	//	byTitle := r["by_title"].([]interface{})
	//	for _, t := range byTitle {
	//		title := t.(bson.M)
	//		response.Title = append(response.Title, models.TitleItem{
	//			Value:  title["_id"].(string),
	//			Number: title["num_tutorial"].(int64),
	//		})
	//	}
	//}

	return nil, nil
}

func (s *service) GetTypeStatistics(ctx *gin.Context, query *models.SearchRequest) (*models.TypeStatisticsResponse, error) {
	query.Index = "asset"
	searchQuery, err := s.assetService.BuildSearchQuery(*query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)
	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$group": bson.M{
				"_id":    "$service",
				"number": bson.M{"$sum": 1},
			},
		},
		{
			"$match": bson.M{
				"_id": bson.M{"$ne": nil},
			},
		},
		{
			"$sort": bson.M{
				"number": -1,
			},
		},
	}

	// 返回 *mongo.Cursor
	cursor, err := s.repo.Aggregate(ctx, "asset", pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 直接将结果写入结构体切片
	var items []models.TypeItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return &models.TypeStatisticsResponse{
		Service: items,
	}, nil
}

func (s *service) GetIconStatistics(ctx *gin.Context, query *models.SearchRequest) (*models.IconStatisticsResponse, error) {
	skip := (query.PageIndex - 1) * query.PageSize
	searchQuery, err := s.assetService.BuildSearchQuery(*query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)
	pipeline := []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"faviconmmh3": 1,
			"iconcontent": 1,
		}},
		{"$group": bson.M{
			"_id":         "$faviconmmh3",
			"number":      bson.M{"$sum": 1},
			"iconcontent": bson.M{"$first": "$iconcontent"},
		}},
		{"$match": bson.M{
			"$and": []bson.M{
				{"_id": bson.M{"$ne": ""}},
				{"_id": bson.M{"$ne": nil}},
			},
		}},
		{"$sort": bson.M{"number": -1}},
		{"$skip": skip},
		{"$limit": query.PageSize},
	}

	// 返回 *mongo.Cursor
	cursor, err := s.repo.Aggregate(ctx, "asset", pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 直接将结果写入结构体切片
	var items []models.IconItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return &models.IconStatisticsResponse{
		Icon: items,
	}, nil
}

func (s *service) GetAppStatistics(ctx *gin.Context, query *models.SearchRequest) (*models.AppStatisticsResponse, error) {
	query.Index = "asset"
	searchQuery, err := s.assetService.BuildSearchQuery(*query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	// 定义管道（只处理 technologies）
	pipeline := []bson.M{
		{"$match": filter}, // 基于过滤条件匹配数据
		{"$project": bson.M{
			"technologies": 1, // 只保留必要的字段
		}},
		{"$unwind": "$technologies"}, // 展开 technologies 数组
		{"$group": bson.M{
			"_id":    "$technologies",   // 按技术栈分组
			"number": bson.M{"$sum": 1}, // 统计数量
		}},
		{"$sort": bson.M{"number": -1}}, // 按数量倒序排序
		{
			"$limit": 200,
		},
	}

	// 执行聚合查询
	cursor, err := s.repo.Aggregate(ctx, "asset", pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 直接将结果写入结构体切片
	var items []models.AppItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	// 返回统计结果
	return &models.AppStatisticsResponse{
		Product: items,
	}, nil
}
