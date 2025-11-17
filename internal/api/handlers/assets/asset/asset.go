package asset

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/assets/asset"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var assetService asset.Service

func init() {
	assetService = asset.NewService()
}

// GetAssets godoc
// @Summary 获取资产列表
// @Description 获取资产列表，支持分页、排序和过滤
// @Tags 资产
// @Accept json
// @Produce json
// @Param query body models.SearchRequest true "查询参数"
// @Success 200 {object} response.SuccessResponse{data=object{list=[]models.Asset}}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/asset [post]
func GetAssets(c *gin.Context) {
	var query models.SearchRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	assets, err := assetService.GetAssets(c, query)
	if err != nil {
		response.InternalServerError(c, "api.asset.list.failed", err)
		return
	}

	response.Success(c, gin.H{
		"list": assets,
	}, "api.asset.list.success")
}

// GetAssetByID godoc
// @Summary 获取单个资产
// @Description 根据ID获取单个资产的详细信息
// @Tags 资产
// @Accept json
// @Produce json
// @Param id path string true "资产ID"
// @Success 200 {object} response.SuccessResponse{data=models.Asset}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 404 {object} response.NotFoundResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/asset/detail [get]
func GetAssetByID(c *gin.Context) {
	var query models.IdRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	asset, err := assetService.GetAssetByID(c, query.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		response.InternalServerError(c, "api.asset.get.failed", err)
		return
	}

	if asset == nil {
		response.NotFound(c, "api.asset.not_found", nil)
		return
	}

	response.Success(c, asset, "api.asset.get.success")
}

// CreateAsset godoc
// @Summary 创建资产
// @Description 创建新的资产
// @Tags 资产
// @Accept json
// @Produce json
// @Param asset body models.Asset true "资产信息"
// @Success 201 {object} response.SuccessResponse{data=models.Asset}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets [post]
func CreateAsset(c *gin.Context) {
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := assetService.CreateAsset(c, &asset); err != nil {
		response.InternalServerError(c, "api.asset.create.failed", err)
		return
	}

	response.Success(c, asset, "api.asset.create.success")
}

// UpdateAsset godoc
// @Summary 更新资产
// @Description 更新指定ID的资产信息
// @Tags 资产
// @Accept json
// @Produce json
// @Param id path string true "资产ID"
// @Param asset body models.Asset true "资产信息"
// @Success 200 {object} response.SuccessResponse{data=models.Asset}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 404 {object} response.NotFoundResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/{id} [put]
func UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := assetService.UpdateAsset(c, id, &asset); err != nil {
		response.InternalServerError(c, "api.asset.update.failed", err)
		return
	}

	response.Success(c, asset, "api.asset.update.success")
}

// DeleteAsset godoc
// @Summary 删除资产
// @Description 删除指定ID的资产
// @Tags 资产
// @Accept json
// @Produce json
// @Param id path string true "资产ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.BadRequestResponse
// @Failure 404 {object} response.NotFoundResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/{id} [delete]
func DeleteAsset(c *gin.Context) {
	id := c.Param("id")
	if err := assetService.DeleteAsset(c, id); err != nil {
		response.InternalServerError(c, "api.asset.delete.failed", err)
		return
	}

	response.Success(c, nil, "api.asset.delete.success")
}

// GetScreenshot godoc
// @Summary 获取资产截图
// @Description 获取指定ID资产的截图
// @Tags 资产
// @Accept json
// @Produce json
// @Param id path string true "资产ID"
// @Success 200 {object} response.SuccessResponse{data=string}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 404 {object} response.NotFoundResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/{id}/screenshot [get]
func GetScreenshot(c *gin.Context) {
	id := c.Param("id")
	screenshot, err := assetService.GetScreenshot(c, id)
	if err != nil {
		response.InternalServerError(c, "api.asset.screenshot.failed", err)
		return
	}

	if screenshot == "" {
		response.NotFound(c, "api.asset.screenshot.not_found", nil)
		return
	}

	response.Success(c, screenshot, "api.asset.screenshot.success")
}

// GetChangeLog godoc
// @Summary 获取资产变更日志
// @Description 获取指定ID资产的变更历史记录
// @Tags 资产
// @Accept json
// @Produce json
// @Param id path string true "资产ID"
// @Success 200 {object} response.SuccessResponse{data=[]models.AssetChangeLog}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 404 {object} response.NotFoundResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/asset/changelog [get]
func GetChangeLog(c *gin.Context) {
	var query models.IdRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	logs, err := assetService.GetChangeLog(c, query.ID)
	if err != nil {
		response.InternalServerError(c, "api.asset.changelog.failed", err)
		return
	}

	response.Success(c, logs, "")
}

// DeduplicateAssets godoc
// @Summary 资产去重
// @Description 对指定类型的资产进行去重处理
// @Tags 资产
// @Accept json
// @Produce json
// @Param assetType path string true "资产类型"
// @Param filter body map[string]interface{} true "过滤条件"
// @Param groupFields body []string true "分组字段"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/deduplicate/{assetType} [post]
func DeduplicateAssets(c *gin.Context) {
	assetType := c.Param("assetType")
	var request struct {
		Filter      bson.M   `json:"filter" binding:"required"`
		GroupFields []string `json:"groupFields" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := assetService.DeduplicateAssets(c, assetType, request.Filter, request.GroupFields); err != nil {
		response.InternalServerError(c, "api.asset.deduplicate.failed", err)
		return
	}

	response.Success(c, nil, "api.asset.deduplicate.success")
}

// GetAssetCardData godoc
// @Summary 获取资产卡片数据
// @Description 获取资产卡片数据列表，支持分页和搜索条件
// @Tags 资产
// @Accept json
// @Produce json
// @Param query body models.SearchRequest true "查询参数"
// @Success 200 {object} response.SuccessResponse{data=object{list=[]models.Asset}}
// @Failure 400 {object} response.BadRequestResponse
// @Failure 500 {object} response.InternalServerErrorResponse
// @Router /api/assets/asset/card [post]
func GetAssetCardData(c *gin.Context) {
	var query models.SearchRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	assets, err := assetService.GetAssetCardData(c, query)
	if err != nil {
		response.InternalServerError(c, "api.asset.card.failed", err)
		return
	}

	response.Success(c, gin.H{
		"list": assets,
	}, "api.asset.card.success")
}
