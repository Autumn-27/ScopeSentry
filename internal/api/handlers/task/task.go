// Package task -----------------------------
// @file      : task.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/4 22:04
// -------------------------------------------
package task

import (
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/task/common"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/task/task"
	"github.com/gin-gonic/gin"
)

var taskService task.Service
var commonService common.Service

// TaskDetailRequest 任务详情请求
type TaskDetailRequest struct {
	ID string `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
}

// List 获取任务列表（分页+模糊搜索）
// @Summary      获取任务列表
// @Description  支持name字段模糊搜索与分页
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      ListRequest  true  "查询参数"
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]models.Task,total=int}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task [post]
func List(c *gin.Context) {
	var req models.ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	go func() {
		err := taskService.TaskProgress(c)
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
		}
	}()
	tasks, total, err := taskService.List(c, req.Search, req.PageIndex, req.PageSize)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		response.InternalServerError(c, "api.task.list.error", err)
		return
	}
	resp := map[string]interface{}{
		"list":  tasks,
		"total": total,
	}
	response.Success(c, resp, "")
}

// TaskDetail 获取任务详情
// @Summary      获取任务详情
// @Description  根据任务ID获取任务的详细信息
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      TaskDetailRequest  true  "任务ID"
// @Success      200  {object}  response.SuccessResponse{data=object{name=string,target=string,ignore=string,node=[]string,allNode=bool,scheduledTasks=bool,hour=int,duplicates=string,template=string,day=int,minute=int,project=[]string,search=string,cycleType=string}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/detail [post]
func TaskDetail(c *gin.Context) {
	var req TaskDetailRequest
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
	taskDetail, err := taskService.GetTaskDetail(c, req.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("TaskDetail %v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	if taskDetail == nil {
		response.NotFound(c, "api.task.not_found", nil)
		return
	}
	response.Success(c, taskDetail, "")
}

// AddTask 添加任务
// @Summary 添加任务
// @Description 添加新的扫描任务
// @Tags Task
// @Accept json
// @Produce json
// @Param request body AddTaskRequest true "任务信息"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/task/add [post]
func AddTask(c *gin.Context) {
	var req models.Task
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}

	// 检查任务名是否已存在
	exists, err := taskService.CheckTaskNameExists(c, req.Name)
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

	// 创建任务
	_, err = commonService.Insert(c, &req)
	if err != nil {
		logger.Error(fmt.Sprintf("AddTask %v", err))
		response.InternalServerError(c, "api.error", err)
		return
	}

	// 如果是定时任务，创建定时任务
	if req.ScheduledTasks {
		err := SchedulerService.CreateScheduledTask(c, &req)
		if err != nil {
			response.InternalServerError(c, "api.error", err)
			return
		}
	}

	response.Success(c, nil, "api.task.add_success")
}

// DeleteTaskRequest 删除任务请求
// 支持批量删除与是否删除资产
// ids: 任务ID列表
// delA: 是否同时删除资产
type DeleteTaskRequest struct {
	IDs  []string `json:"ids" binding:"required" example:"[\"507f1f77bcf86cd799439011\",\"507f1f77bcf86cd799439012\"]"`
	DelA bool     `json:"delA" binding:"omitempty" example:"false"`
}

// RetestTaskRequest 重测任务请求
type RetestTaskRequest struct {
	ID string `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
}

// DeleteTask 删除任务
// @Summary      删除任务
// @Description  批量删除任务，支持可选删除其产生的资产，并清理调度与缓存
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      DeleteTaskRequest  true  "删除参数"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/delete [post]
func DeleteTask(c *gin.Context) {
	var req DeleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids is required"))
		return
	}

	if err := taskService.DeleteTasks(c, req.IDs, req.DelA); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// RetestTask 重测任务
// @Summary      重测任务
// @Description  重新创建扫描任务并重置进度
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      RetestTaskRequest  true  "任务ID"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/retest [post]
func RetestTask(c *gin.Context) {
	var req RetestTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if req.ID == "" {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("id is required"))
		return
	}

	if err := taskService.RetestTask(c, req.ID); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// StopTaskRequest 停止任务请求
type StopTaskRequest struct {
	IDs []string `json:"ids" binding:"required" example:"[\"507f1f77bcf86cd799439011\"]"`
}

// StartTaskRequest 开始任务请求
type StartTaskRequest struct {
	IDs []string `json:"ids" binding:"required" example:"[\"507f1f77bcf86cd799439011\"]"`
}

// StopTask 停止任务
// @Summary      停止任务
// @Description  批量停止任务，向节点广播停止并将状态置为2
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      StopTaskRequest  true  "任务ID列表"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/stop [post]
func StopTask(c *gin.Context) {
	var req StopTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids is required"))
		return
	}
	if err := taskService.StopTasks(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// StartTask 开始任务（从暂停恢复）
// @Summary      开始任务
// @Description  批量开始任务，从暂停恢复继续扫描
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      StartTaskRequest  true  "任务ID列表"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/start [post]
func StartTask(c *gin.Context) {
	var req StartTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if len(req.IDs) == 0 {
		response.BadRequest(c, "api.bad_request", fmt.Errorf("ids is required"))
		return
	}
	if err := taskService.StartTasks(c, req.IDs); err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}
	response.Success(c, nil, "api.success")
}

// ProgressInfoRequest 任务进度信息请求
type ProgressInfoRequest struct {
	ID        string `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
	PageIndex int    `json:"pageIndex" binding:"omitempty,min=1" example:"1"`
	PageSize  int    `json:"pageSize" binding:"omitempty,min=1,max=100" example:"10"`
}

// ProgressInfo 获取任务进度信息
// @Summary      获取任务进度信息
// @Description  获取指定任务的详细进度信息，包括各扫描阶段的时间戳
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      ProgressInfoRequest  true  "任务ID和分页参数"
// @Success      200  {object}  response.SuccessResponse{data=object{list=[]models.TaskProgress,total=int}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/progress/info [post]
func ProgressInfo(c *gin.Context) {
	var req ProgressInfoRequest
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

	progressData, err := taskService.GetTaskProgress(c, req.ID, req.PageIndex, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, progressData, "")
}

func init() {
	taskService = task.NewService()
	commonService = common.NewService()
}
