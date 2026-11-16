// ct-----------------------------------------
// @file      : sdk.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/01/24 21:00
// -------------------------------------------
// CT日志监听SDK - 主要接口实现
// ============================================

package ct

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils/ct/models"
	"github.com/Autumn-27/ScopeSentry/internal/utils/ct/types"
)

// CTWatcherSDK CT日志监听器SDK
type CTWatcherSDK struct {
	// 配置
	opts *types.SDKOptions

	// 核心组件
	domainMatcher     *DomainMatcher
	streamProcessor   *StreamProcessor
	checkpointSaver   *AutoCheckpointSaver
	checkpointManager *RedisCheckpointManager

	// 监听器管理
	watchers map[string]*CTWatcher
	wg       sync.WaitGroup
	mu       sync.RWMutex

	// 状态管理
	isRunning int32
	status    *models.SDKStatus
	ctx       context.Context
	cancel    context.CancelFunc

	// 事件处理
	eventBus *models.EventBus

	// 回调函数
	subdomainCallbacks []models.SubdomainCallback
	errorCallbacks     []models.ErrorCallback
}

// NewCTWatcherSDK 创建新的CT监听器SDK
func NewCTWatcherSDK(opts *types.SDKOptions) (*CTWatcherSDK, error) {
	if opts == nil {
		opts = types.DefaultSDKOptions()
	}

	// 验证配置
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sdk := &CTWatcherSDK{
		opts:               opts,
		watchers:           make(map[string]*CTWatcher),
		status:             models.NewSDKStatus(),
		ctx:                ctx,
		cancel:             cancel,
		eventBus:           models.NewEventBus(),
		subdomainCallbacks: make([]models.SubdomainCallback, 0),
		errorCallbacks:     make([]models.ErrorCallback, 0),
	}

	// 初始化Redis检查点管理器
	checkpointManager, err := NewRedisCheckpointManager(opts.Redis.KeyPrefix)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("初始化Redis检查点管理器失败: %w", err)
	}
	sdk.checkpointManager = checkpointManager

	// 初始化自动检查点保存器
	sdk.checkpointSaver = NewAutoCheckpointSaver(checkpointManager, 10*time.Minute)

	// 初始化域名匹配器
	sdk.domainMatcher = NewDomainMatcher(opts.DomainMatcher)

	// 初始化流处理器
	sdk.streamProcessor = NewStreamProcessor(sdk.domainMatcher, sdk.eventBus)

	// 添加初始域名
	for _, domain := range opts.Domains {
		if err := sdk.domainMatcher.AddDomain(domain); err != nil {
			// 记录警告但不中断初始化
			continue
		}
	}

	// 注册事件处理器
	sdk.eventBus.Subscribe(models.EventTypeSubdomainDiscovered, &SubdomainEventHandler{sdk: sdk})
	sdk.eventBus.Subscribe(models.EventTypeCertificateProcessed, &CertificateEventHandler{sdk: sdk})
	sdk.eventBus.Subscribe(models.EventTypeErrorOccurred, &ErrorEventHandler{sdk: sdk})

	return sdk, nil
}

// Start 启动SDK
func (sdk *CTWatcherSDK) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&sdk.isRunning, 0, 1) {
		return fmt.Errorf("SDK已在运行")
	}

	sdk.status.StartTime = time.Now()
	sdk.status.IsRunning = true

	// 创建派生上下文
	sdkCtx, cancel := context.WithCancel(ctx)
	sdk.cancel = cancel

	// 为每个配置的日志服务器创建监听器
	for _, logServer := range sdk.opts.LogServers {
		if err := sdk.addWatcher(logServer); err != nil {
			// 记录错误但继续创建其他监听器
			sdk.handleError(fmt.Errorf("创建监听器失败 [%s]: %w", logServer.URL, err), nil)
			continue
		}
	}

	// 启动检查点保存器
	if err := sdk.checkpointSaver.Start(); err != nil {
		atomic.StoreInt32(&sdk.isRunning, 0)
		return fmt.Errorf("启动检查点保存器失败: %w", err)
	}

	if len(sdk.watchers) == 0 {
		atomic.StoreInt32(&sdk.isRunning, 0)
		sdk.checkpointSaver.Stop()
		return fmt.Errorf("没有成功创建任何监听器")
	}

	// 等待上下文取消或所有监听器退出
	go func() {
		<-sdkCtx.Done()
		sdk.stopInternal()
	}()

	return nil
}

