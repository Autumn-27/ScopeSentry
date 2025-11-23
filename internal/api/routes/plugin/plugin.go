package plugin

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterPluginRoutes 注册插件相关路由
func RegisterPluginRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	pluginRoutes := models.RouteGroup{
		Path: "/plugin",
		Routes: []models.Route{
			// 需要认证和user:read权限的路由
			{
				Method:      "POST",
				Path:        "",
				Handler:     plugin.List,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/detail",
				Handler:     plugin.Detail,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/save",
				Handler:     plugin.Save,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     plugin.Delete,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/reinstall",
				Handler:     plugin.Reinstall,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/recheck",
				Handler:     plugin.Recheck,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/uninstall",
				Handler:     plugin.Uninstall,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/import",
				Handler:     plugin.Import,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/log",
				Handler:     plugin.GetLogs,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/log/clean",
				Handler:     plugin.CleanLogs,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/log/clean/all",
				Handler:     plugin.CleanAllLogs,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/module",
				Handler:     plugin.ListByModule,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/key/check",
				Handler:     plugin.CheckKey,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, pluginRoutes)
}
