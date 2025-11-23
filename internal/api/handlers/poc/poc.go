package poc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/poc"
	"github.com/Autumn-27/ScopeSentry/internal/utils/random"
	"github.com/gin-gonic/gin"
)

var pocService poc.Service

func init() {
	pocService = poc.NewService()
}

// GetPocList godoc
// @Summary 获取POC列表
// @Description 获取POC列表，支持分页和搜索
// @Tags POC管理
// @Accept json
// @Produce json
// @Param request body poc.PocListRequest true "请求参数"
// @Success 200 {object} response.Response{data=poc.PocListResponse}
// @Router /poc [post]
func GetPocList(c *gin.Context) {
	var req models.PocListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(fmt.Sprintf("GetPocList bad request%v", err))
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := pocService.GetPocList(c, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("GetPocList error %v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}

// GetPocDetail godoc
// @Summary 获取POC详情
// @Description 根据ID获取POC详情
// @Tags POC管理
// @Accept json
// @Produce json
// @Param request body poc.PocDetailRequest true "请求参数"
// @Success 200 {object} response.Response{data=poc.Poc}
// @Router /poc/detail [post]
func GetPocDetail(c *gin.Context) {
	var req models.PocDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	result, err := pocService.GetPocDetail(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, result, "api.success")
}

// UpdatePoc godoc
// @Summary 更新POC
// @Description 更新POC信息
// @Tags POC管理
// @Accept json
// @Produce json
// @Param request body poc.PocUpdateRequest true "请求参数"
// @Success 200 {object} response.Response
// @Router /poc/update [post]
func UpdatePoc(c *gin.Context) {
	var req models.PocUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := pocService.UpdatePoc(c, &req); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// AddPoc godoc
// @Summary 添加POC
// @Description 添加新的POC
// @Tags POC管理
// @Accept json
// @Produce json
// @Param request body poc.PocAddRequest true "请求参数"
// @Success 200 {object} response.Response
// @Router /api/poc/add [post]
func AddPoc(c *gin.Context) {
	var req models.PocAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := pocService.AddPoc(c, &req); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// DeletePoc godoc
// @Summary 删除POC
// @Description 删除指定的POC
// @Tags POC管理
// @Accept json
// @Produce json
// @Param request body poc.PocDeleteRequest true "请求参数"
// @Success 200 {object} response.Response
// @Router /poc/delete [post]
func DeletePoc(c *gin.Context) {
	var req models.PocDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := pocService.DeletePoc(c, &req); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// GetAllPocData godoc
// @Summary 获取所有POC数据
// @Description 获取所有POC数据，包含id、name、time、tags字段
// @Tags POC管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /poc/data/all [get]
func GetAllPocData(c *gin.Context) {
	result, err := pocService.GetAllPocData(c)
	if err != nil {
		logger.Error(fmt.Sprintf("GetAllPocData error %v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, map[string]interface{}{
		"list": result,
	}, "api.success")
}

// ImportPoc godoc
// @Summary 导入POC文件
// @Description 上传ZIP文件并导入POC
// @Tags POC管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "POC文件(ZIP格式)"
// @Success 200 {object} response.Response{data=models.PocImportResponse}
// @Router /poc/data/import [post]
func ImportPoc(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(fmt.Sprintf("ImportPoc get file error %v", err))
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 检查文件扩展名
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("file must be zip format"))
		return
	}

	// 检查文件名安全性
	if strings.Contains(file.Filename, "..") || strings.Contains(file.Filename, "/") || strings.Contains(file.Filename, "\\") {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("unsafe filename"))
		return
	}

	// 使用系统临时目录
	tempDir := os.TempDir()
	fileName := random.GenerateString(5) + ".zip"
	filePath := filepath.Join(tempDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Error(fmt.Sprintf("ImportPoc save file error %v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	// 异步处理导入
	go func() {
		result, err := pocService.ImportPoc(c, filePath)
		if err != nil {
			logger.Error(fmt.Sprintf("ImportPoc process error %v", err))
		} else {
			logger.Info(fmt.Sprintf("ImportPoc success: %+v", result))
		}
	}()

	response.Success(c, map[string]interface{}{
		"message": "Importing in progress",
	}, "api.success")
}
