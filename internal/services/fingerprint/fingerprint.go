package fingerprint

import (
	"context"
	"errors"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	repo "github.com/Autumn-27/ScopeSentry-go/internal/repositories/fingerprint"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/node"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	List(ctx context.Context, search string, pageIndex, pageSize int) (list []models.FingerprintRule, total int64, err error)
	Add(ctx context.Context, data models.FingerprintRule) (string, error)
	Update(ctx context.Context, id string, data models.FingerprintRule) error
	Delete(ctx context.Context, ids []string) (int64, error)
}

type service struct {
	repo        repo.Repository
	nodeService node.Service
}

func NewService() Service {
	return &service{
		repo:        repo.NewRepository(),
		nodeService: node.NewService(),
	}
}

func (s *service) List(ctx context.Context, search string, pageIndex, pageSize int) ([]models.FingerprintRule, int64, error) {
	filter := bson.M{"name": bson.M{"$regex": search, "$options": "i"}}
	// 设置查询选项：投影、分页、排序
	opts := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	list, err := s.repo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *service) Add(ctx context.Context, data models.FingerprintRule) (string, error) {
	if data.Rule == "" {
		return "", errors.New("rule is null")
	}
	express, err := helper.StringToPostfix(data.Rule)
	if err != nil || len(express) == 0 {
		return "", errors.New("rule to express error")
	}
	data.Express = express
	data.Amount = 0
	id, err := s.repo.Insert(ctx, &data)
	if err != nil {
		return "", err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "finger"})
	return id, nil
}

func (s *service) Update(ctx context.Context, id string, data models.FingerprintRule) error {
	if data.Rule == "" {
		return errors.New("rule is null")
	}
	express, err := helper.StringToPostfix(data.Rule)
	if err != nil || len(express) == 0 {
		return errors.New("rule to express error")
	}
	update := bson.M{"$set": bson.M{
		"name":            data.Name,
		"rule":            data.Rule,
		"express":         express,
		"category":        data.Category,
		"parent_category": data.ParentCategory,
		"state":           data.State,
	}}
	if err := s.repo.Update(ctx, id, update); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "finger"})
	return nil
}

func (s *service) Delete(ctx context.Context, ids []string) (int64, error) {
	deleted, err := s.repo.DeleteMany(ctx, ids)
	if err != nil {
		return 0, err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "finger"})
	return deleted, nil
}
