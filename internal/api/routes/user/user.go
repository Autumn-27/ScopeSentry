package user

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/user"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	userRoutes := models.RouteGroup{
		Path: "/user",
		Routes: []models.Route{
			// 公开路由
			{
				Method:  "POST",
				Path:    "/login",
				Handler: user.Login,
			},
			// 修改密码（需要认证）
			{
				Method:      "POST",
				Path:        "/changePassword",
				Handler:     user.ChangePassword,
				Middlewares: common.WithAuth(),
			},
			//{
			//	Method:  "POST",
			//	Path:    "/register",
			//	Handler: user.Register,
			//},
			//// 需要认证和user:read权限的路由
			//{
			//	Method:      "GET",
			//	Path:        "",
			//	Handler:     user.List,
			//	Middlewares: common.WithAuthAndPermission(constants.PermissionUserRead),
			//},
			//{
			//	Method:      "GET",
			//	Path:        "/:id",
			//	Handler:     user.Get,
			//	Middlewares: common.WithAuthAndPermission(constants.PermissionUserRead),
			//},
			//// 需要认证和user:update权限的路由
			//{
			//	Method:      "PUT",
			//	Path:        "/:id",
			//	Handler:     user.Update,
			//	Middlewares: common.WithAuthAndPermission(constants.PermissionUserUpdate),
			//},
			//// 需要认证和user:delete权限的路由
			//{
			//	Method:      "DELETE",
			//	Path:        "/:id",
			//	Handler:     user.Delete,
			//	Middlewares: common.WithAuthAndPermission(constants.PermissionUserDelete),
			//},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, userRoutes)
}
