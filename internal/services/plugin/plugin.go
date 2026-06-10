package plugin

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/utils"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/plugins"
	schedulerCore "github.com/Autumn-27/ScopeSentry/internal/scheduler"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/options"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Service 定义插件服务接口
type Service interface {
	List(ctx *gin.Context, req *models.PluginListRequest) (*models.PluginListResponse, error)
	ListByModule(ctx *gin.Context, module string) ([]models.Plugin, error)
	GetPluginByHash(ctx context.Context, hash string) (*models.Plugin, error)
	Detail(ctx *gin.Context, req *models.PluginDetailRequest) (*models.Plugin, error)
	Save(ctx *gin.Context, req *models.PluginSaveRequest) error
	Delete(ctx *gin.Context, req *models.PluginDeleteRequest) error
	GetLogs(ctx *gin.Context, req *models.PluginLogRequest) (string, error)
	CleanLogs(ctx *gin.Context, req *models.PluginLogRequest) error
	CleanAllLogs(ctx *gin.Context) error
	Import(ctx *gin.Context, filePath string, key string) error
	Reinstall(ctx *gin.Context, req *models.PluginReinstallRequest) error
	Recheck(ctx *gin.Context, req *models.PluginRecheckRequest) error
	Uninstall(ctx *gin.Context, req *models.PluginUninstallRequest) error
	CheckKey(ctx *gin.Context, req *models.PluginKeyCheckRequest) error
	SearchRemotePlugins(ctx *gin.Context) (map[string]interface{}, error)
	ImportByData(ctx *gin.Context, req *models.PluginImportByDataRequest) error
	UpdateStatus(ctx *gin.Context, req *models.PluginStatusRequest) error
	Run(ctx *gin.Context, req *models.PluginRunRequest) error
	RunTaskEnd(task models.Task) error
}

// service 实现插件服务接口
type service struct {
	repo        plugin.Repository
	nodeService node.Service
}

// NewService 创建新的插件服务实例
func NewService() Service {
	return &service{
		repo:        plugin.NewRepository(),
		nodeService: node.NewService(),
	}
}

