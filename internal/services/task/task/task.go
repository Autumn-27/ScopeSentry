// Package task -----------------------------
// @file      : task.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/4 22:04
// -------------------------------------------
package task

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/Autumn-27/ScopeSentry-go/internal/logger"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/task/task"
	schedulerCore "github.com/Autumn-27/ScopeSentry-go/internal/scheduler"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/node"
	commonTask "github.com/Autumn-27/ScopeSentry-go/internal/services/task/common"
	schedulerSvc "github.com/Autumn-27/ScopeSentry-go/internal/services/task/scheduler"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义任务服务接口
type Service interface {
	List(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.Task, int64, error)
	CheckTaskNameExists(ctx *gin.Context, name string) (bool, error)
	GetTaskDetail(ctx *gin.Context, taskID string) (*models.Task, error)
	DeleteTasks(ctx *gin.Context, ids []string, delA bool) error
	RetestTask(ctx context.Context, id string) error
	StopTasks(ctx *gin.Context, ids []string) error
	StartTasks(ctx *gin.Context, ids []string) error
	GetTaskProgress(ctx *gin.Context, taskID string, pageIndex, pageSize int) (map[string]interface{}, error)
	TaskProgressNumber(ctx *gin.Context) error
	TaskProgress(ctx *gin.Context) error
	ProcessTaskProgress(ctx *gin.Context, task models.Task) error
}

// service 实现Service接口
type service struct {
	taskRepo         task.Repository
	commonService    commonTask.Service
	schedulerService schedulerSvc.Service
	nodeService      node.Service
}

// NewService 创建任务服务实例
func NewService() Service {
	return &service{
		taskRepo:         task.NewRepository(),
		commonService:    commonTask.NewService(),
		schedulerService: schedulerSvc.NewService(),
		nodeService:      node.NewService(),
	}
}

// List 获取任务列表
func (s *service) List(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.Task, int64, error) {
	filter := bson.M{}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	// 计算总数
	total, err := s.taskRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "creatTime", Value: -1}})

	tasks, err := s.taskRepo.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// CheckTaskNameExists 检查任务名是否存在
func (s *service) CheckTaskNameExists(ctx *gin.Context, name string) (bool, error) {
	return s.taskRepo.CheckTaskNameExists(ctx, name)
}

// GetTaskDetail 获取任务详情
func (s *service) GetTaskDetail(ctx *gin.Context, taskID string) (*models.Task, error) {
	// 将字符串ID转换为ObjectID
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	// 查询任务详情
	task, err := s.taskRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTasks 批量删除任务
func (s *service) DeleteTasks(ctx *gin.Context, ids []string, delA bool) error {
	// 1. 移除调度任务与配置
	// 先尝试从自定义调度器中移除任务
	for _, id := range ids {
		if id == "" {
			continue
		}
		// 移除 gocron 调度任务与 scheduled 集合记录
		_ = schedulerCore.GetGlobalScheduler().RemoveJob(id)

		// 清理 Redis 缓存
		_ = s.taskRepo.ClearTaskCache(ctx, id)
	}

	// 删除 ScheduledTasks 集合中的配置
	err := s.schedulerService.DeleteScheduledTasks(ctx, ids)
	if err != nil {
		logger.Error(fmt.Sprintf("error while deleting tasks: %v", err))
	}

	// 2. 可选：删除资产（依据 taskName）
	if delA {
		// 获取任务名称列表
		var objIDs []primitive.ObjectID
		for _, id := range ids {
			if oid, err := primitive.ObjectIDFromHex(id); err == nil {
				objIDs = append(objIDs, oid)
			}
		}
		if len(objIDs) > 0 {
			// 仅查询 name 字段
			opts := options.Find().SetProjection(bson.M{"name": 1})
			tasks, err := s.taskRepo.Find(ctx, bson.M{"_id": bson.M{"$in": objIDs}}, opts)
			if err == nil {
				nameSet := make(map[string]struct{})
				for _, t := range tasks {
					if t.Name != "" {
						nameSet[t.Name] = struct{}{}
					}
				}
				if len(nameSet) > 0 {
					var names []string
					for n := range nameSet {
						names = append(names, n)
					}
					// 广播到所有节点进行后续清理（兼容原 refresh_config 行为）
					_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "delete_task", Content: strings.Join(ids, ",")})

					// 按 taskName 删除各资产集合数据
					// 注意：部分集合中 taskName 为数组字段，$in 也能匹配
					collections := []string{"asset", "DirScanResult", "SensitiveResult", "SubdomainTakerResult", "UrlScan", "crawler", "subdomain", "vulnerability", "PageMonitoring", "app", "RootDomain", "mp"}
					filter := bson.M{"taskName": bson.M{"$in": names}}
					// 直接通过底层数据库执行批删
					db := mongodb.DB
					for _, coll := range collections {
						_, _ = db.Collection(coll).DeleteMany(ctx.Request.Context(), filter)
					}
				}
			}
		}
	} else {
		// 仍发出节点配置刷新，保持与 Python 端一致
		_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "delete_task", Content: strings.Join(ids, ",")})
	}

	// 3. 删除 task 集合中的任务
	var objIDs2 []primitive.ObjectID
	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs2 = append(objIDs2, oid)
		}
	}
	if len(objIDs2) == 0 {
		return nil
	}
	if _, err := s.taskRepo.DeleteMany(ctx, objIDs2); err != nil {
		return err
	}
	return nil
}

