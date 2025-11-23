package assets

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/handlers/assets/asset"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func registerAssetRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/asset",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     asset.GetAssets,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/detail",
				Handler:     asset.GetAssetByID,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/changelog",
				Handler:     asset.GetChangeLog,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/card",
				Handler:     asset.GetAssetCardData,
				Middlewares: common.WithAuth(),
			},
		},
	}
} 