// Stop 停止SDK
func (sdk *CTWatcherSDK) Stop() error {
	if !atomic.CompareAndSwapInt32(&sdk.isRunning, 1, 0) {
		return nil
	}

	sdk.cancel()
	return nil
}

// stopInternal 内部停止方法
func (sdk *CTWatcherSDK) stopInternal() {
	sdk.mu.Lock()
	defer sdk.mu.Unlock()

	// 停止所有监听器
	for url, watcher := range sdk.watchers {
		if err := watcher.Stop(); err != nil {
			sdk.handleError(fmt.Errorf("停止监听器失败 [%s]: %w", url, err), nil)
		}
	}

	// 停止检查点保存器
	sdk.checkpointSaver.Stop()

	// 关闭检查点管理器
	sdk.checkpointManager.Close()

	sdk.status.IsRunning = false
}

// AddDomains 动态添加关注的域名
func (sdk *CTWatcherSDK) AddDomains(domains []string) error {
	for _, domain := range domains {
		if err := sdk.domainMatcher.AddDomain(domain); err != nil {
			return fmt.Errorf("添加域名失败 [%s]: %w", domain, err)
		}
	}
	return nil
}

// RemoveDomains 动态移除关注的域名
func (sdk *CTWatcherSDK) RemoveDomains(domains []string) error {
	for _, domain := range domains {
		if err := sdk.domainMatcher.RemoveDomain(domain); err != nil {
			return fmt.Errorf("移除域名失败 [%s]: %w", domain, err)
		}
	}
	return nil
}

// GetWatchedDomains 获取当前关注的域名
func (sdk *CTWatcherSDK) GetWatchedDomains() []string {
	return sdk.domainMatcher.GetAllDomains()
}

// GetStatus 获取SDK状态
func (sdk *CTWatcherSDK) GetStatus() *types.Status {
	sdk.mu.RLock()
	defer sdk.mu.RUnlock()

	status := &types.Status{
		IsRunning:             atomic.LoadInt32(&sdk.isRunning) == 1,
		StartTime:             sdk.status.StartTime,
		Uptime:                time.Since(sdk.status.StartTime),
		ProcessedCertificates: sdk.status.TotalProcessedCertificates,
		DiscoveredSubdomains:  sdk.status.TotalDiscoveredSubdomains,
		ActiveWatchers:        len(sdk.watchers),
		RecentErrors:          make([]string, 0),
		WatcherStatuses:       make(map[string]*models.LogStatus),
	}

	// 更新监听器状态
	for url, watcher := range sdk.watchers {
		watcherStatus := watcher.GetStatus()
		status.WatcherStatuses[url] = &models.LogStatus{
			LogServer:      watcherStatus.LogServer,
			IsConnected:    watcherStatus.IsConnected,
			LastIndex:      watcherStatus.LastProcessedIndex,
			ErrorCount:     watcherStatus.ErrorCount,
			ProcessedCount: watcherStatus.ProcessedCount,
			Status:         watcherStatus.Status,
		}
	}

	return status
}

// GetStats 获取统计信息
func (sdk *CTWatcherSDK) GetStats() *models.DomainMatcherStats {
	if sdk.domainMatcher == nil {
		return &models.DomainMatcherStats{}
	}

	return sdk.domainMatcher.GetStats()
}

// OnSubdomainFound 注册子域名发现回调
func (sdk *CTWatcherSDK) OnSubdomainFound(callback models.SubdomainCallback) error {
	if callback == nil {
		return fmt.Errorf("回调函数不能为空")
	}

	sdk.subdomainCallbacks = append(sdk.subdomainCallbacks, callback)
	return nil
}

// OnError 注册错误回调
func (sdk *CTWatcherSDK) OnError(callback models.ErrorCallback) error {
	if callback == nil {
		return fmt.Errorf("回调函数不能为空")
	}

	sdk.errorCallbacks = append(sdk.errorCallbacks, callback)
	return nil
}

