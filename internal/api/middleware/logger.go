package middleware

import (
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求路径
		path := c.Request.URL.Path
		// 请求参数
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		cost := time.Since(start)
		// 状态码
		status := c.Writer.Status()

		// 记录日志
		logger.Info("Request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		)
	}
}
