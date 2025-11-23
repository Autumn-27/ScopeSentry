package dictionary

import (
	dictHandler "github.com/Autumn-27/ScopeSentry/internal/api/handlers/dictionary"
	"github.com/Autumn-27/ScopeSentry/internal/api/routes/common"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

func RegisterDictionaryRoutes(api *gin.RouterGroup) {
	group := models.RouteGroup{
		Path: "/dictionary",
		Groups: []models.RouteGroup{
			{
				Path: "/manage",
				Routes: []models.Route{
					{Method: "GET", Path: "/list", Handler: dictHandler.List, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/create", Handler: dictHandler.Create, Middlewares: common.WithAuth()},
					{Method: "GET", Path: "/download", Handler: dictHandler.Download, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/delete", Handler: dictHandler.Delete, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/save", Handler: dictHandler.Save, Middlewares: common.WithAuth()},
				},
			},
			{
				Path: "/port",
				Routes: []models.Route{
					{Method: "POST", Path: "/data", Handler: dictHandler.GetPortData, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/upgrade", Handler: dictHandler.UpgradePortDict, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/add", Handler: dictHandler.AddPortDict, Middlewares: common.WithAuth()},
					{Method: "POST", Path: "/delete", Handler: dictHandler.DeletePortDict, Middlewares: common.WithAuth()},
				},
			},
		},
	}
	common.RegisterRouteGroup(api, group)
}
