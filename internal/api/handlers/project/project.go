package project

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/project"
	"github.com/gin-gonic/gin"
)

var projectService project.Service

func init() {
	projectService = project.NewService()
}

// GetProjectsByTag godoc
// @Summary 获取按标签分组的项目列表
// @Description 获取所有项目并按标签分组
// @Tags 项目
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]project.TagGroup}
// @Router /api/projects/all [get]
func GetProjectsByTag(c *gin.Context) {
	result, err := projectService.GetProjectsByTag(c)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, gin.H{
		"list": result,
	}, "")
}

// GetProjectsData 获取项目数据（含标签聚合）
// @Summary 获取项目数据
// @Description 返回每个标签下的分页项目列表与标签计数
// @Tags 项目
// @Accept json
// @Produce json
// @Param data body models.ProjectListRequest true "查询参数"
// @Success 200 {object} response.Response{data=models.ProjectListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/project/data [post]
func GetProjectsData(c *gin.Context) {
	var req models.ProjectListRequest
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
	data, err := projectService.GetProjectsData(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, data, "")
}

// GetProjectContent 获取项目内容
// @Summary 获取项目目标与基础信息
// @Tags 项目
// @Accept json
// @Produce json
// @Param data body models.ProjectContentRequest true "项目ID"
// @Success 200 {object} response.Response{data=models.ProjectContentResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/project/content [post]
func GetProjectContent(c *gin.Context) {
	var req models.ProjectContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	data, err := projectService.GetProjectContent(c, req.ID)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	if data == nil {
		response.NotFound(c, "api.not_found", nil)
		return
	}
	response.Success(c, data, "")
}

func AddProject(c *gin.Context) {
	var req models.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	err := projectService.AddProject(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
	return
}

func UpdatePorject(c *gin.Context) {
	var req models.UpdateProject
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	err := projectService.UpdateProject(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
	return
}

func DeleteProject(c *gin.Context) {
	var req models.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	err := projectService.DeleteProjects(c, req.IDs, req.DeleateAssets)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
}
