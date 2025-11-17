// Package node -----------------------------
// @file      : node.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/30 15:36
// -------------------------------------------
package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/node"
	"github.com/gin-gonic/gin"
)

var nodeService node.Service

// GetNode @Summary      获取节点数据
// @Description  获取所有在线节点的数据
// @Tags         node
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]models.NodeData}}
// @Failure      401  {object}  response.UnauthorizedResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /node [get]
func GetNode(c *gin.Context) {
	result, err := nodeService.GetNodeData(c.Request.Context(), false)
	if err != nil {
		response.InternalServerError(c, "api.node.get_data.failed", err)
		return
	}

	response.Success(c, response.DataResponse{List: result}, "")
}

// GetNodeOnline @Summary      获取在线节点
// @Description  获取所有在线节点的数据
// @Tags         node
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]string}}
// @Failure      401  {object}  response.UnauthorizedResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /node/online [get]
func GetNodeOnline(c *gin.Context) {
	result, err := nodeService.GetNodeData(c.Request.Context(), true)
	if err != nil {
		response.InternalServerError(c, "api.node.get_data.failed", err)
		return
	}

	// 仅返回 name 字段
	names := make([]string, 0, len(result))
	for _, item := range result {
		names = append(names, item.Name)
	}

	response.Success(c, response.DataResponse{List: names}, "")
}

// NodeConfigUpdateRequest 节点配置更新请求
type NodeConfigUpdateRequest struct {
	Name          string `json:"name" binding:"required"`
	OldName       string `json:"oldName" binding:"omitempty"`
	ModulesConfig string `json:"ModulesConfig" binding:"required"`
	State         bool   `json:"state" binding:"required"`
}

// ConfigUpdate 更新节点配置
// @Summary 更新节点配置
// @Tags node
// @Accept json
// @Produce json
// @Param data body NodeConfigUpdateRequest true "节点配置"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/node/config/update [post]
func ConfigUpdate(c *gin.Context) {
	var req NodeConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	old := req.OldName
	if strings.TrimSpace(old) == "" {
		old = req.Name
	}
	msg := models.Message{
		Name:    old,
		Type:    "nodeConfig",
		Content: req.Name + "[*]" + fmt.Sprintf("%v", req.State) + "[*]" + req.ModulesConfig,
	}
	if err := nodeService.RefreshConfig(c, msg); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// NodeDeleteRequest 删除节点请求
type NodeDeleteRequest struct {
	Names []string `json:"names" binding:"required"`
}

// Delete 删除节点
// @Summary 删除节点
// @Tags node
// @Accept json
// @Produce json
// @Param data body NodeDeleteRequest true "节点名称列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/node/delete [post]
func Delete(c *gin.Context) {
	var req NodeDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.Names) == 0 {
		response.BadRequest(c, "api.bad_request", nil)
		return
	}
	ctx := context.Background()
	if err := nodeService.DeleteNodes(ctx, req.Names); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// NodeLogRequest 节点日志请求
type NodeLogRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetLogs 获取节点日志
// @Summary 获取节点日志
// @Tags node
// @Accept json
// @Produce json
// @Param data body NodeLogRequest true "节点名称"
// @Success 200 {object} response.Response{data=object{logs=string}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/node/log [post]
func GetLogs(c *gin.Context) {
	var req NodeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	ctx := context.Background()
	logs, err := nodeService.GetNodeLogs(ctx, req.Name)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, gin.H{"logs": logs}, "")
}

// GetNodePlugin 获取节点插件信息
// @Summary 获取节点插件信息
// @Tags node
// @Accept json
// @Produce json
// @Param data body NodeLogRequest true "节点名称"
// @Success 200 {object} response.Response{data=object{list=[]models.NodePluginInfo}}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/node/plugin [post]
func GetNodePlugin(c *gin.Context) {
	var req NodeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	ctx := context.Background()
	plugins, err := nodeService.GetNodePlugin(ctx, req.Name)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, response.DataResponse{List: plugins}, "")
}

// RestartNode 重启节点
// @Summary 重启节点
// @Tags node
// @Accept json
// @Produce json
// @Param data body NodeLogRequest true "节点名称"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/node/restart [post]
func RestartNode(c *gin.Context) {
	var req NodeLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	ctx := context.Background()
	if err := nodeService.RestartNode(ctx, req.Name); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

func init() {
	nodeService = node.NewService()
}
