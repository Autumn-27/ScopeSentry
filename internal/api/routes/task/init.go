// task-------------------------------------
// @file      : task.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/10 15:33
// -------------------------------------------

package task

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	taskRoutes := models.RouteGroup{
		Path: "/task",
		Groups: []models.RouteGroup{
			registerTaskRoutes(),
			registerTemplateRoutes(),
			registerSchedulerTaskRoutes(),
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, taskRoutes)
}
