// subdomain-------------------------------------
// @file      : taker.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/14 21:53
// -------------------------------------------

package subdomain

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/subdomain"
	"github.com/gin-gonic/gin"
)

var takerService subdomain.TakerService

func init() {
	takerService = subdomain.NewTakerService()
}

// GetSubdomainTakerData godoc
// @Summary 获取子域名接管数据
// @Description 获取子域名接管数据列表
// @Tags 子域名接管
// @Accept json
// @Produce json
// @Param request body subdomain.SubdomainTakerRequest true "请求参数"
// @Success 200 {object} response.Response{data=subdomain.SubdomainTakerResponse} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/assets/subdomain/taker [post]
func GetSubdomainTakerData(c *gin.Context) {
	var request models.SearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := takerService.GetSubdomainTakerData(c, request)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}
