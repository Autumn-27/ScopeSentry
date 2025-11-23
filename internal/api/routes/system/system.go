// Package routes -----------------------------
// @file      : system.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/3 22:17
// -------------------------------------------
package system

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/system"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterSystemRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	systemRoutes := models.RouteGroup{
		Path: "/system",
		Routes: []models.Route{
			// 需要认证和user:read权限的路由
			{
				Method:      "GET",
				Path:        "version",
				Handler:     system.GetSystemVersion,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, systemRoutes)
}
