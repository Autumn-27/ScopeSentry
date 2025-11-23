// common-------------------------------------
// @file      : register.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/8 21:10
// -------------------------------------------

package common

import (
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterRouteGroup 注册路由组
func RegisterRouteGroup(group *gin.RouterGroup, routeGroup models.RouteGroup) {
	// 创建子路由组
	subGroup := group.Group(routeGroup.Path)

	// 应用路由组级别的中间件
	for _, handlerFunc := range routeGroup.Middlewares {
		subGroup.Use(handlerFunc)
	}

	// 注册路由
	for _, route := range routeGroup.Routes {
		// 创建路由处理器链
		handlers := make([]gin.HandlerFunc, 0)

		// 添加路由级别的中间件
		handlers = append(handlers, route.Middlewares...)

		// 添加路由处理器
		handlers = append(handlers, route.Handler)

		// 根据HTTP方法注册路由
		switch route.Method {
		case "GET":
			subGroup.GET(route.Path, handlers...)
		case "POST":
			subGroup.POST(route.Path, handlers...)
		case "PUT":
			subGroup.PUT(route.Path, handlers...)
		case "DELETE":
			subGroup.DELETE(route.Path, handlers...)
		case "PATCH":
			subGroup.PATCH(route.Path, handlers...)
		}
	}

	for _, childGroup := range routeGroup.Groups {
		RegisterRouteGroup(subGroup, childGroup)
	}
}