// List 获取插件列表
func (s *service) List(ctx *gin.Context, req *models.PluginListRequest) (*models.PluginListResponse, error) {
	query := bson.M{}
	if req.Type != "all" {
		if req.Type != "server" {
			query["type"] = bson.M{"$ne": "server"}
		} else {
			query["type"] = req.Type
		}
	}
	if req.Search != "" {
		query["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	opts := mongoOptions.Find().
		SetSkip(int64((req.PageIndex - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	plugins, err := s.repo.FindWithPagination(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find plugins: %w", err)
	}

	total, err := s.repo.Count(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count plugins: %w", err)
	}

	return &models.PluginListResponse{
		List:  plugins,
		Total: total,
	}, nil
}

// ListByModule 根据模块获取插件列表
func (s *service) ListByModule(ctx *gin.Context, module string) ([]models.Plugin, error) {
	result, err := s.repo.FindByModule(ctx, module)
	if result == nil {
		return []models.Plugin{}, err
	}
	return result, err
}

// GetPluginByHash 根据 hash 获取单个插件
func (s *service) GetPluginByHash(ctx context.Context, hash string) (*models.Plugin, error) {
	return s.repo.FindByHash(ctx, hash)
}

// Detail 获取插件详情
func (s *service) Detail(ctx *gin.Context, req *models.PluginDetailRequest) (*models.Plugin, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	return s.repo.FindByID(ctx, id)
}

// Save 保存插件
func (s *service) Save(ctx *gin.Context, req *models.PluginSaveRequest) error {
	// 验证密钥
	key, err := ioutil.ReadFile("PLUGINKEY")
	if err != nil {
		return fmt.Errorf("failed to read plugin key: %w", err)
	}
	if req.Key != strings.TrimSpace(string(key)) {
		return fmt.Errorf("invalid plugin key")
	}

	if req.ID == "" {
		// 创建新插件
		plugin := &models.Plugin{
			Name:          req.Name,
			Module:        req.Module,
			Hash:          generatePluginHash(32, req.Type),
			Parameter:     req.Parameter,
			Help:          req.Help,
			Introduction:  req.Introduction,
			Source:        req.Source,
			Version:       req.Version,
			ParameterList: req.ParameterList,
			Status:        true,
			Type:          req.Type,
			IsSystem:      false,
		}
		id, err := s.repo.Create(ctx, plugin)
		if err != nil {
			return fmt.Errorf("failed to create plugin: %w", err)
		}
		if req.Type != "server" {
			go func() {
				msg := models.Message{
					Name:    "all",
					Type:    "install_plugin",
					Content: fmt.Sprintf(`%v`, id),
				}
				err = s.nodeService.RefreshConfig(ctx, msg)
				if err != nil {
					logger.Error("failed to refresh config", zap.Error(err))
				}
			}()
		} else {
			RegisterPlugin(plugin.Source, plugin.Hash)
		}
	} else {
		// 更新现有插件
		id, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			return fmt.Errorf("invalid id format: %w", err)
		}

		update := bson.M{
			"name":          req.Name,
			"module":        req.Module,
			"parameter":     req.Parameter,
			"help":          req.Help,
			"introduction":  req.Introduction,
			"parameterList": req.ParameterList,
			"source":        req.Source,
			"version":       req.Version,
		}
		err = s.repo.Update(ctx, id, update)
		if err != nil {
			return err
		}
		if req.Type != "server" {
			go func() {
				msg := models.Message{
					Name:    "all",
					Type:    "install_plugin",
					Content: fmt.Sprintf(`%v`, req.ID),
				}
				err = s.nodeService.RefreshConfig(ctx, msg)
				if err != nil {
					logger.Error("failed to refresh config", zap.Error(err))
				}
			}()
		} else {
			// 如果是更新 也重新注册插件
			RegisterPlugin(req.Source, req.Hash)
		}
	}
	// 服务端插件 无论是插入还是更新，判断cycle是不是1 如果是1 表示运行1次
	if req.Type == "server" {
		plugins.RunPluginOnce(req.Hash)
	}
	return nil
}

// Delete 删除插件
func (s *service) Delete(ctx *gin.Context, req *models.PluginDeleteRequest) error {
	var hashes []string
	for _, item := range req.Data {
		hashes = append(hashes, item.Hash)
		go func() {
			if item.Module != "" {
				msg := models.Message{
					Name:    "all",
					Type:    "delete_plugin",
					Content: fmt.Sprintf(`%v_%v`, item.Hash, item.Module),
				}
				err := s.nodeService.RefreshConfig(ctx, msg)
				if err != nil {
					logger.Error("failed to refresh config", zap.Error(err))
				}
			} else {
				// module 为空 说明是服务端插件
				// 删除插件
				plugins.GlobalPluginManager.DeletePlugin(item.Hash)
				// 移除计划任务
				err := schedulerCore.GetGlobalScheduler().RemoveJob(item.Hash)
				if err != nil {
					logger.Error(fmt.Sprintf("failed to remove job: %v", err))
				}
			}
		}()
	}
	return s.repo.DeleteByHash(ctx, hashes)
}

// GetLogs 获取插件日志
func (s *service) GetLogs(ctx *gin.Context, req *models.PluginLogRequest) (string, error) {
	logKey := ""
	if req.Type == "server" {
		logKey = fmt.Sprintf("logs:server_plugins:%v", req.Hash)
	} else {
		logKey = fmt.Sprintf("logs:plugins:%v:%v", req.Module, req.Hash)
	}
	logs, err := s.repo.GetLogs(ctx, logKey)
	if err != nil {
		return "", err
	}
	return logs, nil
}

// CleanLogs 清理插件日志
func (s *service) CleanLogs(ctx *gin.Context, req *models.PluginLogRequest) error {
	if req.Hash == "" {
		return fmt.Errorf("module and hash are required")
	}
	var logKey string
	if req.Type == "server" {
		logKey = fmt.Sprintf("logs:server_plugins:%v", req.Hash)
	} else {
		logKey = fmt.Sprintf("logs:plugins:%v:%v", req.Module, req.Hash)
	}
	err := s.repo.CleanLogs(ctx, logKey)
	if err != nil {
		return err
	}
	return nil
}

// CleanAllLogs 清除所有插件日志
func (s *service) CleanAllLogs(ctx *gin.Context) error {
	err := s.repo.CleanAllLogs(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Import 导入插件
func (s *service) Import(ctx *gin.Context, filePath string, reqKey string) error {
	key, err := ioutil.ReadFile("PLUGINKEY")
	if err != nil {
		return fmt.Errorf("failed to read plugin key: %w", err)
	}
	if reqKey != strings.TrimSpace(string(key)) {
		return fmt.Errorf("invalid plugin key")
	}
	// 2. 读取 zip 文件字节内容
	zipData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取 zip 文件失败: %w", err)
	}

	// 3. 初始化 zip.Reader
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("打开 zip 文件失败: %w", err)
	}

	var pluginInfo models.PluginInfo
	var pluginSource string

	// 4. 遍历压缩包中的文件
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		baseName := strings.ToLower(filepath.Base(file.Name))

		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("打开压缩包中文件失败: %w", err)
		}
		content, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return fmt.Errorf("读取压缩包中文件内容失败: %w", err)
		}

		switch baseName {
		case "info.json":
			if err := json.Unmarshal(content, &pluginInfo); err != nil {
				return fmt.Errorf("解析 info.json 失败: %w", err)
			}
		case "plugin.go":
			pluginSource = string(content)
		}
	}

	// 5. 校验 info.json 的关键信息
	if pluginInfo.Name == "" {
		return fmt.Errorf("info.json 缺少 name字段")
	}
	if pluginInfo.Module == "" && pluginInfo.Type != "server" {
		return fmt.Errorf("info.json 缺少 module字段")
	}

	if pluginInfo.Hash == "" {
		pluginInfo.Hash = generatePluginHash(32, pluginInfo.Type)
	}
	var PLUGINS = make(map[string]bool)
	for _, p := range constants.Plugins {
		PLUGINS[p.Hash] = true
	}
	// 6. 插件是否是系统插件 或 module 不合法
	if PLUGINS[pluginInfo.Hash] {
		return fmt.Errorf("插件已存在")
	}
	var pluginsModules = make(map[string]bool)
	for _, module := range constants.PLUGINSMODULES {
		pluginsModules[module] = true
	}
	if pluginInfo.Type != "server" {
		if !pluginsModules[pluginInfo.Module] {
			return fmt.Errorf("模块非法: %s", pluginInfo.Module)
		}
	}
	// 7. 设置其他字段
	pluginInfo.Source = pluginSource
	pluginInfo.IsSystem = false

	// 8. 插入数据库

	// 9. 触发配置刷新

	p := &models.Plugin{
		Name:          pluginInfo.Name,
		Module:        pluginInfo.Module,
		Hash:          pluginInfo.Hash,
		Parameter:     pluginInfo.Parameter,
		Help:          pluginInfo.Help,
		Introduction:  pluginInfo.Introduction,
		ParameterList: pluginInfo.ParameterList,
		Source:        pluginInfo.Source,
		Version:       pluginInfo.Version,
		Type:          pluginInfo.Type,
		IsSystem:      false,
	}
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to create plugin: %w", err)
	}
	go func() {
		if pluginInfo.Type != "server" {
			msg := models.Message{
				Name:    "all",
				Type:    "install_plugin",
				Content: fmt.Sprintf(`%v`, id),
			}
			err = s.nodeService.RefreshConfig(ctx, msg)
			if err != nil {
				logger.Error("failed to refresh config", zap.Error(err))
			}
		}
	}()
	if pluginInfo.Type == "server" {
		RegisterPlugin(pluginInfo.Source, pluginInfo.Hash)
		// 判断是否需要运行一次
		plugins.RunPluginOnce(pluginInfo.Hash)
	}

	return nil
}

