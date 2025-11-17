package scheduler

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	schedulerCore "github.com/Autumn-27/ScopeSentry-go/internal/scheduler"

	schedulerRepo "github.com/Autumn-27/ScopeSentry-go/internal/repositories/task/scheduler"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	GetScheduledData(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.Task, int64, error)
	CheckScheduledTaskNameExists(ctx *gin.Context, name string) (bool, error)
	CreateScheduledTask(ctx *gin.Context, req *models.Task) error
	UpdateScheduledTask(ctx *gin.Context, req *models.Task) error
	DeleteScheduledTasks(ctx *gin.Context, ids []string) error
	GetScheduledTaskDetail(ctx *gin.Context, id string) (*models.Task, error)
	UpdatePageMonitScheduledTask(ctx *gin.Context, req *models.Task) error
	GetPageMonitData(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.PageMonitoringTask, int64, error)
	AddPageMonitTask(ctx *gin.Context, url string) error
	DeletePageMonitTasks(ctx *gin.Context, ids []string) error
}

type service struct {
	schedulerRepo schedulerRepo.Repository
}

func NewService() Service {
	return &service{
		schedulerRepo: schedulerRepo.NewRepository(),
	}
}

// GetScheduledData 获取计划任务数据
func (s *service) GetScheduledData(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.Task, int64, error) {
	// 构建查询条件
	query := bson.M{}
	if search != "" {
		query["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	// 获取总数
	total, err := s.schedulerRepo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))

	tasks, err := s.schedulerRepo.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}

	// 计算周期描述
	for i := range tasks {
		tasks[i].Cycle = s.calculateCycle(&tasks[i])
	}

	return tasks, total, nil
}

// CheckScheduledTaskNameExists 检查计划任务名称是否存在
func (s *service) CheckScheduledTaskNameExists(ctx *gin.Context, name string) (bool, error) {
	query := bson.M{"name": name}
	count, err := s.schedulerRepo.Count(ctx, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateScheduledTask 创建计划任务
func (s *service) CreateScheduledTask(ctx *gin.Context, req *models.Task) error {
	// 设置默认值
	if req.Type == "" {
		req.Type = "scan"
	}
	if req.TargetSource == "" {
		req.TargetSource = "general"
	}
	if req.CycleType == "" {
		req.CycleType = "nhours"
	}
	if req.Day == 0 {
		req.Day = 1
	}
	if req.Minute == 0 {
		req.Minute = 0
	}
	if req.Hour == 0 {
		req.Hour = 1
	}
	if req.Week == 0 {
		req.Week = 1
	}
	req.ID = primitive.NewObjectID()
	req.TaskID = req.ID.Hex()
	// 插入数据库
	if err := s.schedulerRepo.Insert(ctx, req); err != nil {
		return err
	}

	// 如果需要启用调度，添加到调度器
	if req.ScheduledTasks {
		if err := s.addScheduler(req, "scan"); err != nil {
			return err
		}
	}

	return nil
}

// UpdateScheduledTask 更新计划任务
func (s *service) UpdateScheduledTask(ctx *gin.Context, req *models.Task) error {
	// 检查任务是否存在
	existingTask, err := s.GetScheduledTaskDetail(ctx, req.ID.Hex())
	if err != nil {
		return err
	}
	if existingTask == nil {
		return fmt.Errorf("scheduled task not found")
	}

	// 构建更新数据
	if req.Name != "" {
		existingTask.Name = req.Name
	}
	if req.Type != "" {
		existingTask.Type = req.Type
	}
	if len(req.Node) > 0 {
		existingTask.Node = req.Node
	}
	existingTask.AllNode = req.AllNode
	if len(req.Project) > 0 {
		existingTask.Project = req.Project
	}
	if req.TargetSource != "" {
		existingTask.TargetSource = req.TargetSource
	}
	if req.Day > 0 {
		existingTask.Day = req.Day
	}
	if req.Minute >= 0 {
		existingTask.Minute = req.Minute
	}
	if req.Hour >= 0 {
		existingTask.Hour = req.Hour
	}
	if req.Search != "" {
		existingTask.Search = req.Search
	}
	if req.CycleType != "" {
		existingTask.CycleType = req.CycleType
	}
	existingTask.ScheduledTasks = req.ScheduledTasks
	if req.Target != "" {
		existingTask.Target = req.Target
	}
	if req.Ignore != "" {
		existingTask.Ignore = req.Ignore
	}
	if req.Template != "" {
		existingTask.Template = req.Template
	}
	if req.Duplicates != "" {
		existingTask.Duplicates = req.Duplicates
	}
	if req.Week > 0 {
		existingTask.Week = req.Week
	}

	// 更新数据库
	filter := bson.M{"_id": req.ID}
	if err := s.schedulerRepo.UpdateOne(ctx, filter, existingTask); err != nil {
		return err
	}

	// 处理调度器
	if req.ScheduledTasks {
		// 先移除旧任务
		schedulerCore.GetGlobalScheduler().RemoveJob(req.ID.Hex())
		// 添加新任务
		if err := s.addScheduler(existingTask, "scan"); err != nil {
			return err
		}
	} else {
		// 移除调度任务
		schedulerCore.GetGlobalScheduler().RemoveJob(req.ID.Hex())
	}

	return nil
}

// DeleteScheduledTasks 删除计划任务
func (s *service) DeleteScheduledTasks(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建删除条件
	filter := bson.M{
		"$or": []bson.M{
			{"id": bson.M{"$in": ids}},
		},
	}

	// 处理特殊ID
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		if id != "page_monitoring" {
			if objID, err := primitive.ObjectIDFromHex(id); err == nil {
				objectIds = append(objectIds, objID)
			}
		}
	}

	if len(objectIds) > 0 {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"_id": bson.M{"$in": objectIds}})
	}

	// 从调度器中移除任务
	for _, id := range ids {
		schedulerCore.GetGlobalScheduler().RemoveJob(id)
	}

	// 删除数据库记录
	if err := s.schedulerRepo.DeleteMany(ctx, filter); err != nil {
		return err
	}

	return nil
}

