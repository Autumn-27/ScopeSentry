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
	"github.com/Autumn-27/ScopeSentry/internal/constants"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update19() {
	// 资产索引
	asset := mongodb.DB.Collection("asset")
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"host", 1}}},
		{Keys: bson.D{{"rootDomain", 1}, {"time", -1}}}, // 优化按 rootDomain 查询并按 time 排序的性能
	}
	_, err := asset.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create asset indexes: %v", err))
	}
	// 旧指纹备份
	ctx := context.Background()
	fingerprintCollection := mongodb.DB.Collection("FingerprintRules")
	backupPipeline := []bson.M{
		{"$out": "FingerprintRules_bak"},
	}
	cursor, err := fingerprintCollection.Aggregate(ctx, backupPipeline)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to backup FingerprintRules: %v", err))
		return
	}
	defer cursor.Close(ctx)
	logger.Info("FingerprintRules backed up to FingerprintRules_bak successfully")

	// 插入新指纹
	fingerprintCollection = mongodb.DB.Collection("FingerprintRules")
	fingerprint, _ := constants.GetFingerprintData()
	if len(fingerprint) > 0 {
		_, err = fingerprintCollection.InsertMany(context.Background(), fingerprint)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create FingerprintRules: %v", err))
		}
	}

	// 插入指纹版本
	configCollection := mongodb.DB.Collection("config")
	configCollection.InsertOne(context.Background(), bson.M{
		"name":  "FingerVersion",
		"value": "2026-01-20 23:50",
		"type":  "finger",
	})
}