// Reinstall 重装插件
func (s *service) Reinstall(ctx *gin.Context, req *models.PluginReinstallRequest) error {
	if req.Node == "" {
		req.Node = "all"
	}
	if req.Module == "" {
		// 说明是服务端插件
		plg, flag := plugins.GlobalPluginManager.GetPlugin(req.Hash)
		if !flag {
			logger.Error(fmt.Sprintf("<UNK> plugin %s <UNK>", req.Hash))
			return fmt.Errorf("<UNK> plugin %s <UNK>", req.Hash)
		}
		plg.Install()
	} else {
		msg := models.Message{
			Name:    req.Node,
			Type:    "re_install_plugin",
			Content: fmt.Sprintf(`%v_%v`, req.Hash, req.Module),
		}
		err := s.nodeService.RefreshConfig(ctx, msg)
		if err != nil {
			logger.Error("failed to refresh config", zap.Error(err))
		}
	}
	return nil
}

// Recheck 重检插件
func (s *service) Recheck(ctx *gin.Context, req *models.PluginRecheckRequest) error {
	if req.Node == "" {
		req.Node = "all"
	}
	msg := models.Message{
		Name:    req.Node,
		Type:    "re_check_plugin",
		Content: fmt.Sprintf(`%v_%v`, req.Hash, req.Module),
	}
	err := s.nodeService.RefreshConfig(ctx, msg)
	if err != nil {
		logger.Error("failed to refresh config", zap.Error(err))
	}
	return nil
}