// GetScheduledTaskDetail 获取计划任务详情
func (s *service) GetScheduledTaskDetail(ctx *gin.Context, id string) (*models.Task, error) {

	// 如果按id字段没找到，尝试按_id字段查询
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter := bson.M{"_id": objID}
		task, err := s.schedulerRepo.FindOne(ctx, filter)
		if err == nil && task != nil {
			return task, nil
		}
	}

	// 先尝试按id字段查询
	filter := bson.M{"id": id}
	task, err := s.schedulerRepo.FindOne(ctx, filter)
	if err == nil && task != nil {
		return task, nil
	}

	return nil, fmt.Errorf("scheduled task not found")
}

// addScheduler 添加调度任务
func (s *service) addScheduler(task *models.Task, tp string) error {
	// 创建调度任务
	job := &schedulerCore.Job{
		ID:      task.ID.Hex(),
		Name:    task.Name,
		Handler: tp,
		Params:  task.ID.Hex(),
	}

	// 根据周期类型设置任务
	switch task.CycleType {
	case "daily":
		job.CycleType = "cron"
		job.Schedule = fmt.Sprintf("%d %d * * *", task.Minute, task.Hour)
	case "ndays":
		job.CycleType = "interval"
		job.Interval.Days = task.Day
		job.Interval.Hours = task.Hour
		job.Interval.Minutes = task.Minute
	case "nhours":
		job.CycleType = "cron"

		// 防止 0 值
		if task.Hour <= 0 && task.Minute <= 0 {
			task.Hour = 1
			task.Minute = 0
		}

		// 如果分钟 < 0 或 > 59，强制 0
		if task.Minute < 0 || task.Minute > 59 {
			task.Minute = 0
		}

		// 构造 cron 表达式
		if task.Hour > 0 {
			// 每 task.Hour 小时执行一次
			job.Schedule = fmt.Sprintf("%d 0-23/%d * * *", task.Minute, task.Hour)
		} else {
			// 只有分钟
			job.Schedule = fmt.Sprintf("%d * * * *", task.Minute)
		}
	case "weekly":
		job.CycleType = "cron"
		job.Schedule = fmt.Sprintf("%d %d * * %d", task.Minute, task.Hour, task.Week)
	case "monthly":
		job.CycleType = "cron"
		job.Schedule = fmt.Sprintf("%d %d %d * *", task.Minute, task.Hour, task.Day)
	default:
		return fmt.Errorf("unsupported cycle type: %s", task.CycleType)
	}

	// 添加任务到调度器
	if err := schedulerCore.GetGlobalScheduler().AddJob(job); err != nil {
		return fmt.Errorf("failed to add scheduled task: %w", err)
	}

	return nil
}

