// Package taskstate 负责任务状态与进度的纯计算。
//
// 单独成包是为了让状态机可以在不依赖 Mongo/Redis 的情况下被单元测试覆盖
// （任务服务本身会通过 assets 服务间接 import bootstrap，init 阶段就要连数据库）。
package taskstate

import (
	"math"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
)

// Derive 根据已完成目标数推导任务进度（整数百分比）与下一个状态。
//
// count 为 TaskInfo:tmp:<taskId> 集合中的元素数量，taskNum 为任务目标总数。
// 规则：
//   - 任务目标数非法（<=0）时不改状态，避免除零；
//   - 进度 0 且当前为待开始 → 保持待开始（扫描节点还没开始）；
//   - 进度 >0 且 <100 → 运行中（待开始的任务在这里转为运行中）；
//   - 进度达到 100 → 运行完成；
//   - 暂停中的任务保持暂停，不被进度轮询“复活”。
func Derive(currentStatus int, count, taskNum int64) (progress int, status int) {
	if taskNum <= 0 {
		return 0, currentStatus
	}
	ratio := float64(count) / float64(taskNum)
	ratio = math.Round(ratio*100*10) / 10 // 保留一位小数
	if ratio > 100 {
		ratio = 100
	}
	progress = int(ratio)
	switch {
	case progress >= 100:
		status = constants.TaskStatusFinished
	case currentStatus == constants.TaskStatusPaused:
		status = constants.TaskStatusPaused
	case progress > 0:
		status = constants.TaskStatusRunning
	default:
		status = currentStatus
	}
	return progress, status
}