// Uninstall 卸载插件
func (s *service) Uninstall(ctx *gin.Context, req *models.PluginUninstallRequest) error {

	if req.Node == "" {
		req.Node = "all"
	}
	msg := models.Message{
		Name:    req.Node,
		Type:    "uninstall_plugin",
		Content: fmt.Sprintf(`%v_%v`, req.Hash, req.Module),
	}
	err := s.nodeService.RefreshConfig(ctx, msg)
	if err != nil {
		logger.Error("failed to refresh config", zap.Error(err))
	}
	return nil
}

// CheckKey 检查插件密钥
func (s *service) CheckKey(ctx *gin.Context, req *models.PluginKeyCheckRequest) error {
	key, err := ioutil.ReadFile("PLUGINKEY")
	if err != nil {
		return fmt.Errorf("failed to read plugin key: %w", err)
	}
	if req.Key != strings.TrimSpace(string(key)) {
		return fmt.Errorf("invalid plugin key")
	}
	return nil
}

// RemotePluginResponse 远程插件API响应结构
type RemotePluginResponse struct {
	Timestamp int64            `json:"timestamp"`
	Status    string           `json:"status"`
	Message   string           `json:"message"`
	Data      RemotePluginData `json:"data"`
}

type RemotePluginData struct {
	Total int64              `json:"total"`
	Data  []RemotePluginItem `json:"data"`
}

type RemotePluginItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Module       string `json:"module"`
	PriceStatus  int    `json:"priceStatus"`
	Price        *int   `json:"price"`
	Hash         string `json:"hash"`
	Introduction string `json:"introduction"`
	Version      string `json:"version"`
	CreateTime   string `json:"createTime"`
	Username     string `json:"username"`
	IsInstalled  bool   `json:"isInstalled"` // 是否已安装
	NeedUpdate   bool   `json:"needUpdate"`  // 是否需要更新
}

// SearchRemotePlugins 搜索远程插件并对比已安装的插件
func (s *service) SearchRemotePlugins(ctx *gin.Context) (map[string]interface{}, error) {
	var installedPlugins []models.Plugin
	var remotePlugins RemotePluginResponse
	var installedErr, remoteErr error
	var wg sync.WaitGroup

	// 并发获取已安装的插件和远程插件
	wg.Add(2)

	// 获取已安装的插件
	go func() {
		defer wg.Done()
		query := bson.M{}
		opts := mongoOptions.Find()
		installedPlugins, installedErr = s.repo.FindWithPagination(ctx, query, opts)
	}()

	// 获取远程插件
	go func() {
		defer wg.Done()
		client := &fasthttp.Client{
			ReadTimeout: 30 * time.Second,
		}

		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.SetRequestURI("https://api.scope-sentry.top/api/common/plugin/search")
		req.Header.SetMethod("POST")
		req.Header.SetContentType("application/json")

		// 构建请求体
		requestBody := map[string]interface{}{
			"name":        "",
			"priceStatus": 2,
			"module":      "All",
			"tags":        "",
			"page":        1,
			"size":        200,
		}
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			remoteErr = fmt.Errorf("failed to marshal request body: %w", err)
			return
		}
		req.SetBody(bodyBytes)

		err = client.Do(req, resp)
		if err != nil {
			remoteErr = fmt.Errorf("failed to request remote api: %w", err)
			return
		}

		if resp.StatusCode() != http.StatusOK {
			remoteErr = fmt.Errorf("remote api returned status code: %d", resp.StatusCode())
			return
		}

		if err := json.Unmarshal(resp.Body(), &remotePlugins); err != nil {
			remoteErr = fmt.Errorf("failed to unmarshal response: %w", err)
			return
		}
	}()

	wg.Wait()

	// 检查错误
	if installedErr != nil {
		return nil, fmt.Errorf("failed to get installed plugins: %w", installedErr)
	}
	if remoteErr != nil {
		return nil, fmt.Errorf("failed to get remote plugins: %w", remoteErr)
	}

	// 构建已安装插件的hash映射和version映射
	installedHashMap := make(map[string]models.Plugin)
	for _, plugin := range installedPlugins {
		installedHashMap[plugin.Hash] = plugin
	}

	// 为远程插件添加安装状态和更新状态
	resultPlugins := make([]RemotePluginItem, 0, len(remotePlugins.Data.Data))
	for _, remotePlugin := range remotePlugins.Data.Data {
		installedPlugin, isInstalled := installedHashMap[remotePlugin.Hash]

		remotePlugin.IsInstalled = isInstalled

		// 如果已安装，检查版本是否需要更新
		if isInstalled {
			// 只有当远程版本不为空且与已安装版本不同时，才需要更新
			remoteVersion := strings.TrimSpace(remotePlugin.Version)
			installedVersion := strings.TrimSpace(installedPlugin.Version)

			if remoteVersion != "" && remoteVersion != installedVersion {
				remotePlugin.NeedUpdate = true
			} else {
				remotePlugin.NeedUpdate = false
			}
		} else {
			remotePlugin.NeedUpdate = false
		}

		resultPlugins = append(resultPlugins, remotePlugin)
	}

	return map[string]interface{}{
		"total": remotePlugins.Data.Total,
		"data":  resultPlugins,
	}, nil
}

