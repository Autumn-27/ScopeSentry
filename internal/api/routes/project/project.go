// project-------------------------------------
// @file      : project.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/13 21:22
// -------------------------------------------

package project

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/project"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	projectRoutes := models.RouteGroup{
		Path: "/project",
		Routes: []models.Route{
			// 公开路由
			{
				Method:      "GET",
				Path:        "/all",
				Handler:     project.GetProjectsByTag,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/data",
				Handler:     project.GetProjectsData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/content",
				Handler:     project.GetProjectContent,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/add",
				Handler:     project.AddProject,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     project.DeleteProject,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/update",
				Handler:     project.UpdatePorject,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, projectRoutes)
}
