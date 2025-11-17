package page_monitoring

import (
	"context"
	"errors"
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/page_monitoring"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrPageNotFound = errors.New("page not found")
	ErrInvalidData  = errors.New("invalid page monitoring data")
)

type Service interface {
	GetHistory(ctx context.Context, id string) ([]*models.PageMonitoring, error)
	GetContent(ctx context.Context, id string, flag string) (*models.PageMonitoringBody, error)
	GetDiff(ctx context.Context, id string) (*models.PageMonitoringBody, error)
	GetResult(ctx *gin.Context, req models.SearchRequest) (*models.PageMonitoringResultResponse, error)
}

type service struct {
	pageRepo page_monitoring.Repository
}

func NewService() Service {
	return &service{
		pageRepo: page_monitoring.NewRepository(),
	}
}

func (s *service) GetHistory(ctx context.Context, id string) ([]*models.PageMonitoring, error) {
	filter := bson.M{"_id": id}
	opts := options.Find().
		SetSort(bson.D{{"created_at", -1}})
	return s.pageRepo.FindHistory(ctx, filter, opts)
}

func (s *service) GetContent(ctx context.Context, id string, flag string) (*models.PageMonitoringBody, error) {
	filter := bson.M{"_id": id}
	opts := options.FindOne()
	return s.pageRepo.FindContent(ctx, filter, opts)
}

func (s *service) GetDiff(ctx context.Context, id string) (*models.PageMonitoringBody, error) {
	filter := bson.M{"md5": id}
	opts := options.FindOne()
	return s.pageRepo.FindDiff(ctx, filter, opts)
}

// GetResult 获取页面监控结果
func (s *service) GetResult(ctx *gin.Context, req models.SearchRequest) (*models.PageMonitoringResultResponse, error) {
	// 构建查询条件
	req.Index = "PageMonitoring"
	searchQuery, err := helper.GetSearchQuery(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build search query: %w", err)
	}
	searchQuery["hash"] = map[string]interface{}{"$size": 2}
	filter := bson.M(searchQuery)

	// 设置查询选项
	opts := options.Find().
		SetSort(bson.D{{"time", -1}}).
		SetSkip(int64((req.PageIndex - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	// 获取总数
	total, err := s.pageRepo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// 获取结果列表
	results, err := s.pageRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	return &models.PageMonitoringResultResponse{
		List:  results,
		Total: total,
	}, nil
}
