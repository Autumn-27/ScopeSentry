package bootstrap

import (
	"context"
	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"log"
	"time"
)

// StartupEvent 启动事件处理函数
type StartupEvent func(ctx context.Context) error

var startupEvents []StartupEvent

// RegisterStartupEvent 注册启动事件
func RegisterStartupEvent(event StartupEvent) {
	startupEvents = append(startupEvents, event)
}

// ExecuteStartupEvents 执行所有启动事件
func ExecuteStartupEvents() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, event := range startupEvents {
		if err := event(ctx); err != nil {
			return err
		}
	}
	return nil
}

// 注册默认的启动事件
func init() {
	// 加载配置
	//RegisterStartupEvent(func(ctx context.Context) error {
	//	return config.Load()
	//})

	//// 初始化日志
	//RegisterStartupEvent(func(ctx context.Context) error {
	//	return logger.Init()
	//})

	// 初始化MongoDB
	//RegisterStartupEvent(func(ctx context.Context) error {
	//	return mongodb.Init()
	//})

	// 初始化Redis
	//RegisterStartupEvent(func(ctx context.Context) error {
	//	return redis.Init()
	//})

	//// 初始化数据库
	RegisterStartupEvent(func(ctx context.Context) error {
		return mongodb.CreateDatabase()
	})

	// 这里可以添加更多的启动事件
	// 例如：初始化缓存、加载插件、启动定时任务等
	err := ExecuteStartupEvents()
	if err != nil {
		log.Fatalf("Error during startup events: %v", err)
		return
	}
}
