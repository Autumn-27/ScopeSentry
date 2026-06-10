// @title           ScopeSentry API
// @version         1.0
// @description     ScopeSentry 是一个安全扫描平台
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /api

// @securityDefinitions.apikey  ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description               Bearer token for authentication

package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/worker"

	"github.com/Autumn-27/ScopeSentry/internal/bootstrap"
	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/scheduler"
	"github.com/Autumn-27/ScopeSentry/internal/update"

	"github.com/Autumn-27/ScopeSentry/internal/logger"

	_ "github.com/Autumn-27/ScopeSentry/internal/database/mongodb"

	_ "github.com/Autumn-27/ScopeSentry/internal/database/redis"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/constants"

	_ "github.com/Autumn-27/ScopeSentry/internal/bootstrap"

	"github.com/Autumn-27/ScopeSentry/internal/api/routes"
	scopesentrymcp "github.com/Autumn-27/ScopeSentry/internal/mcp"
	redisLogSubscriber "github.com/Autumn-27/ScopeSentry/internal/services/redis_log_subscriber"

	"github.com/Autumn-27/ScopeSentry/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

//go:embed static/*
var embeddedFiles embed.FS

func main() {
	// 设置 Gin 模式
	if config.GlobalConfig.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	err := os.Setenv("TZ", config.GlobalConfig.System.Timezone)
	if err != nil {
		return
	}
	constants.Version = "1.9.1"
	fmt.Printf("version: %v\n", constants.Version)
	PrintPlugin()
	err = update.Update()
	if err != nil {
		logger.Error("Update failed", zap.Error(err))
	}
	bootstrap.GetProjectList()
	// 初始化计划任务
	scheduler.InitializeGlobalScheduler()
	scheduler.GetGlobalScheduler().Start()
	defer scheduler.GetGlobalScheduler().Stop()

	// 启动 Redis 日志订阅服务（在后台运行）
	go redisLogSubscriber.SubscribeLogChannel()
	// 创建路由
	router := routes.SetupRouter()
	scopesentrymcp.RegisterRoutes(router)
	gin.DisableConsoleColor()
	//router.Use(logger.GinLogger(), logger.GinRecovery())

	// 准备前端文件系统
	frontendFS, _ := fs.Sub(embeddedFiles, "static")
	// 准备静态资源文件系统（assets目录）
	assetsFS, _ := fs.Sub(embeddedFiles, "static/assets")
	// 准备其他静态资源文件系统
	cssFS, _ := fs.Sub(embeddedFiles, "static/css")
	dracoFS, _ := fs.Sub(embeddedFiles, "static/draco")
	jsFS, _ := fs.Sub(embeddedFiles, "static/js")
	libFS, _ := fs.Sub(embeddedFiles, "static/lib")
	pluginsFS, _ := fs.Sub(embeddedFiles, "static/plugins")
	pngFS, _ := fs.Sub(embeddedFiles, "static/png")

	// 注册 Swagger 路由（在静态文件之前，避免冲突）
	docs.SwaggerInfo.BasePath = "/api"
	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册静态资源路由（/assets 映射到 static/assets 目录）
	router.StaticFS("/assets", http.FS(assetsFS)) // /assets 对应前端 js/css

	// 注册其他静态资源路由
	router.StaticFS("/css", http.FS(cssFS))         // /css 对应 static/css 目录
	router.StaticFS("/draco", http.FS(dracoFS))     // /draco 对应 static/draco 目录
	router.StaticFS("/js", http.FS(jsFS))           // /js 对应 static/js 目录
	router.StaticFS("/lib", http.FS(libFS))         // /lib 对应 static/lib 目录
	router.StaticFS("/plugins", http.FS(pluginsFS)) // /plugins 对应 static/plugins 目录
	router.StaticFS("/png", http.FS(pngFS))         // /png 对应 static/png 目录

	// 注册 /map 路由，映射到 map.html
	router.GET("/map", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendFS, "map.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read map.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 注册 /favicon.ico 路由
	router.GET("/favicon.ico", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendFS, "favicon.ico")
		if err != nil {
			c.String(http.StatusNotFound, "Favicon not found")
			return
		}
		c.Data(http.StatusOK, "image/x-icon", data)
	})

	// 注册上传文件静态路由（/images 映射到外部 uploads 目录）
	router.Static("/images", config.GlobalConfig.System.ImgDir) // /uploads 对应上传的图片文件

	// 注册根路径路由，直接返回文件内容，避免重定向
	router.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(frontendFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 处理所有未匹配的路由，返回前端首页（用于SPA路由）
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 排除 API 路径和静态资源路径
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/mcp") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		// 如果请求的是静态资源但未找到，返回404而不是首页
		if strings.HasPrefix(path, "/assets") || strings.HasPrefix(path, "/uploads") ||
			strings.HasPrefix(path, "/css") || strings.HasPrefix(path, "/draco") ||
			strings.HasPrefix(path, "/js") || strings.HasPrefix(path, "/lib") ||
			strings.HasPrefix(path, "/plugins") || strings.HasPrefix(path, "/png") ||
			strings.HasPrefix(path, "/images") {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		// 其他路径返回前端首页（SPA路由）
		data, err := fs.ReadFile(frontendFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GlobalConfig.Server.Port),
		Handler: router,
	}

	// 运行资产处理
	worker.Run()

	// 在 goroutine 中启动服务器
	go func() {
		logger.Info("Starting server", zap.Int("port", config.GlobalConfig.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}

func Banner() {
	banner := "   _____                         _____            _              \n  / ____|                       / ____|          | |             \n | (___   ___ ___  _ __   ___  | (___   ___ _ __ | |_ _ __ _   _ \n  \\___ \\ / __/ _ \\| '_ \\ / _ \\  \\___ \\ / _ \\ '_ \\| __| '__| | | |\n  ____) | (_| (_) | |_) |  __/  ____) |  __/ | | | |_| |  | |_| |\n |_____/ \\___\\___/| .__/ \\___| |_____/ \\___|_| |_|\\__|_|   \\__, |\n                  | |                                       __/ |\n                  |_|                                      |___/ "
	fmt.Println(banner)
}

func PrintPlugin() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✨✨✨ IMPORTANT NOTICE: Please review the Plugin Key below ✨✨✨")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("🔑 Plugin Key: %s\n", config.GlobalConfig.System.PluginKey)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ Ensure the Plugin Key is correctly copied!")
}
