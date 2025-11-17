// configuration-------------------------------------
// @file      : system.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/29 20:57
// -------------------------------------------

package configuration

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/configuration"
	"github.com/gin-gonic/gin"
)

var systemService *configuration.Service

// GetSystemData @Summary 获取系统配置
// @Tags        configuration
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} response.SuccessResponse{data=object}
// @Router      /configuration/system/data [get]
func GetSystemData(c *gin.Context) {
	data, err := systemService.GetSystemData(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "api.success")
}

// SaveSystemData @Summary 保存系统配置
// @Tags        configuration
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       body body map[string]interface{} true "键值对"
// @Success     200 {object} response.SuccessResponse
// @Router      /configuration/system/save [post]
func SaveSystemData(c *gin.Context) {
	var kv map[string]interface{}
	if err := c.ShouldBindJSON(&kv); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := systemService.SaveSystemData(c, kv); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	systemService = configuration.NewService()
}