// calculateCycle 计算周期描述
func (s *service) calculateCycle(task *models.Task) string {
	week := task.Week
	day := task.Day
	hour := task.Hour
	minute := task.Minute
	cycleType := task.CycleType

	if task.TaskID == "page_monitoring" {
		return fmt.Sprintf("%d hour", hour)
	}

	switch cycleType {
	case "daily":
		return fmt.Sprintf("Every day at %d:%02d", hour, minute)
	case "ndays":
		return fmt.Sprintf("Every %d days at %d:%02d", day, hour, minute)
	case "nhours":
		return fmt.Sprintf("Every %dh %dm", hour, minute)
	case "weekly":
		return fmt.Sprintf("Every week on day %d at %d:%02d", week, hour, minute)
	case "monthly":
		return fmt.Sprintf("Every month on day %d at %d:%02d", day, hour, minute)
	default:
		return "Unknown cycle"
	}
}

func (s *service) UpdatePageMonitScheduledTask(ctx *gin.Context, req *models.Task) error {

	filter := bson.M{"id": "page_monitoring"}
	existingTask, err := s.schedulerRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}
	if existingTask == nil {
		return fmt.Errorf("scheduled task not found")
	}
	existingTask.CycleType = "nhours"
	existingTask.AllNode = req.AllNode
	existingTask.ScheduledTasks = req.ScheduledTasks
	existingTask.Hour = req.Hour
	existingTask.ScheduledTasks = req.ScheduledTasks
	filter = bson.M{"_id": existingTask.ID}

	if err := s.schedulerRepo.UpdateOne(ctx, filter, existingTask); err != nil {
		return err
	}
	if existingTask.ScheduledTasks {
		// 先移除旧任务
		schedulerCore.GetGlobalScheduler().RemoveJob(existingTask.ID.Hex())
		// 添加新任务
		if err := s.addScheduler(existingTask, "page_monitoring"); err != nil {
			return err
		}
	} else {
		// 移除调度任务
		schedulerCore.GetGlobalScheduler().RemoveJob(existingTask.ID.Hex())
	}

	return nil
}

// GetPageMonitData 获取页面监控数据
func (s *service) GetPageMonitData(ctx *gin.Context, search string, pageIndex, pageSize int) ([]models.PageMonitoringTask, int64, error) {
	// 构建查询条件
	query := bson.M{}
	if search != "" {
		query["url"] = bson.M{"$regex": search, "$options": "i"}
	}

	// 获取总数
	total, err := s.schedulerRepo.CountPageMonit(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	opts := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))

	tasks, err := s.schedulerRepo.FindPageMonit(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// AddPageMonitTask 添加页面监控任务
func (s *service) AddPageMonitTask(ctx *gin.Context, url string) error {
	// 创建页面监控任务
	task := &models.PageMonitoringTask{
		ID:      primitive.NewObjectID(),
		URL:     url,
		Hash:    []string{},
		MD5:     calculateMD5(url),
		State:   1,
		Project: "",
		Time:    "",
	}

	// 插入数据库
	if err := s.schedulerRepo.InsertPageMonit(ctx, task); err != nil {
		return err
	}

	return nil
}

// DeletePageMonitTasks 删除页面监控任务
func (s *service) DeletePageMonitTasks(ctx *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建删除条件
	var objectIds []primitive.ObjectID
	for _, id := range ids {
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIds = append(objectIds, objID)
		}
	}

	if len(objectIds) == 0 {
		return fmt.Errorf("invalid ids")
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	// 删除数据库记录
	if err := s.schedulerRepo.DeletePageMonit(ctx, filter); err != nil {
		return err
	}

	return nil
}

// calculateMD5 计算MD5值
func calculateMD5(url string) string {
	return helper.CalculateMD5FromContent(url)
}
