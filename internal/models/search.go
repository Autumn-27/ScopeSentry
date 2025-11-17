// models-------------------------------------
// @file      : search.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/7 19:27
// -------------------------------------------

package models

// SearchRequest 定义了搜索请求的参数结构
type SearchRequest struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	// Index 搜索类型，对应SearchKey中的key
	Index            string `json:"index"`
	SearchExpression string `json:"search"`
	// Filter 过滤条件
	Filter map[string][]interface{} `json:"filter"`
	// FuzzyQuery 模糊查询条件
	FuzzyQuery map[string]string `json:"-"`
	Sort       map[string]string `json:"sort"`
	Sid        string            `json:"sid"`
}