// RetestTask 重新测试任务
func (s *service) RetestTask(ctx context.Context, id string) error {
	// 获取任务
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	t, err := s.taskRepo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("task not found")
	}

	// 创建扫描任务
	if _, err := s.commonService.CreateTaskScan(ctx, *t, id, false); err != nil {
		return err
	}

	// 更新任务状态
	update := bson.M{"$set": bson.M{
		"progress":  0,
		"creatTime": helper.GetNowTimeString(),
		"endTime":   "",
		"status":    1,
	}}
	if err := s.taskRepo.UpdateOne(ctx, bson.M{"_id": objectID}, update); err != nil {
		return err
	}
	return nil
}

// StopTasks 停止任务：广播停任务并将状态置为2
func (s *service) StopTasks(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		// 广播停止任务
		_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "stop_task", Content: id})
		// 状态置为2
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			_ = s.taskRepo.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": 2}})
		}
	}
	return nil
}

// StartTasks 开始任务：从暂停恢复继续扫描
func (s *service) StartTasks(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		t, err := s.taskRepo.FindByID(ctx, oid)
		if err != nil || t == nil {
			continue
		}
		if t.Progress == 100 {
			// 已完成，不再启动
			continue
		}
		t.IsStart = true
		if _, err := s.commonService.CreateTaskScan(ctx, *t, id, true); err != nil {
			continue
		}
		_ = s.taskRepo.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": 1}})
	}
	return nil
}

// GetTaskProgress 获取任务进度信息
func (s *service) GetTaskProgress(ctx *gin.Context, taskID string, pageIndex, pageSize int) (map[string]interface{}, error) {
	// 1. 获取任务基本信息
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := s.taskRepo.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	// 2. 处理目标列表和分页
	targets, err := s.processTargetsWithPagination(ctx, task.Target, task.Ignore, pageIndex, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to process targets: %w", err)
	}

	// 3. 从Redis获取进度数据
	progressData, err := s.getProgressFromRedis(ctx, taskID, targets)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress data: %w", err)
	}

	return map[string]interface{}{
		"list":  progressData,
		"total": task.TaskNum,
	}, nil
}

// processTargetsWithPagination 处理目标列表并分页
func (s *service) processTargetsWithPagination(ctx *gin.Context, targetStr, ignoreStr string, pageIndex, pageSize int) ([]string, error) {
	if targetStr == "" {
		return []string{}, nil
	}

	// 分割目标字符串
	targetLines := strings.Split(targetStr, "\n")
	var allTargets []string

	// 处理每个目标行
	for _, line := range targetLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 使用 helper 获取目标列表（处理忽略规则）
		targets, err := helper.GetTargetList(line, ignoreStr)
		if err != nil {
			continue
		}
		allTargets = append(allTargets, targets...)
	}

	// 分页处理
	startIndex := (pageIndex - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(allTargets) {
		return []string{}, nil
	}

	if endIndex > len(allTargets) {
		endIndex = len(allTargets)
	}

	return allTargets[startIndex:endIndex], nil
}

