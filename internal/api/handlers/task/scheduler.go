// task-------------------------------------
// @file      : scheduler.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/6/8 22:58
// -------------------------------------------

package task

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/scheduler"

	"github.com/gin-gonic/gin"
)

var SchedulerService scheduler.Service

func init() {
	SchedulerService = scheduler.NewService()
}

// GetScheduledData 获取计划任务列表
// @Summary      获取计划任务列表
// @Description  支持name字段模糊搜索与分页
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.ScheduledTaskListRequest  true  "查询参数"
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]models.Task,total=int}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled [post]
func GetScheduledData(c *gin.Context) {
	var req models.ScheduledTaskListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 设置默认分页参数
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	tasks, total, err := SchedulerService.GetScheduledData(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	resp := map[string]interface{}{
		"list":  tasks,
		"total": total,
	}
	response.Success(c, resp, "")
}

// CreateScheduledTask 创建计划任务
// @Summary      创建计划任务
// @Description  创建新的计划任务
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.ScheduledTaskAddRequest  true  "任务信息"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/add [post]
func CreateScheduledTask(c *gin.Context) {
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 检查任务名是否已存在
	exists, err := SchedulerService.CheckScheduledTaskNameExists(c, req.Name)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	if exists {
		response.BadRequest(c, "api.task.name_exists", nil)
		return
	}

	// 验证必要参数
	if req.Name == "" || len(req.Node) == 0 {
		response.BadRequest(c, "api.task.target_null", nil)
		return
	}

	if err := SchedulerService.CreateScheduledTask(c, &req); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// UpdateScheduledTask 更新计划任务
// @Summary      更新计划任务
// @Description  更新指定ID的计划任务
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.ScheduledTaskUpdateRequest  true  "任务信息"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/update [post]
func UpdateScheduledTask(c *gin.Context) {
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if err := SchedulerService.UpdateScheduledTask(c, &req); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// DeleteScheduledTask 删除计划任务
// @Summary      删除计划任务
// @Description  批量删除计划任务
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.ScheduledTaskDeleteRequest  true  "任务ID列表"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/delete [post]
func DeleteScheduledTask(c *gin.Context) {
	var req models.ScheduledTaskDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids is required"))
		return
	}

	if err := SchedulerService.DeleteScheduledTasks(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// GetScheduledTaskDetail 获取计划任务详情
// @Summary      获取计划任务详情
// @Description  根据任务ID获取计划任务的详细信息
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.ScheduledTaskDetailRequest  true  "任务ID"
// @Success      200  {object}  response.SuccessResponse{data=models.ScheduledTask}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/detail [post]
func GetScheduledTaskDetail(c *gin.Context) {
	var req models.ScheduledTaskDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 检查ID是否提供
	if req.ID == "" {
		response.BadRequest(c, "api.task.id_missing", nil)
		return
	}

	// 获取任务详情
	taskDetail, err := SchedulerService.GetScheduledTaskDetail(c, req.ID)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	if taskDetail == nil {
		response.NotFound(c, "api.task.not_found", nil)
		return
	}
	response.Success(c, taskDetail, "")
}

// UpdatePageMonitScheduledTask
// @Router       /api/task/scheduled/pagemonit/update [post]
func UpdatePageMonitScheduledTask(c *gin.Context) {
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	err := SchedulerService.UpdatePageMonitScheduledTask(c, &req)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// GetPageMonitData 获取页面监控数据列表
// @Summary      获取页面监控数据列表
// @Description  支持url字段模糊搜索与分页
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.PageMonitoringListRequest  true  "查询参数"
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]models.PageMonitoringTask,total=int64}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/pagemonit/data [post]
func GetPageMonitData(c *gin.Context) {
	var req models.PageMonitoringListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 设置默认分页参数
	if req.PageIndex <= 0 {
		req.PageIndex = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	tasks, total, err := SchedulerService.GetPageMonitData(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	resp := map[string]interface{}{
		"list":  tasks,
		"total": total,
	}
	response.Success(c, resp, "")
}

// AddPageMonitTask 添加页面监控任务
// @Summary      添加页面监控任务
// @Description  添加新的页面监控任务
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.PageMonitoringAddRequest  true  "URL信息"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/pagemonit/add [post]
func AddPageMonitTask(c *gin.Context) {
	var req models.PageMonitoringAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if req.URL == "" {
		response.BadRequest(c, "api.task.url_null", nil)
		return
	}

	if err := SchedulerService.AddPageMonitTask(c, req.URL); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}

// DeletePageMonitTask 删除页面监控任务
// @Summary      删除页面监控任务
// @Description  批量删除页面监控任务
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      models.PageMonitoringDeleteRequest  true  "任务ID列表"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/scheduled/pagemonit/delete [post]
func DeletePageMonitTask(c *gin.Context) {
	var req models.PageMonitoringDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids is required"))
		return
	}

	if err := SchedulerService.DeletePageMonitTasks(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, nil, "api.success")
}
