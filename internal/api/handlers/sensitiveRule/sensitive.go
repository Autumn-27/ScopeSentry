// sensitiveRule-------------------------------------
// @file      : sensitive.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/28 22:08
// -------------------------------------------

package sensitiveRule

import (
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	svc "github.com/Autumn-27/ScopeSentry/internal/services/sensitive_rule"
	"github.com/gin-gonic/gin"
)

var sensitiveService svc.Service

type listRequest struct {
	Search    string `json:"search" binding:"omitempty"`
	PageIndex int    `json:"pageIndex" binding:"required,min=1"`
	PageSize  int    `json:"pageSize" binding:"required,min=1,max=100"`
}

type updateRequest struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Regular string `json:"regular" binding:"required"`
	Color   string `json:"color" binding:"required"`
	State   bool   `json:"state" binding:"required"`
}

type addRequest struct {
	Name    string `json:"name" binding:"required"`
	Regular string `json:"regular" binding:"required"`
	Color   string `json:"color" binding:"required"`
	State   bool   `json:"state" binding:"required"`
}

type idsStateRequest struct {
	IDs   []string `json:"ids" `
	State bool     `json:"state"`
}

type idsRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,dive,required"`
}

// Data 列表
func Data(c *gin.Context) {
	var req listRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	list, total, err := sensitiveService.RuleList(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total}, "api.success")
}

// Update 更新
func Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := sensitiveService.RuleUpdate(c, req.ID, req.Name, req.Regular, req.Color, req.State); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// Add 新增
func Add(c *gin.Context) {
	var req addRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := sensitiveService.RuleAdd(c, req.Name, req.Regular, req.Color, req.State); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// UpdateState 批量改状态
func UpdateState(c *gin.Context) {
	var req idsStateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if _, err := sensitiveService.RuleUpdateState(c, req.IDs, req.State); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// Delete 批量删除
func Delete(c *gin.Context) {
	var req idsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if _, err := sensitiveService.RuleDelete(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	sensitiveService = svc.NewService()
}
