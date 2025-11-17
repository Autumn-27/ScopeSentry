// task-------------------------------------
// @file      : scheduler.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/19 20:13
// -------------------------------------------

package task

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/task"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
)

func registerSchedulerTaskRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/scheduled",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     task.GetScheduledData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/add",
				Handler:     task.CreateScheduledTask,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/detail",
				Handler:     task.GetScheduledTaskDetail,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/update",
				Handler:     task.UpdateScheduledTask,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     task.DeleteScheduledTask,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/pagemonit/data",
				Handler:     task.GetPageMonitData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/pagemonit/add",
				Handler:     task.AddPageMonitTask,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/pagemonit/delete",
				Handler:     task.DeletePageMonitTask,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/pagemonit/update",
				Handler:     task.UpdatePageMonitScheduledTask,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
