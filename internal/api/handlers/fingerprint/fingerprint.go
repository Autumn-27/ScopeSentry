// fingerprint-------------------------------------
// @file      : fingerprint.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/28 20:30
// -------------------------------------------

package fingerprint

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	service "github.com/Autumn-27/ScopeSentry-go/internal/services/fingerprint"
	"github.com/gin-gonic/gin"
)

var fingerprintService service.Service

type listRequest struct {
	Search    string `json:"search" binding:"omitempty" example:"nginx"`
	PageIndex int    `json:"pageIndex" binding:"required,min=1" example:"1"`
	PageSize  int    `json:"pageSize" binding:"required,min=1,max=100" example:"10"`
}

type updateRequest struct {
	ID             string `json:"id" binding:"required" example:"6642b0..."`
	Name           string `json:"name" binding:"required" example:"Nginx"`
	Rule           string `json:"rule" binding:"required" example:"body~=nginx && header~=server"`
	Category       string `json:"category" binding:"required" example:"web"`
	ParentCategory string `json:"parent_category" binding:"omitempty"`
	State          bool   `json:"state"`
}

type addRequest struct {
	Name           string `json:"name" binding:"required" example:"Nginx"`
	Rule           string `json:"rule" binding:"required" example:"body~=nginx && header~=server"`
	Category       string `json:"category" binding:"required" example:"web"`
	ParentCategory string `json:"parent_category" binding:"omitempty"`
	State          bool   `json:"state" binding:"required" example:"1"`
}

type deleteRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,dive,required"`
}

// Data @Summary      指纹规则查询
// @Description  按名称模糊搜索并分页
// @Tags         fingerprint
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      listRequest  true  "查询参数"
// @Success      200      {object}  response.SuccessResponse{data=object{list=[]models.FingerprintRule,total=int}}
// @Failure      400      {object}  response.BadRequestResponse
// @Failure      500      {object}  response.InternalServerErrorResponse
// @Router       /fingerprint/data [post]
func Data(c *gin.Context) {
	var req listRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	list, total, err := fingerprintService.List(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total}, "api.success")
}

// Update @Summary      更新指纹规则
// @Tags         fingerprint
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      updateRequest  true  "更新数据"
// @Success      200      {object}  response.SuccessResponse
// @Failure      400      {object}  response.BadRequestResponse
// @Failure      500      {object}  response.InternalServerErrorResponse
// @Router       /fingerprint/update [post]
func Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	data := models.FingerprintRule{
		Name:           req.Name,
		Rule:           req.Rule,
		Category:       req.Category,
		ParentCategory: req.ParentCategory,
		State:          req.State,
	}
	if err := fingerprintService.Update(c, req.ID, data); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// Add @Summary      新增指纹规则
// @Tags         fingerprint
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      addRequest  true  "新增数据"
// @Success      200      {object}  response.SuccessResponse{data=object{id=string}}
// @Failure      400      {object}  response.BadRequestResponse
// @Failure      500      {object}  response.InternalServerErrorResponse
// @Router       /fingerprint/add [post]
func Add(c *gin.Context) {
	var req addRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	data := models.FingerprintRule{
		Name:           req.Name,
		Rule:           req.Rule,
		Category:       req.Category,
		ParentCategory: req.ParentCategory,
		State:          req.State,
	}
	id, err := fingerprintService.Add(c, data)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"id": id}, "api.success")
}

// Delete @Summary      批量删除指纹规则
// @Tags         fingerprint
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      deleteRequest  true  "要删除的ID列表"
// @Success      200      {object}  response.SuccessResponse
// @Failure      400      {object}  response.BadRequestResponse
// @Failure      500      {object}  response.InternalServerErrorResponse
// @Router       /fingerprint/delete [post]
func Delete(c *gin.Context) {
	var req deleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	_, err := fingerprintService.Delete(c, req.IDs)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	fingerprintService = service.NewService()
}
