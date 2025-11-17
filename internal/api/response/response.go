package response

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"net/http"
	"runtime"

	"github.com/Autumn-27/ScopeSentry-go/internal/i18n"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // HTTP 状态码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// SuccessResponse 成功响应示例
type SuccessResponse struct {
	Code    int         `json:"code" example:"200"`     // 成功状态码
	Message string      `json:"message" example:"操作成功"` // 成功消息
	Data    interface{} `json:"data,omitempty"`         // 响应数据
}

// BadRequestResponse 请求错误响应示例
type BadRequestResponse struct {
	Code    int    `json:"code" example:"400"`     // 请求错误状态码
	Message string `json:"message" example:"请求错误"` // 错误消息
}

// UnauthorizedResponse 未授权响应示例
type UnauthorizedResponse struct {
	Code    int    `json:"code" example:"401"`      // 未授权状态码
	Message string `json:"message" example:"未授权访问"` // 错误消息
}

// NotFoundResponse 未找到响应示例
type NotFoundResponse struct {
	Code    int    `json:"code" example:"404"`      // 未找到状态码
	Message string `json:"message" example:"资源未找到"` // 错误消息
}

// InternalServerErrorResponse 服务器错误响应示例
type InternalServerErrorResponse struct {
	Code    int    `json:"code" example:"500"`        // 服务器错误状态码
	Message string `json:"message" example:"服务器内部错误"` // 错误消息
}

type PageResponse struct {
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// getLocale 获取请求的语言
func getLocale(c *gin.Context) string {
	locale := c.GetHeader("Accept-Language")
	if locale == "" {
		locale = "zh-CN" // 默认中文
	}
	return locale
}

func Success(c *gin.Context, data interface{}, msgKey string) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: i18n.Translate(getLocale(c), msgKey),
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}, msgKey string) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: i18n.Translate(getLocale(c), msgKey),
		Data:    data,
	})
}

func BadRequest(c *gin.Context, msgKey string, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: i18n.Translate(getLocale(c), msgKey),
	})
}

func NotFound(c *gin.Context, msgKey string, err error) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: i18n.Translate(getLocale(c), msgKey),
	})
}

func InternalServerError(c *gin.Context, msgKey string, err error) {
	pc, file, line, ok := runtime.Caller(1)
	callerInfo := ""
	if ok {
		fn := runtime.FuncForPC(pc)
		callerInfo = fmt.Sprintf("%s:%d (%s)", file, line, fn.Name())
	}

	logger.Error(fmt.Sprintf("InternalServerError at %s: %v", callerInfo, err))
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: i18n.Translate(getLocale(c), msgKey),
		Data:    fmt.Sprintf("%v", err),
	})
}

func Unauthorized(c *gin.Context, msgKey string, err error) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: i18n.Translate(getLocale(c), msgKey),
		Data:    fmt.Sprintf("%v", err),
	})
}