// getProgressFromRedis 从Redis获取进度数据
func (s *service) getProgressFromRedis(ctx *gin.Context, taskID string, targets []string) ([]models.TaskProgressDetail, error) {
	var progressList []models.TaskProgressDetail

	for _, target := range targets {
		progress := s.createProgressResult(target)

		// 从Redis获取该目标的进度信息
		redisKey := fmt.Sprintf("TaskInfo:progress:%s:%s", taskID, target)
		progressData, err := s.taskRepo.GetProgressFromRedis(ctx, redisKey)
		if err == nil && progressData != nil {
			s.populateProgressFromRedis(progress, progressData)
		}

		progressList = append(progressList, *progress)
	}

	return progressList, nil
}

// createProgressResult 创建进度结果结构
func (s *service) createProgressResult(target string) *models.TaskProgressDetail {
	return &models.TaskProgressDetail{
		ID:                  primitive.NewObjectID(),
		TargetHandler:       []string{"", ""},
		SubdomainScan:       []string{"", ""},
		SubdomainSecurity:   []string{"", ""},
		PortScanPreparation: []string{"", ""},
		PortScan:            []string{"", ""},
		PortFingerprint:     []string{"", ""},
		AssetMapping:        []string{"", ""},
		AssetHandle:         []string{"", ""},
		URLScan:             []string{"", ""},
		WebCrawler:          []string{"", ""},
		URLSecurity:         []string{"", ""},
		DirScan:             []string{"", ""},
		VulnerabilityScan:   []string{"", ""},
		All:                 []string{"", ""},
		Target:              target,
		Node:                "",
	}
}

// populateProgressFromRedis 从Redis数据填充进度信息
func (s *service) populateProgressFromRedis(progress *models.TaskProgressDetail, redisData map[string]string) {
	// 映射Redis字段到进度结构
	fieldMappings := map[string]*[]string{
		"TargetHandler_start":       &progress.TargetHandler,
		"SubdomainScan_start":       &progress.SubdomainScan,
		"SubdomainSecurity_start":   &progress.SubdomainSecurity,
		"PortScanPreparation_start": &progress.PortScanPreparation,
		"PortScan_start":            &progress.PortScan,
		"PortFingerprint_start":     &progress.PortFingerprint,
		"AssetMapping_start":        &progress.AssetMapping,
		"AssetHandle_start":         &progress.AssetHandle,
		"URLScan_start":             &progress.URLScan,
		"WebCrawler_start":          &progress.WebCrawler,
		"URLSecurity_start":         &progress.URLSecurity,
		"DirScan_start":             &progress.DirScan,
		"VulnerabilityScan_start":   &progress.VulnerabilityScan,
		"scan_start":                &progress.All,
	}

	// 填充开始时间
	for redisField, progressField := range fieldMappings {
		if startTime, exists := redisData[redisField]; exists {
			(*progressField)[0] = startTime
		}
	}

	// 填充结束时间
	endFieldMappings := map[string]*[]string{
		"TargetHandler_end":       &progress.TargetHandler,
		"SubdomainScan_end":       &progress.SubdomainScan,
		"SubdomainSecurity_end":   &progress.SubdomainSecurity,
		"PortScanPreparation_end": &progress.PortScanPreparation,
		"PortScan_end":            &progress.PortScan,
		"PortFingerprint_end":     &progress.PortFingerprint,
		"AssetMapping_end":        &progress.AssetMapping,
		"AssetHandle_end":         &progress.AssetHandle,
		"URLScan_end":             &progress.URLScan,
		"WebCrawler_end":          &progress.WebCrawler,
		"URLSecurity_end":         &progress.URLSecurity,
		"DirScan_end":             &progress.DirScan,
		"VulnerabilityScan_end":   &progress.VulnerabilityScan,
		"scan_end":                &progress.All,
	}

	for redisField, progressField := range endFieldMappings {
		if endTime, exists := redisData[redisField]; exists {
			(*progressField)[1] = endTime
		}
	}

	// 设置节点信息
	if node, exists := redisData["node"]; exists {
		progress.Node = node
	}
}

