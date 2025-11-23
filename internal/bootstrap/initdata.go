// config-------------------------------------
// @file      : initdata.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/13 22:34
// -------------------------------------------
package bootstrap

import (
	"context"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProjectList 存储项目ID和名称的映射关系
var ProjectList = make(map[string]string)

// GetProjectList 从MongoDB获取项目列表并更新到ProjectList中
func GetProjectList() error {
	collection := mongodb.DB.Collection("project")
	opts := options.Find().SetProjection(bson.M{"_id": 1, "name": 1})
	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result struct {
			ID   string `bson:"_id"`
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err != nil {
			return err
		}
		ProjectList[result.ID] = result.Name
	}

	return cursor.Err()
}
