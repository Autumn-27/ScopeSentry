package crawler

import (
	"context"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/assets/crawler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetCrawlers(ctx context.Context, query models.SearchRequest) ([]models.CrawlerResult, error)
}

type service struct {
	repo crawler.Repository
}

func NewService() Service {
	return &service{
		repo: crawler.NewRepository(),
	}
}

func (s *service) GetCrawlers(ctx context.Context, query models.SearchRequest) ([]models.CrawlerResult, error) {
	// 构建查询条件
	query.Index = "crawler"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	projection := bson.M{
		"_id":    1,
		"method": 1,
		"body":   1,
		"url":    1,
		"tags":   1,
		"time":   1,
	}
	// 分页查询
	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.D{{Key: "time", Value: -1}})

	return s.repo.Find(ctx, filter, opts)
}
