package url

import (
	"context"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/assets/url"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetURLs(ctx context.Context, query models.SearchRequest) ([]models.URL, error)
	GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error)
	GetTaskTargetByIDs(ctx *gin.Context, ids []string) (string, error)
}

type service struct {
	repo url.Repository
}

func NewService() Service {
	return &service{
		repo: url.NewRepository(),
	}
}

func (s *service) GetTaskTargetByIDs(ctx *gin.Context, ids []string) (string, error) {
	if len(ids) == 0 {
		return "", nil
	}

	// 将字符串ID转换为ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", err
		}
		objectIDs = append(objectIDs, objectID)
	}

	// 构建查询条件
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	projection := bson.M{
		"output": 1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection)

	urlResults, err := s.repo.Find(ctx.Request.Context(), filter, opts)
	if err != nil {
		return "", err
	}

	// 将 Host 拼接成字符串
	var builder strings.Builder
	for i, result := range urlResults {
		builder.WriteString(result.Output)
		if i != len(urlResults)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}

func (s *service) GetURLs(ctx context.Context, query models.SearchRequest) ([]models.URL, error) {
	// 构建查询条件
	query.Index = "UrlScan"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	projection := bson.M{
		"_id":        1,
		"input":      1,
		"source":     1,
		"status":     1,
		"length":     1,
		"outputtype": 1,
		"output":     1,
		"time":       1,
		"tags":       1,
	}
	// 构建排序
	sortBy := bson.D{{Key: "_id", Value: -1}}
	if val, ok := query.Sort["length"]; ok {
		dir := -1
		if val == "ascending" {
			dir = 1
		}
		sortBy = bson.D{{Key: "length", Value: dir}}
	}
	// 分页查询
	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize)).
		SetSort(sortBy)

	return s.repo.Find(ctx, filter, opts)
}

func (s *service) GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error) {
	query.Index = "UrlScan"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return "", err
	}
	filter := bson.M(searchQuery)
	projection := bson.M{
		"output": 1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection)

	if query.PageSize != 0 {
		opts = opts.SetLimit(int64(query.PageSize))
	}

	urlResults, err := s.repo.Find(ctx.Request.Context(), filter, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i, result := range urlResults {
		builder.WriteString(result.Output)
		if i != len(urlResults)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}
