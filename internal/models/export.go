// models-------------------------------------
// @file      : export.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/25 14:41
// -------------------------------------------

package models

import "go.mongodb.org/mongo-driver/bson"

// ExportRequest 导出请求结构体
type ExportRequest struct {
	Index            string `json:"index" binding:"required" example:"asset"`
	SearchExpression string `json:"search"`
	// Filter 过滤条件
	Filter map[string][]interface{} `json:"filter"`
	// FuzzyQuery 模糊查询条件
	FuzzyQuery map[string]string `json:"-"`
	Quantity   int               `json:"quantity" binding:"required" example:"100"`
	Type       string            `json:"type" binding:"required" example:"search"`
	FileType   string            `json:"filetype" binding:"required" example:"csv"`
	Field      []string          `json:"field" binding:"required" example:"['host','ip']"`
}

// ExportResponse 导出响应结构体
type ExportResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ExportTask struct {
	FileName   string `json:"file_name" bson:"file_name"`
	CreateTime string `json:"create_time" bson:"create_time"`
	Quantity   int    `json:"quantity"`
	DataType   string `json:"data_type" bson:"data_type"`
	FileSize   string `json:"file_size" bson:"file_size"`
	FileType   string `json:"file_type" bson:"file_type"`
	State      int    `json:"state"`
	EndTime    string `json:"end_time" bson:"end_time"`
}

type ExportRun struct {
	Filter   bson.M
	FileName string `json:"file_name"`
	Index    string `json:"index"`
	FileType string `json:"file_type"`
	Quantity int    `json:"quantity"`
	Field    []string
}
