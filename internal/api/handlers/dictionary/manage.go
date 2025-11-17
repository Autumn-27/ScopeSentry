package dictionary

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	dictservice "github.com/Autumn-27/ScopeSentry-go/internal/services/dictionary"
	"github.com/gin-gonic/gin"
)

var manageService dictservice.ManageService

func init() {
	manageService = dictservice.NewManageService()
}

// DeleteRequest 删除请求
type DeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// List 获取所有字典元数据
// @Summary 获取字典文件列表
// @Tags 字典管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object{list=[]DictionaryMeta}}
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/manage/list [get]
func List(c *gin.Context) {
	metas, err := manageService.List(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	resp := map[string]interface{}{
		"list": metas,
	}
	response.Success(c, resp, "api.success")
}

// Create 新增一个字典文件
// @Summary 新增字典文件
// @Tags 字典管理
// @Accept mpfd
// @Produce json
// @Param file formData file true "文件"
// @Param name formData string true "名称"
// @Param category formData string true "分类，如 dir、subdomain"
// @Success 200 {object} response.Response{data=object{message=string,code=int}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/manage/create [post]
func Create(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("missing file: %w", err))
		return
	}
	name := c.PostForm("name")
	category := c.PostForm("category")
	if name == "" || category == "" {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("name and category are required"))
		return
	}

	// 读取文件内容
	opened, err := fileHeader.Open()
	if err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	defer opened.Close()
	content, err := io.ReadAll(opened)
	if err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := manageService.Create(c, name, category, content); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"message": "file added successfully", "code": 200}, "api.success")
}

// Download 下载字典文件
// @Summary 下载字典文件
// @Tags 字典管理
// @Accept json
// @Produce application/octet-stream
// @Param id query string true "文件ID"
// @Security ApiKeyAuth
// @Router /api/dictionary/manage/download [get]
func Download(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("id is required"))
		return
	}

	filename, data, err := manageService.Download(c, id)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(200, "application/octet-stream", data)
}

// Delete 删除多个字典文件
// @Summary 删除字典文件
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param request body DeleteRequest true "请求参数"
// @Success 200 {object} response.Response{data=object{message=string,code=int}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/manage/delete [post]
func Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids cannot be empty"))
		return
	}

	if err := manageService.Delete(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"code": 200, "message": "delete file successful"}, "api.success")
}

// Save 替换指定字典文件内容
// @Summary 更新字典文件
// @Tags 字典管理
// @Accept mpfd
// @Produce json
// @Param id query string true "文件ID"
// @Param file formData file true "文件"
// @Success 200 {object} response.Response{data=object{message=string,code=int}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/manage/save [post]
func Save(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("id is required"))
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("missing file: %w", err))
		return
	}
	var content []byte
	if err := readFileContent(fileHeader, &content); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := manageService.Save(c, id, content); err != nil {
		logger.Error(err.Error())
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"code": 200, "message": "upload successful"}, "api.success")
}

func readFileContent(fileHeader *multipart.FileHeader, out *[]byte) error {
	f, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	*out = data
	return nil
}
