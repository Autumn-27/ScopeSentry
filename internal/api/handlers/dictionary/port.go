package dictionary

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	dictservice "github.com/Autumn-27/ScopeSentry/internal/services/dictionary"
	"github.com/gin-gonic/gin"
)

var portService dictservice.PortService

func init() {
	portService = dictservice.NewPortService()
}

// PortDataRequest 请求体
// 用于 /data 分页与搜索
// swagger:ignore
// 仅用于绑定
type PortDataRequest struct {
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	Search    string `json:"search"`
}

// GetPortData 获取端口字典数据
// @Summary 获取端口字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param request body PortDataRequest true "请求参数"
// @Success 200 {object} response.Response{data=object{list=[]models.PortDoc,total=int}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/port/data [post]
func GetPortData(c *gin.Context) {
	var req PortDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := portService.Get(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total}, "api.success")
}

// UpgradePortDict 更新端口字典
// @Summary 更新端口字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param request body models.PortDoc true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/port/upgrade [post]
func UpgradePortDict(c *gin.Context) {
	var req models.UpdatePort
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if req.ID == "" {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("id is required"))
		return
	}
	if err := portService.Update(c, req.ID, req.Name, req.Value); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"code": 200, "message": "SensitiveRule updated successfully"}, "api.success")
}

// AddPortDict 新增端口字典
// @Summary 新增端口字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param request body models.Port true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/port/add [post]
func AddPortDict(c *gin.Context) {
	var req models.Port
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := portService.Add(c, req.Name, req.Value); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	response.Success(c, gin.H{"code": 200, "message": "Port dict added successfully"}, "api.success")
}

// DeletePortDict 删除端口字典
// @Summary 删除端口字典
// @Tags 字典管理
// @Accept json
// @Produce json
// @Param request body object{ids=[]string} true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security ApiKeyAuth
// @Router /api/dictionary/port/delete [post]
func DeletePortDict(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids cannot be empty"))
		return
	}
	if err := portService.Delete(c, req.IDs); err != nil {
		logger.Error(err.Error())
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"code": 200, "message": "Port dict deleted successfully"}, "api.success")
}
