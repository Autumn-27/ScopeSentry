package role

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(r *gin.Engine) {
	// 创建API路由组
	//api := r.Group("/api")
	//{
	//	// 角色相关路由
	//	roleGroup := api.Group("/roles")
	//	{
	//		// 需要认证和角色管理权限的路由
	//		roleGroup.Use(middleware.AuthMiddleware())
	//		roleGroup.Use(middleware.PermissionMiddleware(constants.PermissionRoleRead))
	//		{
	//			roleGroup.GET("", role.List)
	//			roleGroup.GET("/:id", role.Get)
	//		}
	//
	//		// 需要创建角色权限的路由
	//		roleGroup.Use(middleware.PermissionMiddleware(constants.PermissionRoleCreate))
	//		{
	//			roleGroup.POST("", role.Create)
	//		}
	//
	//		// 需要更新角色权限的路由
	//		roleGroup.Use(middleware.PermissionMiddleware(constants.PermissionRoleUpdate))
	//		{
	//			roleGroup.PUT("/:id", role.Update)
	//		}
	//
	//		// 需要删除角色权限的路由
	//		roleGroup.Use(middleware.PermissionMiddleware(constants.PermissionRoleDelete))
	//		{
	//			roleGroup.DELETE("/:id", role.Delete)
	//		}
	//	}
	//}
}
