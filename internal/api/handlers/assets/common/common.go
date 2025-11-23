// common-------------------------------------
// @file      : common.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/10 23:33
// -------------------------------------------

package common

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"github.com/gin-gonic/gin"
)

var commonService common.Service

func init() {
	commonService = common.NewService()
}

// DeleteData 删除数据
// @Summary 删除数据
// @Description 根据ID列表删除指定集合中的数据
// @Tags 通用操作
// @Accept json
// @Produce json
// @Param request body common.DeleteRequest true "删除请求参数"
// @Success 200 {object} response.Response
// @Router /api/asset/common/delete [post]
func DeleteData(c *gin.Context) {
	var req models.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := commonService.DeleteData(c, &req); err != nil {
		switch err {
		case common.ErrInvalidCollection:
			response.NotFound(c, "api.not_found", err)
		case common.ErrNoValidIDs, common.ErrNoDocumentsDeleted:
			response.NotFound(c, "api.not_found", err)
		default:
			response.InternalServerError(c, "api.error", err)
		}
		return
	}

	response.Success(c, nil, "api.success")
}

// AddTag 添加标签
// @Summary 添加标签
// @Description 为指定文档添加标签
// @Tags 通用操作
// @Accept json
// @Produce json
// @Param request body common.TagRequest true "标签请求参数"
// @Success 200 {object} response.Response
// @Router /api/assets/common/add_tag [post]
func AddTag(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := commonService.AddTag(c, &req); err != nil {
		switch err {
		case common.ErrInvalidRequest, common.ErrInvalidID:
			response.NotFound(c, "api.not_found", err)
		default:
			logger.Error(fmt.Sprintf("add_tag error: %v", err))
			response.InternalServerError(c, "api.error", err)
		}
		return
	}

	response.Success(c, nil, "api.success")
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 从指定文档中删除标签
// @Tags 通用操作
// @Accept json
// @Produce json
// @Param request body common.TagRequest true "标签请求参数"
// @Success 200 {object} response.Response
// @Router /api/common/delete_tag [post]
func DeleteTag(c *gin.Context) {
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := commonService.DeleteTag(c, &req); err != nil {
		switch err {
		case common.ErrInvalidRequest, common.ErrInvalidID:
			response.NotFound(c, "api.not_found", err)
		default:
			response.InternalServerError(c, "api.error", err)
		}
		return
	}

	response.Success(c, nil, "api.success")
}

// UpdateStatus 更新状态
// @Summary 更新状态
// @Description 更新指定文档的状态
// @Tags 通用操作
// @Accept json
// @Produce json
// @Param request body common.StatusRequest true "状态请求参数"
// @Success 200 {object} response.Response
// @Router /api/common/update_status [post]
func UpdateStatus(c *gin.Context) {
	var req models.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := commonService.UpdateStatus(c, &req); err != nil {
		switch err {
		case common.ErrInvalidRequest, common.ErrInvalidID:
			response.NotFound(c, "api.not_found", err)
		default:
			response.InternalServerError(c, "api.error", err)
		}
		return
	}

	response.Success(c, nil, "api.success")
}

// TotalData 获取总数
// @Summary 获取总数
// @Description 获取指定集合的数据总数
// @Tags 通用操作
// @Accept json
// @Produce json
// @Param request body common.TotalRequest true "总数请求参数"
// @Success 200 {object} response.Response
// @Router /api/assets/common/total [post]
func TotalData(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	count, err := commonService.TotalData(c, &req)
	if err != nil {
		switch err {
		case common.ErrInvalidCollection:
			response.NotFound(c, "api.not_found", err)
		default:
			response.InternalServerError(c, "api.error", err)
		}
		return
	}

	response.Success(c, gin.H{"total": count}, "api.success")
}
