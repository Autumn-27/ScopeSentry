package redis

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GlobalConfig.Redis.IP, config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       0,
	})

	// 测试连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Error(fmt.Sprintf("failed to connect to Redis: %v", err))
		return
	}

	Client = client
	logger.Info("Redis connected successfully")
	return
}

func Close() error {
	if Client != nil {
		if err := Client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %v", err)
		}
	}
	return nil
}
