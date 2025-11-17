package mp

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/bootstrap"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/mp"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetMPData(ctx *gin.Context, query models.SearchRequest) (*models.MPResponse, error)
}

type service struct {
	repo mp.Repository
}

func NewService() Service {
	return &service{
		repo: mp.NewRepository(),
	}
}

func (s *service) GetMPData(ctx *gin.Context, query models.SearchRequest) (*models.MPResponse, error) {
	// 构建查询条件
	query.Index = "mp"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, fmt.Errorf("failed to build search query: %w", err)
	}
	filter := bson.M(searchQuery)

	// 获取总数
	totalCount, err := s.repo.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	if totalCount == 0 {
		return &models.MPResponse{
			List:  []models.MPResult{},
			Total: 0,
		}, nil
	}

	// 设置分页和排序选项
	opts := options.Find().
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.D{{"time", -1}})

	// 查询数据
	results, err := s.repo.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// 处理项目名称映射
	for i := range results {
		if projectName, exists := bootstrap.ProjectList[results[i].Project]; exists {
			results[i].Project = projectName
		} else {
			results[i].Project = ""
		}
	}

	return &models.MPResponse{
		List:  results,
		Total: totalCount,
	}, nil
} 