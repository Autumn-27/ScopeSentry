package middleware

import (
	"net/http"

	"github.com/Autumn-27/ScopeSentry-go/internal/constants"
	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		// 获取角色对应的权限列表
		permissions, exists := constants.RolePermissions[role.(string)]
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		// 检查是否具有所需权限
		hasPermission := false
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
} 