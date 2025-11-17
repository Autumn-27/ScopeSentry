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
	ctx, _ := context.WithCancel(context.Background())
	go IconHandle(ctx)
	go ScreenshotHandle(ctx)
	go IPAssetHandle(ctx)
	time.Sleep(3 * time.Second)
}
