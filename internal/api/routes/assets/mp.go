// assets-------------------------------------
// @file      : mp.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/14 22:56
// -------------------------------------------

package assets

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/assets/mp"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
)

func registerMpRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/mp",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     mp.GetMPData,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
