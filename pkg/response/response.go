// Package response -----------------------------
// @file      : response.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/4 22:04
// -------------------------------------------
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// Ok 成功响应
func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
	})
}

// OkWithMessage 带消息的成功响应
func OkWithMessage(message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
	})
}

// OkWithData 带数据的成功响应
func OkWithData(data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// OkWithMessageAndData 带消息和数据的成功响应
func OkWithMessageAndData(message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    500,
		Message: "fail",
	})
}

// FailWithMessage 带消息的失败响应
func FailWithMessage(message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    500,
		Message: message,
	})
}

// FailWithCodeAndMessage 带状态码和消息的失败响应
func FailWithCodeAndMessage(code int, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// FailWithCodeAndMessageAndData 带状态码、消息和数据的失败响应
func FailWithCodeAndMessageAndData(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
