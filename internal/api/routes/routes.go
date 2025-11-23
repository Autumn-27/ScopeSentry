package routes

import (
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/api/middleware"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/assets"
	confRoutes "github.com/Autumn-27/ScopeSentry/internal/api/routes/configuration"
	dictRoutes "github.com/Autumn-27/ScopeSentry/internal/api/routes/dictionary"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/export"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/fingerprint"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/poc"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/project"
	sroutes "github.com/Autumn-27/ScopeSentry/internal/api/routes/sensitive_rule"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/system"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/task"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/user"
	"github.com/gin-gonic/gin"

	"github.com/Autumn-27/ScopeSentry/internal/api/routes/node"
)

// loggerWithSkipPaths 自定义日志中间件，支持跳过指定前缀的路径
func loggerWithSkipPaths(skipPrefixes []string) gin.HandlerFunc {
	// 创建默认的日志中间件
	defaultLogger := gin.Logger()

	return func(c *gin.Context) {
		// 检查路径是否需要跳过
		path := c.Request.URL.Path
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(path, prefix) {
				// 如果路径需要跳过，直接执行下一个中间件，不记录日志
				c.Next()
				return
			}
		}

		// 如果路径不需要跳过，执行默认的日志中间件
		defaultLogger(c)
	}
}

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	router := gin.New()

	// 配置日志中间件，跳过 /assets 和 /uploads 路径的日志
	router.Use(loggerWithSkipPaths([]string{"/assets", "/uploads", "/images"}))

	// 添加恢复中间件
	router.Use(gin.Recovery())

	// 添加全局中间件
	router.Use(middleware.I18nMiddleware())

	// 注册API路由
	api := router.Group("/api")
	{
		// 注册用户路由
		user.RegisterUserRoutes(api)

		// 注册system路由
		system.RegisterSystemRoutes(api)

		// 注册configuration路由
		confRoutes.RegisterConfigurationRoutes(api)

		// 注册asset路由
		assets.RegisterAssetsRoutes(api)

		// 注册fingerprint路由
		fingerprint.RegisterFingerprintRoutes(api)

		// 注册SensitiveRule路由
		sroutes.RegisterSensitiveRuleRoutes(api)

		// 注册节点路由
		node.RegisterNodeRoutes(api)

		task.RegisterTaskRoutes(api)

		project.RegisterProjectRoutes(api)

		poc.RegisterPocRoutes(api)

		export.RegisterExportRoutes(api)

		plugin.RegisterPluginRoutes(api)

		// 注册字典路由
		dictRoutes.RegisterDictionaryRoutes(api)
	}

	return router
}

// 中间件组合函数
