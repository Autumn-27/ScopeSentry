package task

import (
	taskHandler "github.com/Autumn-27/ScopeSentry/internal/api/handlers/task"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func registerTemplateRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/template",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     taskHandler.TemplateList,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/detail",
				Handler:     taskHandler.TemplateDetail,
				Middlewares: common.WithAuth(),
			},
			{
				Method:  "POST",
				Path:    "/delete",
				Handler: taskHandler.Delete,
			},
			{
				Method:      "POST",
				Path:        "/save",
				Handler:     taskHandler.Save,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
