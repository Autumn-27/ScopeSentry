// export-------------------------------------
// @file      : export.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/25 14:34
// -------------------------------------------

package export

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/export"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"path/filepath"
)

var exportService export.Service

func init() {
	exportService = export.NewService()
}

// Export 处理导出请求
// @Summary 导出数据
// @Description 根据请求参数导出数据到文件
// @Tags export
// @Accept json
// @Produce json
// @Param request body ExportRequest true "导出请求参数"
// @Success 200 {object} ExportResponse
// @Router /api/export [post]
func Export(c *gin.Context) {
	var req models.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 验证必要参数
	if req.Index == "" || req.Quantity == 0 || req.Type == "" {
		response.BadRequest(c, "api.invalid_params", fmt.Errorf("index, quantity, and type are required"))
		return
	}

	// 验证文件类型
	if req.FileType != "xlsx" && req.FileType != "json" {
		response.BadRequest(c, "api.invalid_file_type", fmt.Errorf("filetype must be csv or json"))
		return
	}

	// 验证字段
	if len(req.Field) == 0 {
		response.BadRequest(c, "api.invalid_fields", fmt.Errorf("fields cannot be empty"))
		return
	}
	searchQuery, err := helper.GetSearchQuery(models.SearchRequest{
		Index:            req.Index,
		SearchExpression: req.SearchExpression,
		Filter:           req.Filter,
		FuzzyQuery:       req.FuzzyQuery,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("get search query error:%s", err.Error()))
		response.InternalServerError(c, "api.export_failed", err)
		return
	}
	if req.Index == "PageMonitoring" {
		searchQuery["hash"] = map[string]interface{}{"$size": 2}
	}
	filter := bson.M(searchQuery)

	// 生成文件名
	fileName := fmt.Sprintf("%s.%s", uuid.New().String(), req.FileType)

	// 创建导出任务
	err = exportService.CreateExportTask(c, &models.ExportTask{
		FileName:   fileName,
		CreateTime: helper.GetNowTimeString(),
		Quantity:   req.Quantity,
		DataType:   req.Index,
		FileType:   req.FileType,
		State:      0,
	})

	if err != nil {
		response.InternalServerError(c, "api.export_failed", err)
		return
	}
	run := models.ExportRun{
		Filter:   filter,
		FileName: fileName,
		Index:    req.Index,
		FileType: req.FileType,
		Field:    req.Field,
		Quantity: req.Quantity,
	}
	// 异步处理导出任务
	go exportService.ProcessExportTask(c, run)

	response.Success(c, models.ExportResponse{
		Message: "Successfully added data export task",
		Code:    200,
	}, "api.success")
}

// GetExportFields 获取导出字段
// @Summary 获取导出字段
// @Description 根据索引获取可导出的字段列表
// @Tags export
// @Accept json
// @Produce json
// @Param request body map[string]string true "请求参数"
// @Success 200 {object} map[string]interface{}
// @Router /api/export/fields [post]
func GetExportFields(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	index := req["index"]
	fields, err := exportService.GetExportFields(c, index)
	if err != nil {
		response.InternalServerError(c, "api.get_fields_failed", err)
		return
	}

	response.Success(c, map[string]interface{}{
		"field": fields,
	}, "api.success")
}

// GetExportRecords 获取导出记录
// @Summary 获取导出记录
// @Description 获取所有导出任务的记录
// @Tags export
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/export/records [get]
func GetExportRecords(c *gin.Context) {
	records, err := exportService.GetExportRecords(c)
	if err != nil {
		response.InternalServerError(c, "api.get_records_failed", err)
		return
	}

	response.Success(c, map[string]interface{}{
		"list": records,
	}, "api.success")
}

// DeleteExport 删除导出文件
// @Summary 删除导出文件
// @Description 删除指定的导出文件
// @Tags export
// @Accept json
// @Produce json
// @Param request body map[string][]string true "请求参数"
// @Success 200 {object} ExportResponse
// @Router /api/export/delete [post]
func DeleteExport(c *gin.Context) {
	var req map[string][]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	ids := req["ids"]
	if len(ids) == 0 {
		response.NotFound(c, "api.no_files_selected", nil)
		return
	}

	err := exportService.DeleteExportFiles(c, ids)
	if err != nil {
		response.InternalServerError(c, "api.delete_failed", err)
		return
	}

	response.Success(c, models.ExportResponse{
		Message: "Export files deleted successfully",
		Code:    200,
	}, "api.success")
}

// DownloadExport 下载导出文件
// @Summary 下载导出文件
// @Description 下载指定的导出文件
// @Tags export
// @Produce application/octet-stream
// @Param filename query string true "文件名"
// @Success 200 {file} binary
// @Router /api/export/download [get]
func DownloadExport(c *gin.Context) {
	fileName := c.Query("filename")
	if !exportService.IsValidFileName(fileName) {
		response.BadRequest(c, "api.invalid_filename", nil)
		return
	}
	filePath := filepath.Join(config.GlobalConfig.System.ExeDir, "files", "export", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "api.file_not_found", nil)
		return
	}

	c.File(filePath)
}
