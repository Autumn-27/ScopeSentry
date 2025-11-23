// update-------------------------------------
// @file      : update181.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/11/23 17:17
// -------------------------------------------

package update

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update181() {
	// 资产bodyhash索引
	asset := mongodb.DB.Collection("asset")
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"host", 1}}},
		{Keys: bson.D{{"rootDomain", 1}, {"time", -1}}}, // 优化按 rootDomain 查询并按 time 排序的性能
	}
	_, err := asset.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create asset indexes: %v", err))
	}
}
