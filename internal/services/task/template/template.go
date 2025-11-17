package template

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"

	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/task/template"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义模板服务接口
type Service interface {
	List(ctx *gin.Context, pageIndex, pageSize int, query string) (*models.TemplateList, error)
	Detail(ctx *gin.Context, id string) (*models.ScanTemplate, error)
	Save(ctx *gin.Context, id string, result *models.ScanTemplate) error
	Delete(ctx *gin.Context, ids []string) error
}

type service struct {
	repo template.Repository
}

// NewService 创建模板服务实例
func NewService() Service {
	return &service{
		repo: template.NewRepository(),
	}
}

// List 获取模板列表
func (s *service) List(ctx *gin.Context, pageIndex, pageSize int, query string) (*models.TemplateList, error) {
	filter := bson.M{
		"name": bson.M{"$regex": query},
	}

	// 设置查询选项
	opts := options.Find().
		SetProjection(bson.M{
			"_id":  1,
			"name": 1,
		}).
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))

	// 获取总数
	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count templates: %w", err)
	}

	// 获取列表
	list, err := s.repo.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find templates: %w", err)
	}

	return &models.TemplateList{
		List:  list,
		Total: total,
	}, nil
}

// Detail 获取模板详情
func (s *service) Detail(ctx *gin.Context, id string) (*models.ScanTemplate, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	filter := bson.M{"_id": objID}

	result, err := s.repo.FindOne(ctx, filter, options.FindOne())
	if err != nil {
		return nil, fmt.Errorf("failed to find template: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("template not found")
	}

	return result, nil
}

// Save 保存模板
func (s *service) Save(ctx *gin.Context, id string, result *models.ScanTemplate) error {
	if id == "" {
		// 插入新模板
		_, err := s.repo.InsertOne(ctx, result)
		if err != nil {
			return fmt.Errorf("failed to insert template: %w", err)
		}
		return nil
	}

	// 更新现有模板
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": result}

	err = s.repo.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	return nil
}

// Delete 删除模板
func (s *service) Delete(ctx *gin.Context, ids []string) error {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid id format: %w", err)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	err := s.repo.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete templates: %w", err)
	}

	return nil
}
