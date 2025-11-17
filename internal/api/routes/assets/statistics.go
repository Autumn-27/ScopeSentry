package assets

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/assets/statistics"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
)

func registerStatisticsRoutes() models.RouteGroup {
	return models.RouteGroup{
		Path: "/statistics",
		Routes: []models.Route{
			{
				Method:      "GET",
				Path:        "",
				Handler:     statistics.GetAssetStatisticsData,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/port",
				Handler:     statistics.GetAssetPortStatistics,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/service",
				Handler:     statistics.GetAssetServiceStatistics,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/app",
				Handler:     statistics.GetAssetAppStatistics,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/icon",
				Handler:     statistics.GetAssetIconStatistics,
				Middlewares: common.WithAuth(),
			},
		},
	}
} 