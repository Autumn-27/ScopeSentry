package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/scheduler/jobs"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// JobHandler 表示任务处理函数的类型
type JobHandler func(params string, nextTime string) error

// Job 表示一个计划任务
type Job struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Schedule    string    `bson:"schedule"` // cron 表达式
	Handler     string    `bson:"handler"`  // 处理函数名称
	Params      string    `bson:"params"`
	LastRunTime time.Time `bson:"last_run_time"`
	NextRunTime time.Time `bson:"next_run_time"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	CycleType   string    `bson:"cycle_type"` // 周期类型：cron 或 interval
	Interval    struct {
		Days    int `bson:"days"`
		Hours   int `bson:"hours"`
		Minutes int `bson:"minutes"`
	} `bson:"interval"`
}

// Scheduler 表示计划任务调度器
type Scheduler struct {
	scheduler      *gocron.Scheduler
	db             *mongo.Database
	jobsCollection *mongo.Collection
	logger         *zap.Logger
	jobs           map[string]*gocron.Job
	handlers       map[string]JobHandler // 存储注册的处理函数
}

var (
	globalScheduler *Scheduler
	once            sync.Once
)

// GetGlobalScheduler 获取全局调度器实例
func GetGlobalScheduler() *Scheduler {
	return globalScheduler
}

// InitializeGlobalScheduler 初始化全局调度器
func InitializeGlobalScheduler() {
	once.Do(func() {
		globalScheduler = NewScheduler()
		globalScheduler.RegisterHandler("scan", func(params string, nextTime string) error {
			return jobs.Scan(params, nextTime)
		})
		globalScheduler.RegisterHandler("page_monitoring", func(params string, nextTime string) error {
			return jobs.CreatePageMonitoringTask()
		})
	})
}

// NewScheduler 创建一个新的调度器实例
func NewScheduler() *Scheduler {
	loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		// 处理错误，例如使用默认时区
		loc = time.Local
	}
	return &Scheduler{
		scheduler:      gocron.NewScheduler(loc),
		db:             mongodb.DB,
		jobsCollection: mongodb.DB.Collection("scheduled"),
		logger:         logger.Logger,
		jobs:           make(map[string]*gocron.Job),
		handlers:       make(map[string]JobHandler),
	}
}

// RegisterHandler 注册一个任务处理函数
func (s *Scheduler) RegisterHandler(name string, handler JobHandler) {
	s.handlers[name] = handler
	s.logger.Info("Registered job handler", zap.String("handler_name", name))
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.scheduler.StartAsync()
	s.loadJobs()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.scheduler.Stop()
}

// AddJob 添加一个新的任务
func (s *Scheduler) AddJob(job *Job) error {
	// 检查处理函数是否已注册
	if _, exists := s.handlers[job.Handler]; !exists {
		return fmt.Errorf("handler %s not registered", job.Handler)
	}

	// 将任务保存到 MongoDB
	job.CreatedAt = helper.GetNowTime()
	job.UpdatedAt = helper.GetNowTime()

	_, err := s.jobsCollection.InsertOne(context.Background(), job)
	if err != nil {
		return err
	}

	var schedJob *gocron.Job
	var err2 error

	// 根据周期类型添加任务
	if job.CycleType == "interval" {
		// 使用 interval 方式添加任务
		schedJob, err2 = s.scheduler.Every(job.Interval.Days).Days().
			At(fmt.Sprintf("%02d:%02d", job.Interval.Hours, job.Interval.Minutes)).
			Do(s.executeJob, job)
	} else {
		// 使用 cron 方式添加任务
		schedJob, err2 = s.scheduler.Cron(job.Schedule).Do(s.executeJob, job)
	}

	if err2 != nil {
		return err2
	}

	s.jobs[job.ID] = schedJob

	// 添加日志记录
	s.logger.Info("Successfully added new job",
		zap.String("job_id", job.ID),
		zap.String("job_name", job.Name),
		zap.String("cycle_type", job.CycleType),
		zap.Time("next_run_time", schedJob.NextRun()),
	)

	return nil
}

// RemoveJob 移除一个任务
func (s *Scheduler) RemoveJob(jobID string) error {
	// 从 MongoDB 中删除任务
	_, err := s.jobsCollection.DeleteOne(context.Background(), bson.M{"_id": jobID})
	if err != nil {
		return err
	}

	// 从调度器中移除任务
	if job, exists := s.jobs[jobID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.jobs, jobID)
	}

	return nil
}

// UpdateJob 更新任务
func (s *Scheduler) UpdateJob(job *Job) error {
	// 检查处理函数是否已注册
	if _, exists := s.handlers[job.Handler]; !exists {
		return fmt.Errorf("handler %s not registered", job.Handler)
	}

	job.UpdatedAt = helper.GetNowTime()

	// 更新 MongoDB 中的任务
	_, err := s.jobsCollection.ReplaceOne(context.Background(), bson.M{"_id": job.ID}, job)
	if err != nil {
		return err
	}

	// 从调度器中移除旧任务
	if oldJob, exists := s.jobs[job.ID]; exists {
		s.scheduler.RemoveByReference(oldJob)
	}

	// 添加新任务到调度器
	schedJob, err := s.scheduler.Cron(job.Schedule).Do(s.executeJob, job)
	if err != nil {
		return err
	}

	s.jobs[job.ID] = schedJob
	return nil
}

// loadJobs 从 MongoDB 加载所有任务
func (s *Scheduler) loadJobs() {
	cursor, err := s.jobsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		s.logger.Error("Failed to load jobs from MongoDB", zap.Error(err))
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var job Job
		if err := cursor.Decode(&job); err != nil {
			s.logger.Error("Failed to decode job", zap.Error(err))
			continue
		}

		// 检查处理函数是否已注册
		if _, exists := s.handlers[job.Handler]; !exists {
			s.logger.Error("Handler not registered for job",
				zap.String("job_id", job.ID),
				zap.String("handler", job.Handler),
			)
			continue
		}

		schedJob, err := s.scheduler.Cron(job.Schedule).Do(s.executeJob, &job)
		if err != nil {
			s.logger.Error("Failed to schedule job", zap.Error(err))
			continue
		}
		s.jobs[job.ID] = schedJob
	}
}

// executeJob 执行任务
func (s *Scheduler) executeJob(job *Job) {
	s.logger.Info("Executing job",
		zap.String("job_id", job.ID),
		zap.String("job_name", job.Name),
		zap.String("handler", job.Handler),
	)

	// 获取并执行处理函数
	handler, exists := s.handlers[job.Handler]
	if !exists {
		s.logger.Error("Handler not found",
			zap.String("job_id", job.ID),
			zap.String("handler", job.Handler),
		)
		return
	}

	// 更新最后执行时间
	job.LastRunTime = helper.GetNowTime()
	job.NextRunTime = s.jobs[job.ID].NextRun()

	// 执行处理函数
	if err := handler(job.Params, helper.FormatTime(job.NextRunTime)); err != nil {
		s.logger.Error("Job execution failed",
			zap.String("job_id", job.ID),
			zap.String("job_name", job.Name),
			zap.Error(err),
		)
	}

	_, err := s.jobsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": job.ID},
		bson.M{
			"$set": bson.M{
				"last_run_time": job.LastRunTime,
				"next_run_time": job.NextRunTime,
			},
		},
	)
	if err != nil {
		s.logger.Error("Failed to update job execution time", zap.Error(err))
	}
}
