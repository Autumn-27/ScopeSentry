package page_monitoring

import (
	"errors"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/page_monitoring"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/gin-gonic/gin"
)

var pageMonitoringService page_monitoring.Service

func init() {
	pageMonitoringService = page_monitoring.NewService()
}

// GetResult @Summary      获取页面监控结果
// @Description  获取页面监控结果列表接口
// @Tags         page_monitoring
// @Accept       json
// @Produce      json
// @Param        resultRequest  body      models.SearchRequest  true  "结果请求"
// @Success      200  {object}  response.SuccessResponse{data=models.PageMonitoringResultResponse}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/assets/page-monitoring [post]
func GetResult(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := pageMonitoringService.GetResult(c, req)
	if err != nil {
		logger.Error(fmt.Sprintf("page-monitoring GetResult error: %s", err))
		response.InternalServerError(c, "api.page_monitoring.result.error", err)
		return
	}

	response.Success(c, result, "api.page_monitoring.result.success")
}

// GetDiff @Summary      获取页面监控差异
// @Description  获取页面监控差异接口
// @Tags         page_monitoring
// @Accept       json
// @Produce      json
// @Param        diffRequest  body      models.PageMonitoringDiffRequest  true  "差异请求"
// @Success      200  {object}  response.SuccessResponse{data=models.PageMonitoringBody}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /page-monitoring/diff [post]
func GetDiff(c *gin.Context) {
	var req models.PageMonitoringDiffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if req.ID == "" {
		response.BadRequest(c, "api.page_monitoring.id_missing", nil)
		return
	}

	result, err := pageMonitoringService.GetDiff(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, page_monitoring.ErrPageNotFound) {
			response.NotFound(c, "api.page_monitoring.diff.not_found", err)
			return
		}
		logger.Error(fmt.Sprintf("page-monitoring GetDiff error: %s", err))
		response.InternalServerError(c, "api.page_monitoring.diff.error", err)
		return
	}

	response.Success(c, result, "api.page_monitoring.diff.success")
}

// GetHistory @Summary      获取页面监控历史记录
// @Description  获取页面监控历史记录接口
// @Tags         page_monitoring
// @Accept       json
// @Produce      json
// @Param        historyRequest  body      models.PageMonitoringHistoryRequest  true  "历史记录请求"
// @Success      200  {object}  response.SuccessResponse{data=[]models.PageMonitoring}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /page-monitoring/history [post]
func GetHistory(c *gin.Context) {
	var req models.PageMonitoringHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if req.ID == "" {
		response.BadRequest(c, "api.page_monitoring.id_missing", nil)
		return
	}

	result, err := pageMonitoringService.GetHistory(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, page_monitoring.ErrPageNotFound) {
			response.NotFound(c, "api.page_monitoring.history.not_found", err)
			return
		}
		response.InternalServerError(c, "api.page_monitoring.history.error", err)
		return
	}

	response.Success(c, result, "api.page_monitoring.history.success")
}

// GetContent @Summary      获取页面监控内容
// @Description  获取页面监控内容接口
// @Tags         page_monitoring
// @Accept       json
// @Produce      json
// @Param        contentRequest  body      models.PageMonitoringContentRequest  true  "内容请求"
// @Success      200  {object}  response.SuccessResponse{data=models.PageMonitoringBody}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /page-monitoring/content [post]
func GetContent(c *gin.Context) {
	var req models.PageMonitoringContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if req.ID == "" {
		response.BadRequest(c, "api.page_monitoring.id_missing", nil)
		return
	}

	if req.Flag != "current" && req.Flag != "previous" {
		response.BadRequest(c, "api.page_monitoring.invalid_flag", nil)
		return
	}

	result, err := pageMonitoringService.GetContent(c.Request.Context(), req.ID, req.Flag)
	if err != nil {
		if err == page_monitoring.ErrPageNotFound {
			response.NotFound(c, "api.page_monitoring.content.not_found", err)
			return
		}
		response.InternalServerError(c, "api.page_monitoring.content.error", err)
		return
	}

	response.Success(c, result, "api.page_monitoring.content.success")
}
