// jobs-------------------------------------
// @file      : scan].go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/6/10 20:33
// -------------------------------------------

package jobs

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/services/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	taskCommon "github.com/Autumn-27/ScopeSentry/internal/services/task/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// scanService 处理计划任务数据访问
type scanService struct {
	collection               *mongo.Collection
	pageMonitoringCollection *mongo.Collection
}

var taskCommonService taskCommon.Service
var ScanSvc *scanService
var NodeService node.Service
var CommonService common.Service

func init() {
	ScanSvc = &scanService{
		collection:               mongodb.DB.Collection("ScheduledTasks"),
		pageMonitoringCollection: mongodb.DB.Collection("PageMonitoring"),
	}
	taskCommonService = taskCommon.NewService()
	NodeService = node.NewService()
	CommonService = common.NewService()
}

func Scan(params string, nextTime string) error {
	// 获取计划任务数据
	ID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": ID}
	opts := options.FindOne()

	var task models.Task
	err = ScanSvc.collection.FindOne(context.Background(), filter, opts).Decode(&task)
	if err != nil {
		return fmt.Errorf("failed to get scheduled task: %w", err)
	}
	task.Name += "-" + helper.GetNowTimeString()
	_, err = taskCommonService.Insert(GetTestContext(), &task)
	if err != nil {
		logger.Error(fmt.Sprintf("AddTask Scan%v", err))
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"lastTime": helper.GetNowTimeString(),
			"nextTime": nextTime,
		},
	}

	_, err = ScanSvc.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Update ScheduledTasks error: %v", err))
		return err
	}

	return nil
}

func GetTestContext() *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}
