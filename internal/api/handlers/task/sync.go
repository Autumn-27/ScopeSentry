package task

import (
	"errors"

	"github.com/Autumn-27/ScopeSentry/internal/api/response"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/tasksync"
	"github.com/gin-gonic/gin"
)

var projectService project.Service

func init() {
	projectService = project.NewService()
}

// SyncToProject 把所选任务的目标同步到项目（已有项目或新建项目）。
//
// @Summary      同步任务目标到项目
// @Description  把任务目标并入已有项目，或按名称/标签新建项目，并重算项目资产
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        data  body      tasksync.Request  true  "同步参数"
// @Success      200  {object}  response.SuccessResponse{data=object{project=string,added=int,targets=int}}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /api/task/sync [post]
func SyncToProject(c *gin.Context) {
	var req tasksync.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "api.bad_request", err)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequestWithError(c, "api.bad_request", err)
		return
	}

	tasks, err := taskService.GetTasksByIDs(c, req.IDS)
	if err != nil {
		response.BadRequestWithError(c, "api.bad_request", err)
		return
	}
	targets, err := tasksync.ExtractTargets(tasks)
	if err != nil {
		response.BadRequestWithError(c, "api.bad_request", err)
		return
	}

	projectID, added, err := projectService.SyncTargets(c, req.Option, req.Project, req.Name, req.Tag, targets)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) || errors.Is(err, project.ErrProjectNameTaken) {
			response.BadRequestWithError(c, "api.bad_request", err)
			return
		}
		response.InternalServerError(c, "api.error", err)
		return
	}

	response.Success(c, gin.H{
		"project": projectID,
		"added":   added,
		"targets": len(targets),
	}, "api.success")
}
