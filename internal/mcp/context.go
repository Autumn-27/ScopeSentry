package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ginContext(ctx context.Context) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	c.Request = req
	return c
}
