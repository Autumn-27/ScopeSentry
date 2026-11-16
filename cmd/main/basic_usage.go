//go:build example

// examples----------------------------------
// @file      : basic_usage.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/01/24 21:00
// -------------------------------------------
// CT日志监听SDK - 基础使用示例
// ============================================

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils/ct"
	"github.com/Autumn-27/ScopeSentry/internal/utils/ct/models"
	"github.com/Autumn-27/ScopeSentry/internal/utils/ct/types"
)

func main() {
	// 创建SDK配置
	opts := &types.SDKOptions{
		// CT日志服务器配置
		LogServers: []*types.LogServerConfig{
			{
				URL:           "https://ct.googleapis.com/logs/us1/argon2026h1/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.googleapis.com/logs/us1/argon2026h2/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.googleapis.com/logs/us1/argon2027h1/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.googleapis.com/logs/eu1/xenon2026h1/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.googleapis.com/logs/eu1/xenon2026h2/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.googleapis.com/logs/eu1/xenon2027h1/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.cloudflare.com/logs/nimbus2026/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
			{
				URL:           "https://ct.cloudflare.com/logs/nimbus2027/",
				BatchSize:     500,
				ParallelFetch: 10,
				BufferSize:    1000,
				Timeout:       30 * time.Second,
				MaxRetries:    3,
			},
		},
		// 初始关注的域名
		Domains: []string{
			"example.com",
			"google.com",
			"github.com",
			"scope-sentry.top",
			"jgeai.site",
		},
		// 域名匹配器配置
		DomainMatcher: &types.DomainMatcherConfig{
			EnableBloomFilter: true,
			BloomFilterSize:   1000000,
			BloomFilterHash:   5,
			EnableLRUCache:    true,
			LRUCacheSize:      10000,
			CaseSensitive:     false,
			MaxDomainLength:   253,
		},
		// 性能配置
		Performance: &types.PerformanceConfig{
			MaxWorkers:        10,
			ChannelBufferSize: 100,
			ProcessingTimeout: 5 * time.Second,
			BatchProcessSize:  50,
		},
		// Redis配置
		Redis: &types.RedisConfig{
			KeyPrefix: "ct_watcher:",
		},
	}

	// 创建SDK实例
	sdk, err := ct.NewCTWatcherSDK(opts)
	if err != nil {
		log.Fatalf("创建SDK失败: %v", err)
	}
	defer sdk.Close()

	// 注册子域名发现回调
	sdk.OnSubdomainFound(func(event models.SubdomainEvent) {
		fmt.Printf("[发现子域名] %s -> %s | 证书索引: %d | 时间: %s\n",
			event.RootDomain,
			event.Subdomain,
			event.Certificate.Index,
			event.DiscoveredAt.Format("2006-01-02 15:04:05"))
	})

	// 注册错误回调
	sdk.OnError(func(err error, context map[string]interface{}) {
		errorType := "unknown"
		if ctxType, ok := context["error_type"]; ok {
			if et, ok := ctxType.(string); ok {
				errorType = et
			}
		}
		fmt.Printf("[错误] 类型: %s | 消息: %v\n", errorType, err)
	})

	// 创建上下文用于控制生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动SDK
	fmt.Println("启动CT日志监听器...")
	if err := sdk.Start(ctx); err != nil {
		log.Fatalf("启动SDK失败: %v", err)
	}

	fmt.Println("CT监听器已启动，正在监听子域名发现...")
	fmt.Println("按Ctrl+C停止监听")

	// 启动状态监控协程
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		// 用于计算增量QPS的变量
		var lastProcessedCount int64 = 0
		var lastReportTime = time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				currentTime := time.Now()
				status := sdk.GetStatus()
				stats := sdk.GetStats()

				// 计算总体平均QPS
				uptimeSeconds := status.Uptime.Seconds()
				avgQps := float64(0)
				if uptimeSeconds > 0 {
					avgQps = float64(status.ProcessedCertificates) / uptimeSeconds
				}

				// 计算增量QPS（最近30秒内的实际处理速度）
				timeDiff := currentTime.Sub(lastReportTime).Seconds()
				countDiff := status.ProcessedCertificates - lastProcessedCount
				incrementalQps := float64(0)
				if timeDiff > 0 {
					incrementalQps = float64(countDiff) / timeDiff
				}

				// 更新上一次的状态
				lastProcessedCount = status.ProcessedCertificates
				lastReportTime = currentTime

				fmt.Printf("\n=== 监听器状态报告 ===\n")
				fmt.Printf("运行时间: %v\n", status.Uptime)
				fmt.Printf("总处理证书数: %d\n", status.ProcessedCertificates)
				fmt.Printf("总发现子域名数: %d\n", status.DiscoveredSubdomains)
				fmt.Printf("平均QPS: %.2f certs/sec\n", avgQps)
				fmt.Printf("增量QPS: %.2f certs/sec (最近30秒)\n", incrementalQps)
				fmt.Printf("活跃监听器数: %d\n", status.ActiveWatchers)
				fmt.Printf("域名总数: %d\n", stats.TotalDomains)

				// 显示每个监听器的处理统计
				if len(status.WatcherStatuses) > 0 {
					fmt.Printf("\n--- 监听器详情 ---\n")
					totalWatcherProcessed := int64(0)
					for url, watcherStatus := range status.WatcherStatuses {
						fmt.Printf("监听器: %s\n", url)
						fmt.Printf("  状态: %s\n", watcherStatus.Status)
						fmt.Printf("  处理证书数: %d\n", watcherStatus.ProcessedCount)
						fmt.Printf("  最后索引: %d\n", watcherStatus.LastIndex)
						fmt.Printf("  错误次数: %d\n", watcherStatus.ErrorCount)
						totalWatcherProcessed += watcherStatus.ProcessedCount
					}
					fmt.Printf("\n监听器总处理证书数: %d\n", totalWatcherProcessed)
				}

				fmt.Printf("========================\n\n")
			}
		}
	}()

	// 动态添加域名示例
	go func() {
		time.Sleep(2 * time.Minute)
		fmt.Println("动态添加域名: test.com")
		if err := sdk.AddDomains([]string{"test.com"}); err != nil {
			fmt.Printf("添加域名失败: %v\n", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-sigChan
	fmt.Println("\n收到停止信号，正在停止CT监听器...")

	// 取消上下文，触发优雅关闭
	cancel()

	// 等待一段时间确保清理完成
	time.Sleep(5 * time.Second)

	fmt.Println("CT监听器已完全停止")
}

// 使用示例：
// func main() {
//     runCTExample()
// }
