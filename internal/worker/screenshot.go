// worker-------------------------------------
// @file      : screenshot.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/11/2 17:58
// -------------------------------------------

package worker

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils"
	"path"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ScreenshotHandle(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("任务发生错误，已恢复:", r)
			// 出现 panic 后，如果未取消，则重启任务循环
			select {
			case <-ctx.Done():
				fmt.Println("任务已取消，ScreenshotHandle不再重启")
			default:
				go ScreenshotHandle(ctx)
			}
		}
	}()

	// 初始化 common repository
	repo := common.NewRepository()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务被取消，ScreenshotHandle停止执行")
			return
		default:

			// 批量查询 icon 集合数据，每次查询 100 条
			limit := int64(100)
			opts := options.Find().SetLimit(limit)

			// 查询数据
			results, err := repo.FindMany(ctx, "screenshot", bson.M{}, opts)
			if err != nil {
				logger.Error(fmt.Sprintf("查询 screenshot 数据失败: %v\n", err))
				time.Sleep(10 * time.Second) // 出错后等待 5 秒再重试
				continue
			}

			// 如果没有数据，说明处理完毕
			if len(results) == 0 {
				time.Sleep(60 * time.Second) // 等待一段时间后再次检查
				continue
			}

			// 收集要删除的 _id
			var ids []primitive.ObjectID

			// 打印并收集数据
			for _, item := range results {
				// 提取 fav3 和 content 字段
				hash, _ := item["hash"].(string)
				content, _ := item["content"].(string)
				// 收集 _id 用于批量删除
				if id, ok := item["_id"].(primitive.ObjectID); ok {
					ids = append(ids, id)
				}
				filepath := path.Join(config.GlobalConfig.System.ScreenshotDir, fmt.Sprintf("%v.png", hash))
				decoded, err := base64.StdEncoding.DecodeString(content)
				if err != nil {
					logger.Error(fmt.Sprintf("<UNK> screenshot base64 decode <UNK>: %v %v\n", err, content))
					continue
				}
				err = utils.WriteFile(filepath, decoded)
				if err != nil {
					logger.Error(fmt.Sprintf("<UNK> screenshot WriteFile<UNK>: %v\n", err))
					continue
				}

			}

			// 批量删除数据
			if len(ids) > 0 {
				filter := bson.M{"_id": bson.M{"$in": ids}}
				_, err := repo.DeleteMany(ctx, "screenshot", filter)
				if err != nil {
					fmt.Printf("删除 icon 数据失败: %v\n", err)
				}
			}

			// 处理完一批后短暂休眠，避免过于频繁的数据库操作
			time.Sleep(3 * time.Second)
		}
	}
}
