package url

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/url"
	"github.com/gin-gonic/gin"
)

var urlService url.Service

func init() {
	urlService = url.NewService()
}

// GetURLs godoc
// @Summary 获取URL列表
// @Description 获取URL列表，支持分页、排序和过滤
// @Tags URL
// @Accept json
// @Produce json
// @Param query body models.SearchRequest true "查询参数"
// @Success 200 {object} response.SuccessResponse{data=object{list=[]models.URL}}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/url [post]
func GetURLs(c *gin.Context) {
	var query models.SearchRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	urls, err := urlService.GetURLs(c, query)
	if err != nil {
		response.InternalServerError(c, "api.url.list.failed", err)
		return
	}

	response.Success(c, gin.H{
		"list": urls,
	}, "api.url.list.success")
}
