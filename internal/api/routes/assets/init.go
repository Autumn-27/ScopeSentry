// assets-------------------------------------
// @file      : asset.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/8 19:51
// -------------------------------------------

package assets

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterAssetsRoutes(api *gin.RouterGroup) {
	nodeRoutes := models.RouteGroup{
		Path: "/assets",
		Groups: []models.RouteGroup{
			registerAssetRoutes(),
			registerStatisticsRoutes(),
			registerRootDomainRoutes(),
			registerCommonRoutes(),
			registerSubDomainRoutes(),
			registerAppRoutes(),
			registerMpRoutes(),
			registerUrlRoutes(),
			registerIPRoutes(),
			registerCrawlerRoutes(),
			registerSensitivveRoutes(),
			registerDirscanRoutes(),
			registerVulnerabilityRoutes(),
			registerPageMonitoringRoutes(),
		},
	}

	common.RegisterRouteGroup(api, nodeRoutes)
}
