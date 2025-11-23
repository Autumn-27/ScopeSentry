package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// 添加简化方法
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func init() {
	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.GlobalConfig.Log.Level)); err != nil {
		fmt.Println(err)
		return
	}

	// 加载时区
	location, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		// 如果加载失败，使用本地时区
		location = time.Local
		fmt.Println("load Timezone error")
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 彩色等级
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			// 使用时区格式化时间
			t = t.In(location)
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		}, // ISO时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 构建输出器
	var cores []zapcore.Core

	// 始终输出到控制台
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	cores = append(cores, consoleCore)

	// 如果配置了日志文件路径，则也输出到文件
	if config.GlobalConfig.Log.Filename != "" {
		logDir := filepath.Dir(config.GlobalConfig.Log.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Println(err)
			return
		}

		fileWriter := &lumberjack.Logger{
			Filename: config.GlobalConfig.Log.Filename,
			MaxSize:  300, // 单位 MB
			Compress: true,
		}

		fileCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(fileWriter), level)
		cores = append(cores, fileCore)
	}

	// 合并所有输出
	core := zapcore.NewTee(cores...)

	// 创建 Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return
}
