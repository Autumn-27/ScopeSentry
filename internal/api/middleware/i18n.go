package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// I18nMiddleware 国际化中间件
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取语言
		lang := c.GetHeader("Accept-Language")
		if lang == "" {
			lang = "zh" // 默认中文
		}

		// 解析语言标签
		tag, err := language.Parse(lang)
		if err != nil {
			tag = language.Chinese
		}

		// 设置语言到上下文
		c.Set("language", tag.String())

		// 设置响应头
		c.Header("Content-Language", tag.String())

		c.Next()
	}
}
