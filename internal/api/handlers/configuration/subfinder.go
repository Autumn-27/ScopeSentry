// configuration-------------------------------------
// @file      : subfinder.go
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

var configService *configuration.Service

type saveContentRequest struct {
	Content string `json:"content" binding:"required" example:"your-config-here"`
}

// GetSubfinderData @Summary 获取 Subfinder 配置
// @Tags        configuration
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} response.SuccessResponse{data=object{content=string}}
// @Router      /configuration/subfinder/data [get]
func GetSubfinderData(c *gin.Context) {
	content, err := configService.GetSubfinderContent(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"content": content}, "api.success")
}

// SaveSubfinderData @Summary 保存 Subfinder 配置
// @Tags        configuration
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       body body saveContentRequest true "配置内容"
// @Success     200 {object} response.SuccessResponse
// @Router      /configuration/subfinder/save [post]
func SaveSubfinderData(c *gin.Context) {
	var req saveContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := configService.SaveSubfinderContent(c, req.Content); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	configService = configuration.NewService()
}
