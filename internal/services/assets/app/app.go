package app

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/bootstrap"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/assets/app"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetAppData(ctx *gin.Context, query models.SearchRequest) (*models.AppResponse, error)
}

type service struct {
	repo app.Repository
}

func NewService() Service {
	return &service{
		repo: app.NewRepository(),
	}
}

func (s *service) GetAppData(ctx *gin.Context, query models.SearchRequest) (*models.AppResponse, error) {
	// 构建查询条件
	query.Index = "app"
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
		return &models.AppResponse{
			List:  []models.AppResult{},
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
	for i := range results {
		if projectName, exists := bootstrap.ProjectList[results[i].Project]; exists {
			results[i].Project = projectName
		} else {
			results[i].Project = ""
		}
	}
	return &models.AppResponse{
		List:  results,
		Total: totalCount,
	}, nil
}
