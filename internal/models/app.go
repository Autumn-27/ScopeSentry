package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// App 表示一个应用程序
type App struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`                 // 应用名称
	PackageName string             `bson:"package_name" json:"package_name"` // 包名
	Version     string             `bson:"version" json:"version"`           // 版本号
	Platform    string             `bson:"platform" json:"platform"`         // 平台（android/ios）
	Size        int64              `bson:"size" json:"size"`                 // 应用大小（字节）
	Icon        string             `bson:"icon" json:"icon"`                 // 应用图标URL
	Description string             `bson:"description" json:"description"`   // 应用描述
	Developer   string             `bson:"developer" json:"developer"`       // 开发者
	Category    string             `bson:"category" json:"category"`         // 应用分类
	Tags        []string           `bson:"tags" json:"tags"`                 // 标签
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`     // 更新时间
}

// AppResult 应用结果
type AppResult struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Icp         string             `bson:"icp" json:"icp"`
	Company     string             `bson:"company" json:"company"`
	Project     string             `bson:"project" json:"project"`
	Time        string             `bson:"time" json:"time"`
	Tags        []string           `bson:"tags" json:"tags"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	BundleID    string             `bson:"bundleID" json:"bundleID"`
	Apk         string             `bson:"apk" json:"apk"`
	URL         string             `bson:"url" json:"url"`
}

// AppResponse 应用响应
type AppResponse struct {
	List  []AppResult `json:"list"`
	Total int64       `json:"total"`
}
