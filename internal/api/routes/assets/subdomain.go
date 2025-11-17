// assets-------------------------------------
// @file      : subdomain.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/13 23:32
// -------------------------------------------

package assets

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/assets/subdomain"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
)

func registerSubDomainRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/subdomain",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     subdomain.GetSubdomains,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/taker",
				Handler:     subdomain.GetSubdomainTakerData,
				Middlewares: common.WithAuth(),
			},
		},
	}
}
