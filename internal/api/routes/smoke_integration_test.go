//go:build integration

package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/database/redis"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"go.mongodb.org/mongo-driver/bson"
)

// API 冒烟：走真实路由 + 中间件 + Mongo/Redis，验证「登录 -> 带 token 访问受保护接口」
// 这条链路是通的，以及未认证请求确实被拦下。等价于前端 E2E，但不依赖浏览器。
func TestAPISmokeLoginAndProtectedRoute(t *testing.T) {
	if mongodb.DB == nil || redis.Client == nil {
		t.Skip("mongo/redis not initialized")
	}

	const (
		username = "e2e-smoke"
		password = "e2e-smoke-pass"
	)
	users := mongodb.DB.Collection("user")
	_, _ = users.DeleteMany(t.Context(), bson.M{"username": username})
	t.Cleanup(func() { _, _ = users.DeleteMany(t.Context(), bson.M{"username": username}) })
	if _, err := users.InsertOne(t.Context(), bson.M{
		"username": username,
		"password": helper.Sha256Hex(password),
		"email":    "e2e@example.invalid",
		"role":     "super_admin",
		"status":   "active",
	}); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	router := SetupRouter()

	// 1) 健康检查
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("healthz = %d: %s", w.Code, w.Body.String())
	}

	// 2) 未认证访问受保护接口 -> 401
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/assets/statistics", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated statistics = %d, want 401", w.Code)
	}

	// 3) 登录拿 token
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	w = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", w.Code, w.Body.String())
	}
	var loginResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login: %v (%s)", err, w.Body.String())
	}
	data, _ := loginResp["data"].(map[string]interface{})
	token, _ := data["access_token"].(string)
	if token == "" {
		t.Fatalf("login payload missing access_token: %s", w.Body.String())
	}

	// 4) 带 token 访问受保护接口 -> 200
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/assets/statistics", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("authenticated statistics = %d: %s", w.Code, w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("assetCount")) {
		t.Fatalf("statistics payload missing assetCount: %s", w.Body.String())
	}
}