// addWatcher 添加监听器
func (sdk *CTWatcherSDK) addWatcher(logServer *types.LogServerConfig) error {
	sdk.mu.Lock()
	defer sdk.mu.Unlock()

	if _, exists := sdk.watchers[logServer.URL]; exists {
		return fmt.Errorf("监听器已存在: %s", logServer.URL)
	}

	watcher := NewCTWatcher(
		logServer,
		sdk.opts.Performance,
		sdk.streamProcessor,
		sdk.checkpointSaver,
	)

	sdk.watchers[logServer.URL] = watcher

	// 启动监听器
	if err := watcher.Start(); err != nil {
		delete(sdk.watchers, logServer.URL)
		return fmt.Errorf("启动监听器失败: %w", err)
	}

	return nil
}

// SubdomainEventHandler 子域名事件处理器
type SubdomainEventHandler struct {
	sdk *CTWatcherSDK
}

func (h *SubdomainEventHandler) HandleEvent(event models.Event) error {
	if subdomainEvent, ok := event.(*models.SubdomainDiscoveredEvent); ok {
		atomic.AddInt64(&h.sdk.status.TotalDiscoveredSubdomains, 1)

		// 触发用户注册的回调
		for _, callback := range h.sdk.subdomainCallbacks {
			go func(cb models.SubdomainCallback) {
				defer func() {
					if r := recover(); r != nil {
						h.sdk.handleError(fmt.Errorf("子域名回调执行失败: %v", r), nil)
					}
				}()
				cb(subdomainEvent.SubdomainEvent)
			}(callback)
		}
	}
	return nil
}

// CertificateEventHandler 证书处理事件处理器
type CertificateEventHandler struct {
	sdk *CTWatcherSDK
}

func (h *CertificateEventHandler) HandleEvent(event models.Event) error {
	if certEvent, ok := event.(*models.CertificateProcessedEvent); ok {
		atomic.AddInt64(&h.sdk.status.TotalProcessedCertificates, 1)

		// 这里可以添加证书处理回调，如果将来需要的话
		_ = certEvent // 暂时未使用，避免编译警告
	}
	return nil
}

// ErrorEventHandler 错误事件处理器
type ErrorEventHandler struct {
	sdk *CTWatcherSDK
}

func (h *ErrorEventHandler) HandleEvent(event models.Event) error {
	if errorEvent, ok := event.(*models.ErrorOccurredEvent); ok {
		atomic.AddInt64(&h.sdk.status.TotalErrors, 1)

		// 触发用户注册的回调
		for _, callback := range h.sdk.errorCallbacks {
			go func(cb models.ErrorCallback) {
				defer func() {
					if r := recover(); r != nil {
						h.sdk.handleError(fmt.Errorf("错误回调执行失败: %v", r), nil)
					}
				}()
				cb(fmt.Errorf("%s", errorEvent.ErrorMessage), map[string]interface{}{
					"error_type": errorEvent.ErrorType,
					"severity":   errorEvent.Severity,
				})
			}(callback)
		}
	}
	return nil
}

// handleError 处理内部错误
func (sdk *CTWatcherSDK) handleError(err error, context map[string]interface{}) {
	errorEvent := &models.ErrorOccurredEvent{
		BaseEvent: models.BaseEvent{
			Type:      models.EventTypeErrorOccurred,
			ID:        models.GenerateEventID(),
			Timestamp: time.Now(),
		},
		ErrorType:    "internal_error",
		ErrorMessage: err.Error(),
		Context:      context,
		Severity:     "medium",
		Recoverable:  true,
	}

	sdk.eventBus.Publish(errorEvent)
}

// IsRunning 检查SDK是否正在运行
func (sdk *CTWatcherSDK) IsRunning() bool {
	return atomic.LoadInt32(&sdk.isRunning) == 1
}

// Close 关闭SDK并释放资源
func (sdk *CTWatcherSDK) Close() error {
	if sdk.IsRunning() {
		sdk.Stop()
	}

	if sdk.checkpointManager != nil {
		sdk.checkpointManager.Close()
	}

	return nil
}
