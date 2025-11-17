// helper-------------------------------------
// @file      : time.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/25 15:04
// -------------------------------------------

package helper

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"time"
)

// GetNowTimeString 获取当前时间字符串，按指定时区
func GetNowTimeString() string {
	loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		// 若时区加载失败，使用 UTC 作为默认回退
		loc = time.UTC
	}
	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05")
}

func GetNowTime() time.Time {
	loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		// 若时区加载失败，使用 UTC 作为默认回退
		loc = time.UTC
	}
	return time.Now().In(loc)
}

// FormatTime 将 time.Time 格式化为字符串，使用全局配置的时区
func FormatTime(t time.Time) string {
	loc, err := time.LoadLocation(config.GlobalConfig.System.Timezone)
	if err != nil {
		loc = time.UTC
	}
	return t.In(loc).Format("2006-01-02 15:04:05")
}
