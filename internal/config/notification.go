// config-------------------------------------
// @file      : notification.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/12/22 22:08
// -------------------------------------------

package config

import (
	"encoding/json"
	"github.com/Autumn-27/ScopeSentry/internal/global"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/utils"
	"go.uber.org/zap"
	"strings"
)

var NotificationApi []models.NotificationApi

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}

func Notification(msg string) {
	safeMsg := jsonEscape(msg)
	for _, api := range NotificationApi {
		uri := strings.Replace(api.Url, "*msg*", safeMsg, -1)
		if api.Method == "GET" {
			_, err := utils.Requests.HttpGet(uri)
			if err != nil {
				global.Log.Error(msg, zap.Error(err))
			}
		} else {
			data := strings.Replace(api.Data, "*msg*", safeMsg, -1)
			err, _ := utils.Requests.HttpPost(uri, []byte(data), api.ContentType)
			if err != nil {
				global.Log.Error(msg, zap.Error(err))
			}
		}
	}
}
