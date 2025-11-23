package root_domain

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/root_domain"
	"github.com/gin-gonic/gin"
)

var rootDomainService root_domain.Service

func init() {
	rootDomainService = root_domain.NewService()
}

// GetRootDomainData godoc
// @Summary 获取根域名数据
// @Description 分页获取根域名数据列表
// @Tags 根域名管理
// @Accept json
// @Produce json
// @Param request body GetRootDomainDataRequest true "请求参数"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Security ApiKeyAuth
// @Router /api/assets/root_domain [post]
func GetRootDomainData(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := rootDomainService.GetRootDomainData(c, req)
	if err != nil {
		logger.Error(fmt.Sprintf("Get root domain data error: %s", err.Error()))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}
