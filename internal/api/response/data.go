// Package response -----------------------------
// @file      : data.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/3 12:02
// -------------------------------------------
package response

type DataResponse struct {
	List interface{} `json:"list"` // 使用 interface{} 作为通用字段类型
}
