package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func doRequest(t *testing.T, mongoPing, redisPing PingFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", Handler(mongoPing, redisPing))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	return w
}

func TestHealthOK(t *testing.T) {
	ok := func(context.Context) error { return nil }
	w := doRequest(t, ok, ok)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHealthDegradedWhenMongoFails(t *testing.T) {
	w := doRequest(t, func(context.Context) error { return errors.New("mongo down") }, func(context.Context) error { return nil })
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
	if !contains(w.Body.String(), "mongo down") {
		t.Fatalf("expected error detail in body: %s", w.Body.String())
	}
}

func TestHealthDegradedWhenRedisFails(t *testing.T) {
	w := doRequest(t, func(context.Context) error { return nil }, func(context.Context) error { return errors.New("redis down") })
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestHealthNilPingIsFailure(t *testing.T) {
	w := doRequest(t, nil, func(context.Context) error { return nil })
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("nil ping must be reported as failure, got %d", w.Code)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
