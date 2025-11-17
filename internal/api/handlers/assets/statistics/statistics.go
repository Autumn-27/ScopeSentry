package statistics

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/statistics"
	"github.com/gin-gonic/gin"
)

var statisticsService statistics.Service

func init() {
	statisticsService = statistics.NewService()
}

// GetAssetStatisticsData godoc
// @Summary 获取资产统计数据
// @Description 获取各类资产的统计数量
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=statistics.StatisticsData}
// @Router /api/assets/statistics [get]
func GetAssetStatisticsData(c *gin.Context) {
	data, err := statisticsService.GetStatisticsData(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetAssetPortStatistics godoc
// @Summary 获取端口统计
// @Description 获取资产端口分布统计
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body statistics.PortStatisticsRequest true "查询参数"
// @Success 200 {object} response.Response{data=statistics.PortStatisticsResponse}
// @Router /api/assets/statistics/port [post]
func GetAssetPortStatistics(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	data, err := statisticsService.GetPortStatistics(c, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("GetPortStatistics err:%v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetAssetTitleStatistics godoc
// @Summary 获取标题统计
// @Description 获取资产标题分布统计
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body statistics.TitleStatisticsRequest true "查询参数"
// @Success 200 {object} response.Response{data=statistics.TitleStatisticsResponse}
// @Router /api/v1/assets/statistics/title [post]
func GetAssetTitleStatistics(c *gin.Context) {
	var req models.TitleStatisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	data, err := statisticsService.GetTitleStatistics(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetAssetServiceStatistics godoc
// @Summary 获取服务类型统计
// @Description 获取资产服务类型分布统计
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body statistics.TypeStatisticsRequest true "查询参数"
// @Success 200 {object} response.Response{data=statistics.TypeStatisticsResponse}
// @Router /api/assets/statistics/service [post]
func GetAssetServiceStatistics(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	data, err := statisticsService.GetTypeStatistics(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetAssetIconStatistics godoc
// @Summary 获取图标统计
// @Description 获取资产图标分布统计
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body statistics.IconStatisticsRequest true "查询参数"
// @Success 200 {object} response.Response{data=statistics.IconStatisticsResponse}
// @Router /api/assets/statistics/icon [post]
func GetAssetIconStatistics(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	data, err := statisticsService.GetIconStatistics(c, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("GetIconStatistics err:%v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetAssetAppStatistics godoc
// @Summary 获取应用统计
// @Description 获取资产应用分布统计
// @Tags 资产统计
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body statistics.AppStatisticsRequest true "查询参数"
// @Success 200 {object} response.Response{data=statistics.AppStatisticsResponse}
// @Router /api/assets/statistics/app [post]
func GetAssetAppStatistics(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	data, err := statisticsService.GetAppStatistics(c, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("GetAppStatistics err:%v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}
