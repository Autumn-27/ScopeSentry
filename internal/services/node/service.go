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

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/logger"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/node"
	"go.uber.org/zap"
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

// GetNodeData 获取所有节点数据。
//
// 存活判定见 ClassifyLiveness：容忍节点与服务器之间的时区/时钟差异，
// 解析失败时保留节点而不是把它从列表里抹掉（#254）。
func (s *service) GetNodeData(ctx context.Context, online bool) ([]models.NodeData, error) {
	nodes, err := s.nodeRepo.GetAllNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get node keys: %w", err)
	}

	loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		logger.Warn("Failed to load timezone, fall back to UTC",
			zap.String("timezone", config.GlobalConfig.System.Timezone), zap.Error(err))
		loc = time.UTC
	}
	now := time.Now().In(loc)
	timeout := NodeTimeout()

	result := make([]models.NodeData, 0, len(nodes))
	for _, nodeData := range nodes {
		lv := ClassifyLiveness(nodeData.State, nodeData.UpdateTime, now, loc, timeout)
		switch {
		case lv.ParseFailed:
			logger.Warn("Node heartbeat timestamp unparsable, keeping node",
				zap.String("node", nodeData.Name), zap.String("time", nodeData.UpdateTime))
		case lv.ClockSkew:
			logger.Debug("Node clock is ahead of server",
				zap.String("node", nodeData.Name), zap.Duration("skew", -lv.TimeDiff))
		case lv.UsedUTC:
			logger.Debug("Node heartbeat matched UTC interpretation",
				zap.String("node", nodeData.Name))
		}

		if nodeData.State == "1" && !lv.Alive {
			// 异步更新节点状态，避免阻塞列表
			go s.updateNodeState(context.Background(), nodeData.Name)
			nodeData.State = "3"
		}

		if online && nodeData.State != "1" {
			continue
		}
		result = append(result, nodeData)
	}

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
