package crawler

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/crawler"
	"github.com/gin-gonic/gin"
)

var crawlerService crawler.Service

func init() {
	crawlerService = crawler.NewService()
}

// GetCrawlers godoc
// @Summary 获取爬虫任务列表
// @Description 获取爬虫任务列表，支持分页、排序和过滤
// @Tags Crawler
// @Accept json
// @Produce json
// @Param query body models.SearchRequest true "查询参数"
// @Success 200 {object} response.SuccessResponse{data=object{list=[]models.CrawlerResult}}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/crawler [post]
func GetCrawlers(c *gin.Context) {
	var query models.SearchRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	tasks, err := crawlerService.GetCrawlers(c, query)
	if err != nil {
		response.InternalServerError(c, "api.crawler.tasks.list.failed", err)
		return
	}

	response.Success(c, gin.H{
		"list": tasks,
	}, "api.crawler.list.success")
}
