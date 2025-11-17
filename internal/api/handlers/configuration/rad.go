// configuration-------------------------------------
// @file      : rad.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/29 20:58
// -------------------------------------------

package configuration

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/configuration"
	"github.com/gin-gonic/gin"
)

var radConfigService *configuration.Service

// GetRadData @Summary 获取 Rad 配置
// @Tags        configuration
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} response.SuccessResponse{data=object{content=string}}
// @Router      /configuration/rad/data [get]
func GetRadData(c *gin.Context) {
	content, err := radConfigService.GetRadContent(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"content": content}, "api.success")
}

// SaveRadData @Summary 保存 Rad 配置
// @Tags        configuration
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       body body saveContentRequest true "配置内容"
// @Success     200 {object} response.SuccessResponse
// @Router      /configuration/rad/save [post]
func SaveRadData(c *gin.Context) {
	var req saveContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := radConfigService.SaveRadContent(c, req.Content); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	radConfigService = configuration.NewService()
}
