// assets-------------------------------------
// @file      : app.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/14 22:41
// -------------------------------------------

package assets

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/assets/app"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func registerAppRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/app",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     app.GetAppData,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
