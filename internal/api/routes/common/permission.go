// common-------------------------------------
// @file      : permission.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/8 21:11
// -------------------------------------------

package common

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// WithAuth 添加认证中间件
func WithAuth() []gin.HandlerFunc {
	return []gin.HandlerFunc{middleware.AuthMiddleware()}
}

// WithPermission 添加权限中间件
func WithPermission(permission string) []gin.HandlerFunc {
	return []gin.HandlerFunc{middleware.PermissionMiddleware(permission)}
}

// WithAuthAndPermission 添加认证和权限中间件
func WithAuthAndPermission(permission string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.AuthMiddleware(),
		middleware.PermissionMiddleware(permission),
	}
}
