// worker-------------------------------------
// @file      : start.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/11/2 17:57
// -------------------------------------------

package worker

import (
	"context"
	"time"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	// 退出时释放，避免 context 泄漏
	defer cancel()
	go IconHandle(ctx)
	go ScreenshotHandle(ctx)
	go IPAssetHandle(ctx)
	time.Sleep(3 * time.Second)
}
