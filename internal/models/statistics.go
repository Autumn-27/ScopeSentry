// models-------------------------------------
// @file      : STATISTICS.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/9 21:16
// -------------------------------------------

package models

// StatisticsData 资产统计数据
type StatisticsData struct {
	AssetCount         int64 `json:"assetCount"`
	SubdomainCount     int64 `json:"subdomainCount"`
	SensitiveCount     int64 `json:"sensitiveCount"`
	UrlCount           int64 `json:"urlCount"`
	VulnerabilityCount int64 `json:"vulnerabilityCount"`
}

// PortStatisticsRequest 端口统计请求
type PortStatisticsRequest struct {
	Filter map[string]interface{} `json:"filter" binding:"required" example:"{\"type\":\"http\"}"`
}

// PortStatisticsResponse 端口统计响应
type PortStatisticsResponse struct {
	Port []PortItem `json:"Port"`
}

type PortItem struct {
	Value  string `bson:"_id" json:"value" example:"80"`
	Number int64  `bson:"number" json:"number" example:"100"`
}

// TitleStatisticsRequest 标题统计请求
type TitleStatisticsRequest struct {
	Filter map[string]interface{} `json:"filter" binding:"required" example:"{\"type\":\"http\"}"`
}

// TitleStatisticsResponse 标题统计响应
type TitleStatisticsResponse struct {
	Title []TitleItem `json:"Title"`
}

type TitleItem struct {
	Value  string `json:"value" example:"Example Title"`
	Number int64  `json:"number" example:"50"`
}

// TypeStatisticsRequest 服务类型统计请求
type TypeStatisticsRequest struct {
	Filter map[string]interface{} `json:"filter" binding:"required" example:"{\"type\":\"http\"}"`
}

// TypeStatisticsResponse 服务类型统计响应
type TypeStatisticsResponse struct {
	Service []TypeItem `json:"Service"`
}

type TypeItem struct {
	Value  string `bson:"_id" json:"value" example:"nginx"`
	Number int64  `bson:"number" json:"number" example:"200"`
}

// IconStatisticsRequest 图标统计请求
type IconStatisticsRequest struct {
	Filter   map[string]interface{} `json:"filter" binding:"required" example:"{\"type\":\"http\"}"`
	Page     int                    `json:"page" binding:"required" example:"1"`
	PageSize int                    `json:"page_size" binding:"required" example:"10"`
}

// IconStatisticsResponse 图标统计响应
type IconStatisticsResponse struct {
	Icon []IconItem `json:"Icon"`
}

type IconItem struct {
	Value    string `bson:"iconcontent" json:"value" example:"base64_icon_content"`
	Number   int64  `bson:"number" json:"number" example:"30"`
	IconHash string `bson:"_id" json:"icon_hash" example:"abc123"`
}

// AppStatisticsRequest 应用统计请求
type AppStatisticsRequest struct {
	Filter map[string]interface{} `json:"filter" binding:"required" example:"{\"type\":\"http\"}"`
}

// AppStatisticsResponse 应用统计响应
type AppStatisticsResponse struct {
	Product []AppItem `json:"Product"`
}

type AppItem struct {
	Value  string `bson:"_id" json:"value" example:"WordPress"`
	Number int64  `bson:"number" json:"number" example:"150"`
}

// APP 应用映射表
var APP = map[string]string{
	// 这里需要添加应用映射关系
}
