// Package system -----------------------------
// @file      : system.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/2 23:01
// -------------------------------------------
package system

import (
	"encoding/json"
	"github.com/Autumn-27/ScopeSentry-go/internal/constants"
	"net/http"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/api/response"
	"github.com/Autumn-27/ScopeSentry-go/internal/database/redis"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

type VersionResponse struct {
	Code int `json:"code"`
	Data struct {
		List []VersionInfo `json:"list"`
	} `json:"data"`
}

type VersionInfo struct {
	Name     string `json:"name"`
	CVersion string `json:"cversion"`
	LVersion string `json:"lversion"`
	Msg      string `json:"msg"`
}

type RemoteVersion struct {
	Server    string `json:"server"`
	ServerMsg string `json:"server_msg"`
	Scan      string `json:"scan"`
	ScanMsg   string `json:"scan_msg"`
}

// GetSystemVersion @Summary      获取系统及节点版本信息
// @Description  获取 ScopeSentry 服务端及所有节点的当前版本信息
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.SuccessResponse{data=VersionListResponse}
// @Failure      400  {object}  response.BadRequestResponse
// @Failure      401  {object}  response.UnauthorizedResponse
// @Failure      404  {object}  response.NotFoundResponse
// @Failure      500  {object}  response.InternalServerErrorResponse
// @Router       /system/version [get]
func GetSystemVersion(c *gin.Context) {
	var (
		serverLVersion string
		serverMsg      string
		scanLVersion   string
		scanMsg        string
	)

	client := &fasthttp.Client{
		ReadTimeout: 10 * time.Second,
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("https://gitee.com/constL/scope-sentry/raw/main/version.json")
	req.Header.SetMethod("GET")

	err := client.Do(req, resp)
	if err == nil && resp.StatusCode() == http.StatusOK {
		var version RemoteVersion
		if err := json.Unmarshal(resp.Body(), &version); err == nil {
			serverLVersion = version.Server
			serverMsg = version.ServerMsg
			scanLVersion = version.Scan
			scanMsg = version.ScanMsg
		}
	}

	// 如果gitee失败，尝试从github获取
	if serverLVersion == "" {
		req.SetRequestURI("https://raw.githubusercontent.com/Autumn-27/ScopeSentry/main/version.json")
		err = client.Do(req, resp)
		if err == nil && resp.StatusCode() == http.StatusOK {
			var version RemoteVersion
			if err := json.Unmarshal(resp.Body(), &version); err == nil {
				serverLVersion = version.Server
				serverMsg = version.ServerMsg
				scanLVersion = version.Scan
				scanMsg = version.ScanMsg
			}
		}
	}

	resultList := []VersionInfo{
		{
			Name:     "ScopeSentry-Server",
			CVersion: constants.Version,
			LVersion: serverLVersion,
			Msg:      serverMsg,
		},
	}

	ctx := c.Request.Context()
	keys, err := redis.Client.Keys(ctx, "node:*").Result()
	if err != nil {
		response.InternalServerError(c, "api.system.version.redis_error", err)
		return
	}
	for _, key := range keys {
		name := key[5:]
		hashData, err := redis.Client.HGetAll(ctx, key).Result()
		if err != nil {
			response.InternalServerError(c, "api.system.version.redis_hash_error", err)
			return
		}
		resultList = append(resultList, VersionInfo{
			Name:     name,
			CVersion: hashData["version"],
			LVersion: scanLVersion,
			Msg:      scanMsg,
		})
	}

	response.Success(c, VersionListResponse{
		List: resultList,
	}, "api.system.version.success")
}

// VersionListResponse 版本列表响应结构
// 用于统一响应格式
// @Description 版本信息列表
// @Author Autumn
// @Date 2025-05-02
// @Contact rainy-autumn@outlook.com
type VersionListResponse struct {
	List []VersionInfo `json:"list"`
}
