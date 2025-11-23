package app

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/app"
	"github.com/gin-gonic/gin"
)

var appService app.Service

func init() {
	appService = app.NewService()
}

// GetAppData godoc
// @Summary 获取应用数据
// @Description 获取应用数据列表
// @Tags 应用管理
// @Accept json
// @Produce json
// @Param request body models.SearchRequest true "请求参数"
// @Success 200 {object} response.Response{data=models.AppResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/assets/app [post]
func GetAppData(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := appService.GetAppData(c, req)
	if err != nil {
		logger.Error(fmt.Sprintf("Get app data error: %s", err.Error()))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}
