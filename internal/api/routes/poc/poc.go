package poc

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/poc"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

// RegisterPocRoutes 注册POC相关路由
func RegisterPocRoutes(api *gin.RouterGroup) {
	// 定义POC路由组
	pocRoutes := models.RouteGroup{
		Path: "/poc",
		Routes: []models.Route{
			// 需要认证和user:read权限的路由
			{
				Method:      "POST",
				Path:        "",
				Handler:     poc.GetPocList,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/detail",
				Handler:     poc.GetPocDetail,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/update",
				Handler:     poc.UpdatePoc,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/add",
				Handler:     poc.AddPoc,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     poc.DeletePoc,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "GET",
				Path:        "/data/all",
				Handler:     poc.GetAllPocData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/data/import",
				Handler:     poc.ImportPoc,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册POC路由组
	common.RegisterRouteGroup(api, pocRoutes)
}
