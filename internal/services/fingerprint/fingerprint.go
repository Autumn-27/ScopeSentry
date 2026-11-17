package fingerprint

import (
	"context"
	"errors"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	commonRepo "github.com/Autumn-27/ScopeSentry/internal/repositories/common"
	repo "github.com/Autumn-27/ScopeSentry/internal/repositories/fingerprint"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/utils/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type BatchAddItem struct {
	ID         string
	Name       string
	Content    string
	UpdateTime string
}

type BatchAddResult struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type Service interface {
	List(ctx context.Context, search string, pageIndex, pageSize int) (list []models.FingerprintRule, total int64, err error)
	Add(ctx context.Context, data models.FingerprintRule) (string, error)
	Update(ctx context.Context, id string, data models.FingerprintRule) error
	Delete(ctx context.Context, ids []string) (int64, error)
	GetVersion(ctx context.Context) (string, error)
	BatchAdd(ctx context.Context, items []BatchAddItem) (results []BatchAddResult, successCount, failedCount int, err error)
}

type service struct {
	repo        repo.Repository
	commonRepo  commonRepo.Repository
	nodeService node.Service
}

func NewService() Service {
	return &service{
		repo:        repo.NewRepository(),
		commonRepo:  commonRepo.NewRepository(),
		nodeService: node.NewService(),
	}
}

func (s *service) List(ctx context.Context, search string, pageIndex, pageSize int) ([]models.FingerprintRule, int64, error) {
	filter := bson.M{"name": bson.M{"$regex": search, "$options": "i"}}
	// 设置查询选项：投影、分页、排序
	opts := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "_id", Value: -1}})

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
	// Rule 现在存储的是 YAML 内容，不再转换为 Express
	// express, err := helper.StringToPostfix(data.Rule)
	// if err != nil || len(express) == 0 {
	// 	return "", errors.New("rule to express error")
	// }
	// data.Express = express
	data.Amount = 0
	if data.FingerprintId == "" {
		data.FingerprintId = "SAA_" + random.GenerateRandomString(12)
	}
	// 同名已存在 → 直接更新那条，避免导入同一份 YAML 又生成一条重复指纹
	if existing := s.lookupHexID(ctx, bson.M{"name": data.Name}); shouldReuseFingerprintID(existing) {
		if err := s.repo.Update(ctx, existing, bson.M{"$set": fingerprintSetDoc(data)}); err != nil {
			return "", err
		}
		_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "finger"})
		return existing, nil
	}
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
	// Rule 现在存储的是 YAML 内容，不再转换为 Express
	// express, err := helper.StringToPostfix(data.Rule)
	// if err != nil || len(express) == 0 {
	// 	return errors.New("rule to express error")
	// }
	update := bson.M{"$set": bson.M{
		"name":            data.Name,
		"rule":            data.Rule,
		"category":        data.Category,
		"parent_category": data.ParentCategory,
		"state":           data.State,
		"company":         data.Company,
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

// GetVersion 查询指纹版本
func (s *service) GetVersion(ctx context.Context) (string, error) {
	doc, err := s.commonRepo.FindOne(ctx, "config", bson.M{"name": "FingerVersion"}, bson.M{"_id": 0})
	if err != nil {
		return "", err
	}
	if v, ok := doc["value"].(string); ok {
		return v, nil
	}
	return "", nil
}

// BatchAdd 批量新增指纹规则
func (s *service) BatchAdd(ctx context.Context, items []BatchAddItem) ([]BatchAddResult, int, int, error) {
	results := make([]BatchAddResult, 0, len(items))
	successCount := 0
	failedCount := 0
	updateTimes := make([]string, 0) // 收集所有更新时间

	type FingerprintYaml struct {
		Fingerprint models.FingerprintRule `yaml:"fingerprint"`
	}

	for _, item := range items {
		result := BatchAddResult{
			ID:      item.ID,
			Success: false,
		}

		// 解析 YAML content
		var fp FingerprintYaml
		if err := yaml.Unmarshal([]byte(item.Content), &fp); err != nil {
			result.Error = fmt.Sprintf("YAML解析失败: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}

		// 构建 FingerprintRule
		data := models.FingerprintRule{
			Name:           fp.Fingerprint.Name,
			Rule:           item.Content,
			FingerprintId:  item.ID, // 使用传入的 ID 作为 FingerprintId
			Category:       fp.Fingerprint.Category,
			ParentCategory: fp.Fingerprint.ParentCategory,
			State:          true, // 默认启用
			Company:        fp.Fingerprint.Company,
			Amount:         0,
		}

		// 先判重：同 fingerprint_id 命中则更新；否则看有没有同名旧文档（历史数据
		// 里 Add 生成的随机 ID 与上游 ID 不同），命中就更新旧文档并回填 ID，
		// 避免"一键更新后每个指纹重复两份"。
		byID := s.lookupHexID(ctx, bson.M{"fingerprint_id": item.ID})
		byName := ""
		if byID == "" {
			byName = s.lookupHexID(ctx, bson.M{"name": data.Name})
		}
		match := decideFingerprintMatch(byID, byName)

		var err error
		if match.MatchBy == "_id" {
			err = s.repo.Update(ctx, match.ExistingHexID, bson.M{"$set": fingerprintSetDoc(data)})
		} else {
			err = s.commonRepo.Upsert(ctx, "FingerprintRules", bson.M{"fingerprint_id": item.ID}, data)
		}
		if err != nil {
			result.Error = fmt.Sprintf("更新或插入失败: %v", err)
			results = append(results, result)
			failedCount++
			continue
		}

		result.Success = true
		results = append(results, result)
		successCount++

		// 收集更新时间
		if item.UpdateTime != "" {
			updateTimes = append(updateTimes, item.UpdateTime)
		}
	}

	// 批量操作完成后统一刷新配置
	if successCount > 0 {
		_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "finger"})
	}

	// 获取最新的更新时间并更新配置版本
	if len(updateTimes) > 0 {
		latestUpdateTime := updateTimes[0]
		// 找到最新的时间（时间格式 "YYYY-MM-DD HH:MM" 可以直接字符串比较）
		for _, updateTime := range updateTimes {
			if updateTime > latestUpdateTime {
				latestUpdateTime = updateTime
			}
		}

		// 更新 FingerVersion 配置
		configFilter := bson.M{"name": "FingerVersion"}
		configDoc := bson.M{
			"name":  "FingerVersion",
			"value": latestUpdateTime,
			"type":  "finger",
		}
		if err := s.commonRepo.Upsert(ctx, "config", configFilter, configDoc); err != nil {
			// 记录错误但不影响整体结果
			_ = err
		}
	}

	return results, successCount, failedCount, nil
}
