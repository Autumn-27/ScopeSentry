// ct/models-----------------------------------
// @file      : checkpoint.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/01/24 21:00
// -------------------------------------------
// CT日志监听SDK - 检查点和状态模型
// ============================================

package models

import (
	"time"
)

// Checkpoint 检查点数据结构
type Checkpoint struct {
	// 基本信息
	URL   string `json:"url" redis:"url"`       // CT日志URL
	LogID string `json:"log_id" redis:"log_id"` // 日志ID

	// 处理进度
	LastProcessedIndex int64 `json:"last_processed_index" redis:"last_processed_index"` // 最后处理的证书索引
	STHTreeSize        int64 `json:"sth_tree_size" redis:"sth_tree_size"`               // STH树大小

	// 元数据
	Timestamp  time.Time `json:"timestamp" redis:"timestamp"`     // 检查点保存时间
	Version    int       `json:"version" redis:"version"`         // 检查点版本
	SDKVersion string    `json:"sdk_version" redis:"sdk_version"` // SDK版本
}

// NewCheckpoint 创建新的检查点
func NewCheckpoint(url string) *Checkpoint {
	return &Checkpoint{
		URL:       url,
		Version:   1,
		Timestamp: time.Now(),
	}
}

// Update 更新检查点
func (cp *Checkpoint) Update(lastIndex, sthSize int64) {
	cp.LastProcessedIndex = lastIndex
	cp.STHTreeSize = sthSize
	cp.Timestamp = time.Now()
}

// IsValid 检查检查点是否有效
func (cp *Checkpoint) IsValid() bool {
	return cp.URL != "" && cp.LastProcessedIndex >= 0
}

// CheckpointManager 检查点管理器接口
type CheckpointManager interface {
	// 保存检查点
	SaveCheckpoint(cp *Checkpoint) error

	// 加载检查点
	LoadCheckpoint(url string) (*Checkpoint, error)

	// 删除检查点
	DeleteCheckpoint(url string) error

	// 列出所有检查点
	ListCheckpoints() ([]*Checkpoint, error)

	// 清理过期检查点
	CleanupExpired(olderThan time.Duration) error
}

// WatcherStatus 监听器状态
type WatcherStatus struct {
	// 基本信息
	LogServer *LogServerConfig `json:"log_server"` // 日志服务器配置
	LogID     string           `json:"log_id"`     // 日志ID

	// 连接状态
	IsConnected bool   `json:"is_connected"` // 是否连接成功
	Status      string `json:"status"`       // 当前状态 (connecting/connected/disconnected/error)

	// 处理统计
	LastProcessedIndex int64 `json:"last_processed_index"` // 最后处理的索引
	ProcessedCount     int64 `json:"processed_count"`      // 已处理证书数
	ErrorCount         int64 `json:"error_count"`          // 错误次数
	RetryCount         int64 `json:"retry_count"`          // 重试次数

	// 时间信息
	StartedAt     time.Time  `json:"started_at"`      // 启动时间
	LastActivity  time.Time  `json:"last_activity"`   // 最后活动时间
	LastErrorTime *time.Time `json:"last_error_time"` // 最后错误时间
	NextRetryTime *time.Time `json:"next_retry_time"` // 下次重试时间

	// 性能指标
	ProcessingRate float64       `json:"processing_rate"` // 处理速率 (证书/秒)
	AverageLatency time.Duration `json:"average_latency"` // 平均延迟
}

// NewWatcherStatus 创建新的监听器状态
func NewWatcherStatus(logServer *LogServerConfig) *WatcherStatus {
	now := time.Now()
	return &WatcherStatus{
		LogServer:    logServer,
		Status:       "disconnected",
		StartedAt:    now,
		LastActivity: now,
	}
}

// UpdateStatus 更新状态
func (ws *WatcherStatus) UpdateStatus(status string) {
	ws.Status = status
	ws.LastActivity = time.Now()

	if status == "connected" {
		ws.IsConnected = true
	} else if status == "disconnected" || status == "error" {
		ws.IsConnected = false
	}
}

// RecordProcessing 记录处理统计
func (ws *WatcherStatus) RecordProcessing(count int64) {
	ws.ProcessedCount += count
	ws.LastActivity = time.Now()
}