func (s *service) TaskProgressNumber(ctx *gin.Context) error {
	filter := bson.M{
		"progress": bson.M{"$ne": 100},
		"status":   1,
	}
	// 调用 repository
	tasks, err := s.taskRepo.Find(ctx, filter, nil)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		fmt.Printf("任务ID: %s, 进度: %d, 状态: %d\n", task.ID.Hex(), task.Progress, task.Status)
		id := task.ID.Hex()
		key := fmt.Sprintf("TaskInfo:tmp:%s", id)
		fmt.Printf("<UNK>ID: %s\n", key)
	}
	return nil
}

// TaskProgress 更新任务进度
func (s *service) TaskProgress(ctx *gin.Context) error {
	// 查询所有进度不为100且状态为1（运行中）的任务
	filter := bson.M{
		"progress": bson.M{"$ne": 100},
		"status":   1,
	}

	tasks, err := s.taskRepo.Find(ctx, filter, nil)
	if err != nil {
		return fmt.Errorf("failed to find running tasks: %w", err)
	}

	// 如果没有运行中的任务，直接返回
	if len(tasks) == 0 {
		return nil
	}

	// 处理每个任务
	for _, task := range tasks {
		if err := s.ProcessTaskProgress(ctx, task); err != nil {
			logger.Error(fmt.Sprintf("failed to process task progress for task %s: %v", task.ID.Hex(), err))
			continue
		}
	}

	return nil
}

// processTaskProgress 处理单个任务的进度更新
func (s *service) ProcessTaskProgress(ctx *gin.Context, task models.Task) error {
	taskID := task.ID.Hex()
	tmpKey := fmt.Sprintf("TaskInfo:tmp:%s", taskID)
	timeKey := fmt.Sprintf("TaskInfo:time:%s", taskID)

	// 检查Redis中是否存在临时统计键
	exists, err := s.taskRepo.Exists(ctx, tmpKey)
	if err != nil {
		return fmt.Errorf("failed to check tmp key existence: %w", err)
	}

	if exists {
		// 获取集合中元素的数量
		count, err := s.taskRepo.SCard(ctx, tmpKey)
		if err != nil {
			return fmt.Errorf("failed to get set cardinality: %w", err)
		}

		// 计算进度
		progressTmp := float64(count) / float64(task.TaskNum)
		progressTmp = math.Round(progressTmp*100*10) / 10 // 保留一位小数

		// 确保进度不超过100
		if progressTmp > 100 {
			progressTmp = 100
		}

		// 如果进度达到100，更新任务状态为完成
		if progressTmp == 100 {
			// 获取结束时间
			endTime, err := s.taskRepo.Get(ctx, timeKey)
			if err != nil {
				return fmt.Errorf("failed to get end time: %w", err)
			}

			// 更新任务状态为完成（状态3）并设置结束时间
			update := bson.M{
				"$set": bson.M{
					"endTime":  endTime,
					"status":   3,
					"progress": int(progressTmp),
				},
			}
			if err := s.taskRepo.UpdateOne(ctx, bson.M{"_id": task.ID}, update); err != nil {
				return fmt.Errorf("failed to update task status: %w", err)
			}

			// 删除Redis统计信息
			if err := s.taskRepo.Del(ctx, tmpKey, timeKey); err != nil {
				logger.Error(fmt.Sprintf("failed to delete redis keys for completed task %s: %v", taskID, err))
			}
		} else {
			// 更新任务进度
			update := bson.M{
				"$set": bson.M{
					"progress": int(progressTmp),
				},
			}
			if err := s.taskRepo.UpdateOne(ctx, bson.M{"_id": task.ID}, update); err != nil {
				return fmt.Errorf("failed to update task progress: %w", err)
			}
		}
	} else {
		// Redis键不存在，将进度设为0
		update := bson.M{
			"$set": bson.M{
				"progress": 0,
			},
		}
		if err := s.taskRepo.UpdateOne(ctx, bson.M{"_id": task.ID}, update); err != nil {
			return fmt.Errorf("failed to reset task progress: %w", err)
		}
	}

	return nil
}
