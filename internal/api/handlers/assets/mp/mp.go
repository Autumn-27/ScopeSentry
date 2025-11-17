package mp

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/mp"
	"github.com/gin-gonic/gin"
)

var mpService mp.Service

func init() {
	mpService = mp.NewService()
}

// GetMPData godoc
// @Summary 获取小程序数据
// @Description 获取小程序数据列表
// @Tags 小程序管理
// @Accept json
// @Produce json
// @Param request body models.SearchRequest true "请求参数"
// @Success 200 {object} response.Response{data=models.MPResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/assets/mp [post]
func GetMPData(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := mpService.GetMPData(c, req)
	if err != nil {
		logger.Error(fmt.Sprintf("Get mp data error: %s", err.Error()))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}
