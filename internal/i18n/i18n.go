package i18n

import (
	"embed"
	"encoding/json"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var (
	bundle *i18n.Bundle
	once   sync.Once
)

func init() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.Chinese)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// 加载语言文件
		loadLocaleFiles()
	})
}

func loadLocaleFiles() {
	files, err := localeFS.ReadDir("locales")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := localeFS.ReadFile("locales/" + file.Name())
		if err != nil {
			panic(err)
		}
		bundle.MustParseMessageFileBytes(data, file.Name())
	}
}

func Translate(locale, messageID string) string {
	localizer := i18n.NewLocalizer(bundle, locale)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID // 如果翻译失败，返回原始key
	}
	return msg
}
