// export-------------------------------------
// @file      : export.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/25 15:12
// -------------------------------------------

package export

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/export"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterExportRoutes(api *gin.RouterGroup) {
	// 定义用户路由组
	nodeRoutes := models.RouteGroup{
		Path: "/export",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "",
				Handler:     export.Export,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/fields",
				Handler:     export.GetExportFields,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "GET",
				Path:        "/records",
				Handler:     export.GetExportRecords,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     export.DeleteExport,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "GET",
				Path:        "/download",
				Handler:     export.DownloadExport,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册用户路由组
	common.RegisterRouteGroup(api, nodeRoutes)
}
