package compress

import (
	"compress/gzip"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 静态资源路径（与 cmd/main 里的 StaticFS 对应）
var staticPrefixes = []string{
	"/assets", "/css", "/js", "/lib", "/plugins", "/png", "/draco", "/images",
}

// 已经是压缩格式或不需要压缩的扩展名
var skipCompressExt = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
	".ico": true, ".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
	".mp4": true, ".webm": true, ".zip": true, ".gz": true, ".br": true,
	".pdf": true, ".xlsx": true, ".xls": true, ".docx": true, ".doc": true,
}

func isStaticPath(path string) bool {
	for _, prefix := range staticPrefixes {
		if strings.HasPrefix(path, prefix+"/") || path == prefix {
			return true
		}
	}
	return false
}

// StaticCache 给内嵌前端资源加缓存头：
// 带 hash 的 /assets/* 可以长缓存，index.html 必须每次协商。
func StaticCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		switch {
		case strings.HasPrefix(path, "/assets/"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		case path == "/" || strings.HasSuffix(path, ".html"):
			c.Header("Cache-Control", "no-cache")
		case isStaticPath(path):
			c.Header("Cache-Control", "public, max-age=3600")
		}
		c.Next()
	}
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer      *gzip.Writer
	wroteHeader bool
	skip        bool
}

func (g *gzipResponseWriter) WriteHeader(code int) {
	if !g.wroteHeader {
		g.wroteHeader = true
		header := g.ResponseWriter.Header()
		// 压缩后长度会变，必须删掉原 Content-Length，并声明 Vary
		header.Del("Content-Length")
		header.Set("Content-Encoding", "gzip")
		header.Add("Vary", "Accept-Encoding")
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	if !g.wroteHeader {
		g.WriteHeader(http.StatusOK)
	}
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) finish() {
	if !g.wroteHeader {
		// 没有 body（如 204/304），不要留下 Content-Encoding
		g.ResponseWriter.Header().Del("Content-Encoding")
	}
	_ = g.writer.Close()
}

// Gzip 只压文本类响应（JS/CSS/HTML/JSON/SVG…），
// 跳过 HEAD、Range、图片与已压缩格式，避免破坏静态资源的断点续传。
func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if c.Request.Method == http.MethodHead ||
			c.Request.Header.Get("Range") != "" ||
			c.Request.Header.Get("Upgrade") != "" ||
			strings.HasPrefix(path, "/mcp") ||
			!strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") ||
			skipCompressExt[strings.ToLower(filepath.Ext(path))] {
			c.Next()
			return
		}

		gz := gzip.NewWriter(c.Writer)
		gw := &gzipResponseWriter{ResponseWriter: c.Writer, writer: gz}
		c.Writer = gw
		defer gw.finish()

		c.Next()
	}
}