// ImportByData 通过POST JSON数据导入插件
func (s *service) ImportByData(ctx *gin.Context, req *models.PluginImportByDataRequest) error {
	// 验证密钥
	key, err := ioutil.ReadFile("PLUGINKEY")
	if err != nil {
		return fmt.Errorf("failed to read plugin key: %w", err)
	}
	if req.Key != strings.TrimSpace(string(key)) {
		return fmt.Errorf("invalid plugin key")
	}

	// 解析info.json的json字符串
	var pluginInfo models.PluginInfo
	if err := json.Unmarshal([]byte(req.JSON), &pluginInfo); err != nil {
		return fmt.Errorf("解析 json 字符串失败: %w", err)
	}

	// 校验关键信息
	if pluginInfo.Name == "" {
		return fmt.Errorf("json 中缺少 name 字段")
	}
	if pluginInfo.Type != "server" && pluginInfo.Module == "" {
		return fmt.Errorf("json 中缺少 Module 字段")
	}

	// 如果hash为空，生成新的hash
	if pluginInfo.Hash == "" {
		pluginInfo.Hash = generatePluginHash(32, pluginInfo.Type)
	}

	// 验证module是否合法
	if pluginInfo.Type != "server" {
		var pluginsModules = make(map[string]bool)
		for _, module := range constants.PLUGINSMODULES {
			pluginsModules[module] = true
		}
		if !pluginsModules[pluginInfo.Module] {
			return fmt.Errorf("模块非法: %s", pluginInfo.Module)
		}
	}

	var PLUGINS = make(map[string]bool)
	for _, p := range constants.Plugins {
		PLUGINS[p.Hash] = true
	}
	// 6. 插件是否是系统插件 或 module 不合法
	if PLUGINS[pluginInfo.Hash] {
		req.IsSystem = true
	}
	// 创建插件对象
	p := &models.Plugin{
		Name:          pluginInfo.Name,
		Module:        pluginInfo.Module,
		Hash:          pluginInfo.Hash,
		Parameter:     pluginInfo.Parameter,
		Help:          pluginInfo.Help,
		Introduction:  pluginInfo.Introduction,
		ParameterList: pluginInfo.ParameterList,
		Source:        req.Source,
		Version:       pluginInfo.Version,
		IsSystem:      req.IsSystem,
		Type:          pluginInfo.Type,
	}

	// 使用upsert：如果存在则更新，不存在则插入
	pluginID, err := s.repo.UpsertByHash(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to upsert plugin: %w", err)
	}

	// 如果不是内置插件，需要RefreshConfig
	if !req.IsSystem && pluginInfo.Type != "server" {
		go func() {
			msg := models.Message{
				Name:    "all",
				Type:    "install_plugin",
				Content: fmt.Sprintf(`%v`, pluginID),
			}
			err = s.nodeService.RefreshConfig(ctx, msg)
			if err != nil {
				logger.Error("failed to refresh config", zap.Error(err))
			}
		}()
	}
	if pluginInfo.Type == "server" {
		RegisterPlugin(req.Source, pluginInfo.Hash)
		plugins.RunPluginOnce(pluginInfo.Hash)
	}

	return nil
}

