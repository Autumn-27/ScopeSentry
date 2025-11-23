// Package redis_log_subscriber -----------------------------
// @file      : log_subscriber.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/01/XX XX:XX
// -------------------------------------------
package redis_log_subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	redisdb "github.com/Autumn-27/ScopeSentry/internal/database/redis"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	// GET_LOG_NAME 存储需要获取日志的节点名称列表
	GET_LOG_NAME []string
	logNameMutex sync.RWMutex

	// LOG_INFO 存储每个节点的日志信息
	LOG_INFO     map[string][]string
	logInfoMutex sync.RWMutex
)

func init() {
	LOG_INFO = make(map[string][]string)
	GET_LOG_NAME = make([]string, 0)
}

// SubscribeLogChannel 订阅 Redis 日志通道并处理消息
func SubscribeLogChannel() {
	channelName := "logs"
	logger.Info("Subscribed to channel", zap.String("channel", channelName))

	for {
		trySubscribe(channelName)
		// 如果连接失败，等待1秒后重试
		time.Sleep(1 * time.Second)
	}
}

// trySubscribe 尝试订阅 Redis 通道
func trySubscribe(channelName string) {
	ctx := context.Background()
	pubsub := redisdb.Client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	// 验证订阅
	_, err := pubsub.Receive(ctx)
	if err != nil {
		logger.Error("Failed to subscribe to channel", zap.String("channel", channelName), zap.Error(err))
		return
	}

	ch := pubsub.Channel()
	logger.Info("Successfully subscribed to channel", zap.String("channel", channelName))

	for msg := range ch {
		if err := handleMessage(ctx, msg.Payload); err != nil {
			logger.Error("Failed to handle message", zap.Error(err))
			continue
		}
	}
}

var (
	logCounters   = make(map[string]int)
	logCountersMu sync.Mutex
)

// handleMessage 处理接收到的消息
func handleMessage(ctx context.Context, payload string) error {
	// 解析 JSON 数据
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	logName, ok := data["name"].(string)
	if !ok {
		return fmt.Errorf("invalid log name in message")
	}

	logContent, ok := data["log"].(string)
	if !ok {
		return fmt.Errorf("invalid log content in message")
	}

	//// 检查节点是否在 GET_LOG_NAME 中
	//logNameMutex.RLock()
	//inList := contains(GET_LOG_NAME, logName)
	//logNameMutex.RUnlock()
	//
	//if inList {
	//	logInfoMutex.Lock()
	//	if LOG_INFO[logName] == nil {
	//		LOG_INFO[logName] = make([]string, 0)
	//	}
	//	LOG_INFO[logName] = append(LOG_INFO[logName], logContent)
	//	logInfoMutex.Unlock()
	//}
	logger.Info(fmt.Sprintf("%s %s", logName, logContent))
	// 如果包含 "Register Success"，检查节点任务
	if strings.Contains(logContent, "Register Success") {
		if err := checkNodeTask(ctx, logName); err != nil {
			logger.Error("Failed to check node task", zap.String("node", logName), zap.Error(err))
		}
	}

	// 将日志推送到 Redis list
	logKey := fmt.Sprintf("log:%s", logName)
	if err := redisdb.Client.RPush(ctx, logKey, logContent).Err(); err != nil {
		return fmt.Errorf("failed to push log to Redis: %w", err)
	}

	logCountersMu.Lock()
	defer logCountersMu.Unlock()

	// 更新本地计数器
	logCounters[logName]++
	if logCounters[logName] > config.GlobalConfig.Log.TotalLogs {
		// 超过限制，直接清理 Redis，并重置本地计数器
		if err := redisdb.Client.Del(ctx, logKey).Err(); err != nil {
			logger.Error("Failed to delete log key", zap.String("key", logKey), zap.Error(err))
		}
		logCounters[logName] = 0
	}

	return nil
}

// checkNodeTask 检查节点任务并推送任务数据
func checkNodeTask(ctx context.Context, nodeName string) error {
	// 查询未完成且状态为1的任务
	filter := bson.M{
		"progress": bson.M{"$ne": 100},
		"status":   1,
		"$or": []bson.M{
			{"node": nodeName},
			{"allNode": true},
		},
	}

	opts := options.Find()
	cursor, err := mongodb.DB.Collection("task").Find(ctx, filter, opts)
	if err != nil {
		return fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return fmt.Errorf("failed to decode tasks: %w", err)
	}

	if len(tasks) == 0 {
		return nil
	}

	// 获取任务服务
	taskService := common.NewService()

	// 处理每个任务
	for _, taskDoc := range tasks {
		// 获取任务模板数据
		scanTemplate, err := taskService.GetScanTemplate(context.Background(), &taskDoc)
		if err != nil {
			logger.Error("Failed to get scan template", zap.String("task_id", taskDoc.ID.Hex()), zap.Error(err))
			continue
		}

		// 将任务数据序列化为 JSON
		taskDataJSON, err := json.Marshal(scanTemplate)
		if err != nil {
			logger.Error("Failed to marshal task data", zap.String("task_id", taskDoc.ID.Hex()), zap.Error(err))
			continue
		}

		// 推送到 Redis
		nodeTaskKey := fmt.Sprintf("NodeTask:%s", nodeName)
		if err := redisdb.Client.RPush(ctx, nodeTaskKey, taskDataJSON).Err(); err != nil {
			logger.Error("Failed to push node task", zap.String("node", nodeName), zap.Error(err))
			continue
		}
	}

	return nil
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// AddLogName 添加需要获取日志的节点名称
func AddLogName(name string) {
	logNameMutex.Lock()
	defer logNameMutex.Unlock()

	if !contains(GET_LOG_NAME, name) {
		GET_LOG_NAME = append(GET_LOG_NAME, name)
	}
}

// RemoveLogName 移除需要获取日志的节点名称
func RemoveLogName(name string) {
	logNameMutex.Lock()
	defer logNameMutex.Unlock()

	for i, n := range GET_LOG_NAME {
		if n == name {
			GET_LOG_NAME = append(GET_LOG_NAME[:i], GET_LOG_NAME[i+1:]...)
			break
		}
	}
}

// GetLogInfo 获取节点的日志信息
func GetLogInfo(name string) []string {
	logInfoMutex.RLock()
	defer logInfoMutex.RUnlock()

	logs := LOG_INFO[name]
	if logs == nil {
		return []string{}
	}
	return logs
}

// PopLogInfo 弹出节点的日志信息
func PopLogInfo(name string) string {
	logInfoMutex.Lock()
	defer logInfoMutex.Unlock()

	logs := LOG_INFO[name]
	if logs == nil || len(logs) == 0 {
		return ""
	}

	log := logs[0]
	LOG_INFO[name] = logs[1:]
	return log
}
