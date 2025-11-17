// Package node -----------------------------
// @file      : node.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/3 10:57
// -------------------------------------------
package node

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/node"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterNodeRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	nodeRoutes := models.RouteGroup{
		Path: "/node",
		Routes: []models.Route{
			// 需要认证和user:read权限的路由
			{
				Method:      "GET",
				Path:        "",
				Handler:     node.GetNode,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "GET",
				Path:        "/online",
				Handler:     node.GetNodeOnline,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/config/update",
				Handler:     node.ConfigUpdate,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     node.Delete,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/log",
				Handler:     node.GetLogs,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/plugin",
				Handler:     node.GetNodePlugin,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/restart",
				Handler:     node.RestartNode,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, nodeRoutes)
}
