// models-------------------------------------
// @file      : configuration.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/29 21:33
// -------------------------------------------

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type System struct {
	ModulesConfig string `json:"modulesConfig" bson:"modulesConfig"`
}

type Notification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `json:"name" bson:"name"`               // 请求名称
	Method      string             `json:"method" bson:"method"`           // HTTP 方法，如 GET、POST
	URL         string             `json:"url" bson:"url"`                 // 请求地址
	ContentType string             `json:"contentType" bson:"contentType"` // 内容类型，如 application/json
	Data        string             `json:"data" bson:"data"`               // 请求数据（可为字符串或对象）
	State       bool               `json:"state" bson:"state"`             // 状态（如 active、disabled 等）
}

type UpdateNotification struct {
	ID          string `json:"id" bson:"-"`
	Name        string `json:"name" bson:"name"`               // 请求名称
	Method      string `json:"method" bson:"method"`           // HTTP 方法，如 GET、POST
	URL         string `json:"url" bson:"url"`                 // 请求地址
	ContentType string `json:"contentType" bson:"contentType"` // 内容类型，如 application/json
	Data        string `json:"data" bson:"data"`               // 请求数据（可为字符串或对象）
	State       bool   `json:"state" bson:"state"`             // 状态（如 active、disabled 等）
}

type DepConfig struct {
	Asset                bool   `json:"asset" bson:"asset"`
	Subdomain            bool   `json:"subdomain" bson:"subdomain"`
	SubdomainTakerResult bool   `json:"SubdomainTakerResult" bson:"SubdomainTakerResult"` // 保留原拼写
	UrlScan              bool   `json:"UrlScan" bson:"UrlScan"`
	Crawler              bool   `json:"crawler" bson:"crawler"`
	SensitiveResult      bool   `json:"SensitiveResult" bson:"SensitiveResult"`
	DirScanResult        bool   `json:"DirScanResult" bson:"DirScanResult"`
	Vulnerability        bool   `json:"vulnerability" bson:"vulnerability"`
	PageMonitoring       bool   `json:"PageMonitoring" bson:"PageMonitoring"`
	Hour                 int    `json:"hour" bson:"hour"`
	Flag                 bool   `json:"flag" bson:"flag"`
	RunNow               bool   `json:"runNow" bson:"runNow"`
	Name                 string `bson:"name" json:"name"`
}
