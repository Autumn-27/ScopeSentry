package task

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/task/template"
	"github.com/gin-gonic/gin"
)

var templateService template.Service

func init() {
	templateService = template.NewService()
}

// ListRequest 模板列表请求
type ListRequest struct {
	PageIndex int    `json:"pageIndex" binding:"required" example:"1"`
	PageSize  int    `json:"pageSize" binding:"required" example:"10"`
	Query     string `json:"query" example:"test"`
}

// TemplateList 获取模板列表
// @Summary 获取模板列表
// @Description 获取扫描模板列表
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param request body ListRequest true "请求参数"
// @Success 200 {object} response.Response{data=task.TemplateList} "成功"
// @Router /api/task/template [post]
func TemplateList(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := templateService.List(c, req.PageIndex, req.PageSize, req.Query)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}

// DetailRequest 模板详情请求
type DetailRequest struct {
	ID string `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
}

// TemplateDetail 获取模板详情
// @Summary 获取模板详情
// @Description 获取扫描模板详情
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param request body DetailRequest true "请求参数"
// @Success 200 {object} response.Response{data=task.Template} "成功"
// @Router /api/task/template/detail [post]
func TemplateDetail(c *gin.Context) {
	var req DetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := templateService.Detail(c, req.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("TemplateDetail error", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}

// SaveRequest 保存模板请求
type SaveRequest struct {
	ID     string              `json:"id" example:"507f1f77bcf86cd799439011"`
	Result models.ScanTemplate `json:"result" binding:"required"`
}

// Save 保存模板
// @Summary 保存模板
// @Description 保存或更新扫描模板
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param request body SaveRequest true "请求参数"
// @Success 200 {object} response.Response "成功"
// @Router /api/task/template/save [post]
func Save(c *gin.Context) {
	var req SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	err := templateService.Save(c, req.ID, &req.Result)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// DeleteRequest 删除模板请求
type DeleteRequest struct {
	IDs []string `json:"ids" binding:"required" example:"['507f1f77bcf86cd799439011']"`
}

// Delete 删除模板
// @Summary 删除模板
// @Description 删除扫描模板
// @Tags 模板管理
// @Accept json
// @Produce json
// @Param request body DeleteRequest true "请求参数"
// @Success 200 {object} response.Response "成功"
// @Router /api/task/template/delete [post]
func Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	err := templateService.Delete(c, req.IDs)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}
