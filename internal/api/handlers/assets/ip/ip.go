// ip-------------------------------------
// @file      : ip.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/11/15 21:17
// -------------------------------------------

package ip

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	ipservice "github.com/Autumn-27/ScopeSentry/internal/services/assets/ip"
	"github.com/gin-gonic/gin"
)

var ipService ipservice.Service

func init() {
	ipService = ipservice.NewService()
}

// GetIPAssets @Summary      获取IP资产列表
// @Description  根据查询条件获取IP资产信息
// @Tags         资产
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      models.SearchRequest  true  "查询条件"
// @Success      200      {object}  response.SuccessResponse{data=models.IPAssetResponse}
// @Failure      400      {object}  response.BadRequestResponse
// @Failure      500      {object}  response.InternalServerErrorResponse
// @Router       /api/assets/ip [post]
func GetIPAssets(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := ipService.GetIPAssets(c, req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}
