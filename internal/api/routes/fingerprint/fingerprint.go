// fingerprint-------------------------------------
// @file      : fingerprint.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/28 20:38
// -------------------------------------------

package fingerprint

import (
	handler "github.com/Autumn-27/ScopeSentry-go/internal/api/handlers/fingerprint"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterFingerprintRoutes(api *gin.RouterGroup) {
	// 定义指纹路由组
	fingerprintRoutes := models.RouteGroup{
		Path: "/fingerprint",
		Routes: []models.Route{
			{
				Method:      "POST",
				Path:        "/data",
				Handler:     handler.Data,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/update",
				Handler:     handler.Update,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/add",
				Handler:     handler.Add,
				Middlewares: common.WithAuth(),
			},
			{
				Method:      "POST",
				Path:        "/delete",
				Handler:     handler.Delete,
				Middlewares: common.WithAuth(),
			},
		},
	}

	// 注册指纹路由组
	common.RegisterRouteGroup(api, fingerprintRoutes)
}
