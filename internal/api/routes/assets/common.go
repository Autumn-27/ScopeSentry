package assets

import (
	assetcommon "github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/assets/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
)

func registerCommonRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/common",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "/total",
				Handler:     assetcommon.TotalData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     assetcommon.DeleteData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete_tag",
				Handler:     assetcommon.DeleteTag,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/add_tag",
				Handler:     assetcommon.AddTag,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/update_status",
				Handler:     assetcommon.UpdateStatus,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