const pluginSalt = "ScopeSentry_"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generatePluginHash(length int, tp string) string {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(b)

	// 拼接盐值
	salted := randomString + pluginSalt + tp

	// 计算 MD5
	hash := md5.Sum([]byte(salted))

	// 返回十六进制字符串
	return hex.EncodeToString(hash[:])
}

// UpdateStatus 更新插件状态
func (s *service) UpdateStatus(ctx *gin.Context, req *models.PluginStatusRequest) error {

	update := bson.M{
		"status": req.Status,
	}
	err := s.repo.UpdateByHash(ctx, req.ID, update)
	if err != nil {
		return fmt.Errorf("failed to update plugin status: %w", err)
	}
	if req.Status {
		plg, flag := plugins.GlobalPluginManager.GetPlugin(req.ID)
		if !flag {
			logger.Error(fmt.Sprintf("failed to get plugin status: %v", req.ID))
			return fmt.Errorf("failed to get plugin: %v", req.ID)
		}
		if plg.Cycle() != "" {
			if plg.Cycle() == "1" {
				plugins.RunPluginOnce(plg.GetPluginId())
			} else {
				err := AddJob(plg.GetPluginId(), plg.Cycle())
				if err != nil {
					return err
				}
			}
		}
	} else {
		flag := utils.HasContext(req.ID)
		if flag {
			utils.CancelContext(req.ID)
		}
		err := schedulerCore.GetGlobalScheduler().RemoveJob(req.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to remove job: %v", err))
			return err
		}
	}
	return nil
}

// Run 运行插件一次
func (s *service) Run(ctx *gin.Context, req *models.PluginRunRequest) error {
	// 从 GlobalPluginManager 获取插件实例
	plgRunner, flag := plugins.GlobalPluginManager.GetPlugin(req.Hash)
	if !flag {
		// 从数据库获取插件信息
		plg, err := s.repo.FindByHash(ctx, req.Hash)
		if err != nil {
			return fmt.Errorf("failed to find plugin: %w", err)
		}
		if plg == nil {
			return fmt.Errorf("plugin not found")
		}

		// 只允许运行服务端插件
		if plg.Type != "server" {
			return fmt.Errorf("only server plugins can be run manually")
		}
		// 如果插件未加载，则加载它
		plugin, err := plugins.LoadPlugin(plg.Source, plg.Hash)
		if err != nil {
			return fmt.Errorf("failed to load plugin: %w", err)
		}
		plugins.GlobalPluginManager.RegisterPlugin(req.Hash, plugin)
		err = plugin.Install()
		if err != nil {
			logger.Error(fmt.Sprintf("failed to install plugin: %v", err))
			return fmt.Errorf("failed to install plugin: %w", err)
		}
		plgRunner = plugin.Clone()
	}

	plgRunner.Log("plugin is running manually")
	_, err := mongodb.UpdateOne("plugins", bson.M{"hash": req.Hash}, bson.M{"$set": bson.M{"lastTime": helper.GetNowTimeString()}})
	if err != nil {
		logger.Error("failed to update plugin", zap.Error(err))
	}
	// 执行插件
	go func() {
		err = plgRunner.Execute(options.NewOption(req.Hash))
		if err != nil {
			logger.Error(fmt.Sprintf("plugin execution error: %v", err))
			plgRunner.Log(fmt.Sprintf("plugin execution error: %v", err), "e")
		}
		plgRunner.Log("plugin execution completed")
	}()

	return nil
}

func (s *service) RunTaskEnd(task models.Task) error {
	query := bson.M{"type": "server", "status": true}
	logger.Info(fmt.Sprintf("executing task end plugin: %v", task.Name))
	opts := mongoOptions.Find()

	plgs, err := s.repo.FindWithPagination(context.Background(), query, opts)
	if err != nil {
		return fmt.Errorf("failed to find plugins: %w", err)
	}
	for _, plgTmp := range plgs {
		plg, flag := plugins.GlobalPluginManager.GetPlugin(plgTmp.Hash)
		if flag {
			go func() {
				err := plg.TaskEnd(task, plgTmp.Hash)
				if err != nil {
					logger.Error(fmt.Sprintf("plugin taskend execution error: %v", err))
					return
				}
			}()
		} else {
			logger.Error(fmt.Sprintf("plugin not found: %v", plgTmp.Hash))
		}
	}
	return nil
}
