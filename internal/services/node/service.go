// Package node -----------------------------
// @file      : service.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/30 15:01
// -------------------------------------------
package node

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/node"
	"go.uber.org/zap"
)

const (
	NodeTimeout = 300 // 节点超时时间（秒）
)

type Service interface {
	GetNodeData(ctx context.Context, online bool) ([]models.NodeData, error)
	RefreshConfig(ctx context.Context, msg models.Message) error
	DeleteNodes(ctx context.Context, names []string) error
	GetNodeLogs(ctx context.Context, name string) (string, error)
	GetNodePlugin(ctx context.Context, nodeName string) ([]models.NodePluginInfo, error)
	RestartNode(ctx context.Context, nodeName string) error
}

type service struct {
	nodeRepo node.Repository
}

func NewService() Service {
	return &service{
		nodeRepo: node.NewRepository(),
	}
}

// GetNodeData 获取所有节点数据
func (s *service) GetNodeData(ctx context.Context, online bool) ([]models.NodeData, error) {
	// 获取所有以 node: 开头的键
	nodes, err := s.nodeRepo.GetAllNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get node keys: %w", err)
	}

	var result []models.NodeData
	for _, nodeData := range nodes {
		layout := "2006-01-02 15:04:05"
		// 检查节点状态
		if nodeData.State == "1" {
			loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
			if err != nil {
				logger.Error("Failed to load timezone", zap.String("timezone", config.GlobalConfig.System.Timezone), zap.Error(err))
				continue
			}
			updateTime, err := time.ParseInLocation(layout, nodeData.UpdateTime, loc)
			if err != nil {
				logger.Error("Failed to parse update time",
					zap.String("time", nodeData.UpdateTime),
					zap.Error(err))
				continue
			}
			currentTime := time.Now().In(loc)
			timeDiff := currentTime.Sub(updateTime).Seconds()
			logger.Info("Node time difference",
				zap.Float64("time_diff", timeDiff),
				zap.String("current_time", currentTime.Format(layout)),
				zap.String("update_time", nodeData.UpdateTime))

			if timeDiff > NodeTimeout {
				// 异步更新节点状态
				go s.updateNodeState(ctx, nodeData.Name)
				nodeData.State = "3"
			}
		}
		if online {
			if nodeData.State == "1" {
				result = append(result, nodeData)
			}
		} else {
			result = append(result, nodeData)
		}
	}

	// 按名称排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// updateNodeState 更新节点状态
func (s *service) updateNodeState(ctx context.Context, nodeID string) {
	err := s.nodeRepo.UpdateNodeState(ctx, nodeID, "3")
	if err != nil {
		logger.Error("Failed to update node state",
			zap.String("node_id", nodeID),
			zap.Error(err))
	}
}

func (s *service) RefreshConfig(ctx context.Context, msg models.Message) error {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {

		return fmt.Errorf("models.Message转换失败: %v", err)
	}

	msgStr := string(jsonBytes)
	if msg.Name == "all" {
		nodes, err := s.nodeRepo.GetAllNodes(ctx)
		if err != nil {
			return fmt.Errorf("failed to get node keys: %w", err)
		}
		for _, nodeData := range nodes {
			err := s.nodeRepo.RefreshConfig(ctx, nodeData.Name, msgStr)
			if err != nil {
				return err
			}
		}
	} else {
		err := s.nodeRepo.RefreshConfig(ctx, msg.Name, msgStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteNodes 删除多个节点
func (s *service) DeleteNodes(ctx context.Context, names []string) error {
	return s.nodeRepo.DeleteNodes(ctx, names)
}

// GetNodeLogs 获取节点日志
func (s *service) GetNodeLogs(ctx context.Context, name string) (string, error) {
	return s.nodeRepo.GetNodeLogs(ctx, name)
}

// GetNodePlugin 获取节点插件信息
func (s *service) GetNodePlugin(ctx context.Context, nodeName string) ([]models.NodePluginInfo, error) {
	return s.nodeRepo.GetNodePlugin(ctx, nodeName)
}

// RestartNode 重启节点
func (s *service) RestartNode(ctx context.Context, nodeName string) error {
	msg := models.Message{
		Name:    nodeName,
		Type:    "restart",
		Content: "",
	}
	return s.RefreshConfig(ctx, msg)
}
