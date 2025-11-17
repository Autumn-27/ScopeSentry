// Package logger -----------------------------
// @file      : gin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/24 13:55
// -------------------------------------------
package logger

//
//import (
//	"os"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"go.uber.org/zap"
//)
//
//// GinLogger 中间件，用 zap 替代默认日志
//func GinLogger() gin.HandlerFunc {
//	return gin.LoggerWithWriter(zapWriter{Logger})
//}
//
//// zapWriter 结构，用于将 gin 的日志写入 zap
//type zapWriter struct {
//	l *zap.Logger
//}
//
//func (z zapWriter) Write(p []byte) (n int, err error) {
//	z.l.Info(string(p)) // 你也可以在这里做日志结构化处理
//	return len(p), nil
//}
//
//// GinRecovery 替代默认 panic 恢复逻辑，记录错误日志
//func GinRecovery() gin.HandlerFunc {
//	return gin.CustomRecoveryWithWriter(os.Stderr, func(c *gin.Context, err interface{}) {
//		Logger.Error("Panic recovered",
//			zap.Any("error", err),
//			zap.String("path", c.Request.URL.Path),
//			zap.Time("time", time.Now()),
//		)
//		c.AbortWithStatusJSON(500, gin.H{
//			"message": "Server Error",
//		})
//	})
//}
