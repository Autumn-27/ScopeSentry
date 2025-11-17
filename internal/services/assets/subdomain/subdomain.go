package subdomain

import (
	"context"
	"strings"

	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/gin-gonic/gin"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/subdomain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义子域名服务接口
type Service interface {
	GetSubdomains(ctx context.Context, query models.SearchRequest) ([]models.Subdomain, error)
	GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error)
	GetTaskTargetByIDs(ctx *gin.Context, ids []string) (string, error)
}

type service struct {
	repo subdomain.Repository
}

// NewService 创建子域名服务实例
func NewService() Service {
	return &service{
		repo: subdomain.NewRepository(),
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
		"host": 1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection)

	subDomains, err := s.repo.FindWithPagination(ctx.Request.Context(), filter, opts)
	if err != nil {
		return "", err
	}

	// 将 Host 拼接成字符串
	var builder strings.Builder
	for i, subDomain := range subDomains {
		builder.WriteString(subDomain.Host)
		if i != len(subDomains)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}

// GetSubdomains 获取子域名列表
func (s *service) GetSubdomains(ctx context.Context, query models.SearchRequest) ([]models.Subdomain, error) {
	// 构建查询条件
	query.Index = "subdomain"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	projection := bson.M{
		"_id":   1,
		"host":  1,
		"type":  1,
		"value": 1,
		"ip":    1,
		"time":  1,
		"tags":  1,
	}
	// 分页查询
	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.D{{Key: "time", Value: -1}})

	return s.repo.FindWithPagination(ctx, filter, opts)
}

func (s *service) GetTaskTarget(ctx *gin.Context, query models.SearchRequest) (string, error) {
	query.Index = "subdomain"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return "", err
	}
	filter := bson.M(searchQuery)
	projection := bson.M{
		"host": 1,
	}

	// 执行查询
	opts := options.Find().
		SetProjection(projection)

	if query.PageSize != 0 {
		opts = opts.SetLimit(int64(query.PageSize))
	}

	subDomains, err := s.repo.FindWithPagination(ctx.Request.Context(), filter, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i, rootDomain := range subDomains {
		builder.WriteString(rootDomain.Host)
		if i != len(subDomains)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}
