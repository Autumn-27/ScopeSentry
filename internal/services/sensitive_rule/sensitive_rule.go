// sensitive_rule-------------------------------------
// @file      : sensitive_rule.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/28 22:28
// -------------------------------------------

package sensitive_rule

import (
	"errors"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/sensitive_rule"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	RuleList(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.SensitiveRuleListItem, int64, error)
	RuleUpdate(ctx *gin.Context, id string, name, regular, color string, state bool) error
	RuleAdd(ctx *gin.Context, name, regular, color string, state bool) error
	RuleUpdateState(ctx *gin.Context, ids []string, state bool) (int64, error)
	RuleDelete(ctx *gin.Context, ids []string) (int64, error)
}

type service struct {
	repo        sensitive_rule.Repository
	nodeService node.Service
}

func NewService() Service {
	return &service{
		repo:        sensitive_rule.NewRepository(),
		nodeService: node.NewService(),
	}
}

func (s *service) RuleList(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.SensitiveRuleListItem, int64, error) {
	filter := bson.M{"name": bson.M{"$regex": search, "$options": "i"}}
	total, err := s.repo.RuleCount(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(int64((pageIndex - 1) * pageSize)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "_id", Value: -1}})
	items, err := s.repo.RuleList(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	list := make([]models.SensitiveRuleListItem, 0, len(items))
	for _, it := range items {
		list = append(list, models.SensitiveRuleListItem{
			ID: it.ID.Hex(), Name: it.Name, Regular: it.Regular, Color: it.Color, State: it.State,
		})
	}
	return list, total, nil
}

func (s *service) RuleUpdate(ctx *gin.Context, id string, name, regular, color string, state bool) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"name": name, "regular": regular, "color": color, "state": state}}
	if err := s.repo.RuleUpdate(ctx, oid, update); err != nil {
		return err
	}
	// 广播刷新
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "sensitive"})
	return nil
}

func (s *service) RuleAdd(ctx *gin.Context, name, regular, color string, state bool) error {
	if regular == "" {
		return errors.New("regular is null")
	}
	item := &models.SensitiveRuleItem{Name: name, Regular: regular, Color: color, State: state}
	if _, err := s.repo.RuleInsert(ctx, item); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "sensitive"})
	return nil
}

func (s *service) RuleUpdateState(ctx *gin.Context, ids []string, state bool) (int64, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		objIDs = append(objIDs, oid)
	}
	n, err := s.repo.RuleUpdateStateMany(ctx, objIDs, state)
	if err != nil {
		return 0, err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "sensitive"})
	return n, nil
}

func (s *service) RuleDelete(ctx *gin.Context, ids []string) (int64, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		objIDs = append(objIDs, oid)
	}
	n, err := s.repo.RuleDeleteMany(ctx, objIDs)
	if err != nil {
		return 0, err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "sensitive"})
	return n, nil
}
