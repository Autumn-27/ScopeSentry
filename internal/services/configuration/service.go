// configuration-------------------------------------
// @file      : service.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/29 21:20
// -------------------------------------------

package configuration

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	assetCommon "github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	nservice "github.com/Autumn-27/ScopeSentry/internal/services/node"

	"github.com/Autumn-27/ScopeSentry/internal/repositories/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	repo         common.Repository
	nodeService  nservice.Service
	dedupService assetCommon.DedupService
}

func NewService() *Service {
	return &Service{
		repo:         common.NewRepository(),
		nodeService:  nservice.NewService(),
		dedupService: assetCommon.NewDedupService(),
	}
}

const (
	collConfig       = "config"
	collNotification = "notification"
)

// GetSubfinderContent 读取 SubfinderApiConfig
func (s *Service) GetSubfinderContent(ctx *gin.Context) (string, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "SubfinderApiConfig"}, bson.M{"_id": 0})
	if err != nil {
		return "", err
	}
	if v, ok := doc["value"].(string); ok {
		return v, nil
	}
	return "", nil
}

// SaveSubfinderContent 保存 SubfinderApiConfig 并通知节点刷新
func (s *Service) SaveSubfinderContent(ctx *gin.Context, content string) error {
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "SubfinderApiConfig"}, bson.M{"name": "SubfinderApiConfig", "value": content}); err != nil {
		return err
	}
	// 通知所有节点刷新 subfinder 配置
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "subfinder"})
	return nil
}

// GetRadContent 读取 RadConfig
func (s *Service) GetRadContent(ctx *gin.Context) (string, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "RadConfig"}, bson.M{"_id": 0})
	if err != nil {
		return "", err
	}
	if v, ok := doc["value"].(string); ok {
		return v, nil
	}
	return "", nil
}

// SaveRadContent 保存 RadConfig 并通知节点刷新
func (s *Service) SaveRadContent(ctx *gin.Context, content string) error {
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "RadConfig"}, bson.M{"name": "RadConfig", "value": content}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "rad"})
	return nil
}

// GetSystemData 返回 type=system 的所有键值
func (s *Service) GetSystemData(ctx *gin.Context) (map[string]interface{}, error) {
	list, err := s.repo.FindMany(ctx.Request.Context(), collConfig, bson.M{"type": "system"}, nil)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	for _, m := range list {
		name, _ := m["name"].(string)
		if name == "" {
			continue
		}
		out[name] = m["value"]
	}
	return out, nil
}

// SaveSystemData 保存 system 配置，并通知节点
func (s *Service) SaveSystemData(ctx *gin.Context, kv map[string]interface{}) error {
	for k, v := range kv {
		doc := bson.M{"type": "system", "name": k, "value": v}
		if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"type": "system", "name": k}, doc); err != nil {
			return err
		}
	}
	timezone, _ := kv["timezone"].(string)
	modulesConfig := fmt.Sprintf("%v", kv["ModulesConfig"]) // 兼容字符串或其他类型
	msg := timezone + "[*]" + modulesConfig
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "system", Content: msg})
	return nil
}

// GetDeduplicationConfig 查询去重配置与下一次运行时间（若有）
func (s *Service) GetDeduplicationConfig(ctx *gin.Context) (map[string]interface{}, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "deduplication"}, bson.M{"_id": 0})
	if err != nil {
		// 若不存在，返回空 map
		return map[string]interface{}{}, nil
	}
	return doc, nil
}

// SaveDeduplicationConfig 保存去重配置；如需要立即运行，透传到节点或调度器（此处仅保存与刷新）
func (s *Service) SaveDeduplicationConfig(ctx *gin.Context, cfg models.DepConfig) error {
	cfg.Name = "deduplication"
	runNow := cfg.RunNow
	cfg.RunNow = false
	if err := s.repo.Upsert(ctx.Request.Context(), collConfig, bson.M{"name": "deduplication"}, cfg); err != nil {
		return err
	}
	// 这里按需可扩展调用内部调度器或消息下发；当前版本仅保存
	if runNow {
		go func() {
			err := s.dedupService.DoAssetDeduplication()
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}()
	}
	return nil
}

// GetNotificationList 获取通知列表
func (s *Service) GetNotificationList(ctx *gin.Context) ([]models.Notification, error) {
	var notifications []models.Notification
	err := s.repo.Find(ctx.Request.Context(), collNotification, bson.M{}, nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// AddNotification 新增通知并刷新节点
func (s *Service) AddNotification(ctx *gin.Context, data models.Notification) error {
	_, err := s.repo.InsertOne(ctx.Request.Context(), collNotification, data)
	if err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}

// UpdateNotification 更新通知并刷新节点
func (s *Service) UpdateNotification(ctx *gin.Context, data models.UpdateNotification) error {
	if data.ID == "" {
		return fmt.Errorf("missing id")
	}
	objID, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}
	if err := s.repo.UpdateOne(ctx.Request.Context(), collNotification, bson.M{"_id": objID}, bson.M{"$set": data}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}

// DeleteNotifications 批量删除通知并刷新节点
func (s *Service) DeleteNotifications(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// 将 string 转换为 ObjectID
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid ObjectID: %s, error: %w", id, err)
		}
		objectIDs = append(objectIDs, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	if _, err := s.repo.DeleteMany(ctx.Request.Context(), collNotification, filter); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}

// GetNotificationConfig 读取通知配置
func (s *Service) GetNotificationConfig(ctx *gin.Context) (map[string]interface{}, error) {
	doc, err := s.repo.FindOne(ctx.Request.Context(), collConfig, bson.M{"name": "notification"}, bson.M{})
	if err != nil {
		return map[string]interface{}{}, nil
	}
	delete(doc, "_id")
	delete(doc, "type")
	delete(doc, "name")
	return doc, nil
}

// UpdateNotificationConfig 更新通知配置并刷新节点
func (s *Service) UpdateNotificationConfig(ctx *gin.Context, data map[string]interface{}) error {
	if data == nil {
		return nil
	}
	if err := s.repo.UpdateOne(ctx.Request.Context(), collConfig, bson.M{"name": "notification"}, bson.M{"$set": data}); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "notification"})
	return nil
}
