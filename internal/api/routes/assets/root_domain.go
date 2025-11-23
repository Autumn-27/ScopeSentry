package assets

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/assets/root_domain"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func registerRootDomainRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/root_domain",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     root_domain.GetRootDomainData,
				Middlewares: common.WithAuth(),
			},
		},
	}
} 