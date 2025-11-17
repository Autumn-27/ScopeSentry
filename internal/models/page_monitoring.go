package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// PageMonitoring 页面监控模型
type PageMonitoring struct {
	ID        string    `json:"id" bson:"_id"`
	URL       string    `json:"url" bson:"url"`
	Content   []string  `json:"content" bson:"content"`
	Hash      []string  `json:"hash" bson:"hash"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// PageMonitoringBody 页面监控内容模型
type PageMonitoringBody struct {
	ID      string   `json:"id" bson:"_id"`
	MD5     string   `json:"md5" bson:"md5"`
	Content []string `json:"diff" bson:"content"`
}

// PageMonitoringHistoryRequest 历史记录请求
type PageMonitoringHistoryRequest struct {
	ID string `json:"id"`
}

// PageMonitoringContentRequest 内容请求
type PageMonitoringContentRequest struct {
	ID   string `json:"id"`
	Flag string `json:"flag"` // 1: 获取上一次响应, 2: 获取当前响应
}

// PageMonitoringDiffRequest 差异对比请求
type PageMonitoringDiffRequest struct {
	ID string `json:"id"`
}

// PageMonitoringResultRequest 页面监控结果请求
type PageMonitoringResultRequest struct {
	PageIndex int                    `json:"pageIndex" binding:"required,min=1" example:"1"`
	PageSize  int                    `json:"pageSize" binding:"required,min=1,max=100" example:"10"`
	Query     map[string]interface{} `json:"query" binding:"required" example:"{\"url\":\"example.com\"}"`
}

// PageMonitoringResult 页面监控结果
type PageMonitoringResult struct {
	URL        string             `json:"url" bson:"url"`
	Hash       []string           `json:"hash" bson:"hash"`
	MD5        string             `json:"md5" bson:"md5"`
	Time       string             `json:"time" bson:"time"`
	StatusCode []int              `json:"statusCode" bson:"statusCode"`
	Similarity float64            `json:"similarity" bson:"similarity"`
	Tags       []string           `json:"tags" bson:"tags"`
	ID         primitive.ObjectID `bson:"_id" json:"id"`
}

// PageMonitoringResultResponse 页面监控结果响应
type PageMonitoringResultResponse struct {
	List  []PageMonitoringResult `json:"list"`
	Total int64                  `json:"total"`
}