// RecordError 记录错误
func (ws *WatcherStatus) RecordError() {
	now := time.Now()
	ws.ErrorCount++
	ws.LastErrorTime = &now
	ws.LastActivity = now
}

// SDKStatus SDK整体状态
type SDKStatus struct {
	// 运行状态
	IsRunning bool          `json:"is_running"` // 是否正在运行
	StartTime time.Time     `json:"start_time"` // 启动时间
	Uptime    time.Duration `json:"uptime"`     // 运行时间

	// 全局统计
	TotalProcessedCertificates int64 `json:"total_processed_certificates"` // 总处理证书数
	TotalDiscoveredSubdomains  int64 `json:"total_discovered_subdomains"`  // 总发现子域名数
	TotalErrors                int64 `json:"total_errors"`                 // 总错误数

	// 资源使用
	MemoryUsage    int64   `json:"memory_usage"`    // 内存使用量(bytes)
	CPUUsage       float64 `json:"cpu_usage"`       // CPU使用率(0-1)
	GoroutineCount int     `json:"goroutine_count"` // 协程数量

	// 监听器状态
	ActiveWatchers  int                       `json:"active_watchers"`  // 活跃监听器数
	WatcherStatuses map[string]*WatcherStatus `json:"watcher_statuses"` // 各监听器状态

	// 域名统计
	TotalDomains   int64   `json:"total_domains"`   // 总关注域名数
	MatchedDomains int64   `json:"matched_domains"` // 匹配到的域名数
	MatchRate      float64 `json:"match_rate"`      // 匹配率
}

// NewSDKStatus 创建新的SDK状态
func NewSDKStatus() *SDKStatus {
	return &SDKStatus{
		StartTime:       time.Now(),
		WatcherStatuses: make(map[string]*WatcherStatus),
	}
}

// UpdateStats 更新统计信息
func (ss *SDKStatus) UpdateStats() {
	ss.Uptime = time.Since(ss.StartTime)

	// 计算活跃监听器数
	activeCount := 0
	for _, status := range ss.WatcherStatuses {
		if status.IsConnected {
			activeCount++
		}
	}
	ss.ActiveWatchers = activeCount

	// 计算匹配率
	if ss.TotalProcessedCertificates > 0 {
		ss.MatchRate = float64(ss.TotalDiscoveredSubdomains) / float64(ss.TotalProcessedCertificates)
	}
}

// LogStatus 监听器状态（用于API响应）
type LogStatus struct {
	LogServer      *LogServerConfig `json:"log_server"`      // 日志服务器配置
	IsConnected    bool             `json:"is_connected"`    // 是否连接成功
	LastIndex      int64            `json:"last_index"`      // 最后处理的索引
	ErrorCount     int64            `json:"error_count"`     // 错误次数
	ProcessedCount int64            `json:"processed_count"` // 已处理证书数
	Status         string           `json:"status"`          // 当前状态
}

// LogServerConfig 日志服务器配置（简化版，用于状态显示）
type LogServerConfig struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	Operator    string `json:"operator"`
}

// ProcessingStats 处理统计信息
type ProcessingStats struct {
	// 时间窗口统计
	TimeWindow time.Duration `json:"time_window"` // 统计时间窗口

	// 处理指标
	CertificatesPerSecond float64       `json:"certificates_per_second"` // 证书处理速率
	SubdomainsPerSecond   float64       `json:"subdomains_per_second"`   // 子域名发现速率
	ErrorRate             float64       `json:"error_rate"`              // 错误率
	AverageLatency        time.Duration `json:"average_latency"`         // 平均处理延迟

	// 资源使用
	MemoryUsage int64   `json:"memory_usage"` // 内存使用
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率
	GCStats     GCStats `json:"gc_stats"`     // GC统计
}

// GCStats GC统计信息
type GCStats struct {
	NumGC         uint32    `json:"num_gc"`          // GC次数
	LastGCTime    time.Time `json:"last_gc_time"`    // 最后GC时间
	GCCPUFraction float64   `json:"gc_cpu_fraction"` // GC CPU占用率
}
