// Package jobs -----------------------------
// @file      : page_monitoring.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/30 15:18
// -------------------------------------------
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePageMonitoringTask() error {
	logger.Info("create_page_monitoring_task")
	filter := bson.M{"id": "page_monitoring"}
	opts := options.FindOne()
	var task models.Task
	err := ScanSvc.collection.FindOne(context.Background(), filter, opts).Decode(&task)
	if err != nil {
		return fmt.Errorf("failed to get scheduled task: %w", err)
	}
	if task.AllNode {
		// 选择自动加入后 获取所有节点
		allNode, err := NodeService.GetNodeData(context.Background(), false)
		if err != nil {
			return err
		}
		var nameMap = make(map[string]struct{})
		for _, n := range task.Node {
			nameMap[n] = struct{}{}
		}
		for _, n := range allNode {
			if _, exists := nameMap[n.Name]; !exists {
				nameMap[n.Name] = struct{}{}
				task.Node = append(task.Node, n.Name)
			}
		}
	}
	data, err := getPageMonitoringData()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	err = CommonService.Repo.Del(context.Background(), "TaskInfo:page_monitoring")
	if err != nil {
		return err
	}

	err = CommonService.Repo.LPush(context.Background(), "TaskInfo:page_monitoring", data)
	if err != nil {
		return err
	}
	addRedisTaskData := map[string]interface{}{
		"ID":   "page_monitoring",
		"type": "page_monitoring",
	}
	jsonData, err := json.Marshal(addRedisTaskData)
	if err != nil {
		fmt.Println("JSON 转换错误:", err)
		return err
	}
	for _, d := range data {
		err := CommonService.Repo.RPush(context.Background(), fmt.Sprintf("NodeTask:%v", d), jsonData)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("CommonService.Repo.RPush %v", err))
		}
	}
	return nil
}

// getPageMonitoringData 获取页面监控数据
func getPageMonitoringData() ([]string, error) {
	cursor, err := ScanSvc.pageMonitoringCollection.Find(context.Background(), nil, options.Find().SetProjection(bson.M{"url": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(context.Background())
	var urls []string
	for cursor.Next(context.Background()) {
		var task models.PageMonitoringTask
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}
		urls = append(urls, task.URL)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("getPageMonitoringData cursor error: %w", err)
	}
	return urls, nil
}
