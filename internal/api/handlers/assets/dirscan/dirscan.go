// dirscan-------------------------------------
// @file      : dirscan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/17 15:06
// -------------------------------------------

package dirscan

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/dirscan"
	"github.com/gin-gonic/gin"
)

var dirscanService dirscan.Service

func init() {
	dirscanService = dirscan.NewService()
}

// List @Summary 获取目录扫描数据
// @Tags DirScan
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ListRequest true "分页与筛选参数"
// @Success 200 {object} response.DataResponse{data=[]models.DirScanResult}
// @Router /api/assets/dirscan [post]
func List(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	results, err := dirscanService.List(c, req)
	if err != nil {
		response.InternalServerError(c, "api.dirscan.list.failed", err)
		return
	}

	response.Success(c, gin.H{
		"list": results,
	}, "api.dirscan.list.success")
}
