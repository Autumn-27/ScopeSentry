// common-------------------------------------
// @file      : common.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/9/25 23:40
// -------------------------------------------

package common

import (
	"github.com/Autumn-27/ScopeSentry/internal/repositories/common"
)

type Service struct {
	//taskRepo         task.Repository
	//commonService    commonTask.Service
	//schedulerService schedulerSvc.Service
	//nodeService      node.Service
	Repo common.Repository
}

func NewService() Service {
	return Service{
		Repo: common.NewRepository(),
	}
}
