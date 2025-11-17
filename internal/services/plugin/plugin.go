package plugin

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/constants"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/plugin"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/node"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Service 定义插件服务接口
type Service interface {
	List(ctx *gin.Context, req *models.PluginListRequest) (*models.PluginListResponse, error)
	ListByModule(ctx *gin.Context, module string) ([]models.Plugin, error)
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
	if req.Search != "" {
		query["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	opts := options.Find().
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
	return s.repo.FindByModule(ctx, module)
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
			Name:         req.Name,
			Module:       req.Module,
			Hash:         generatePluginHash(32),
			Parameter:    req.Parameter,
			Help:         req.Help,
			Introduction: req.Introduction,
			Source:       req.Source,
			Version:      req.Version,
			IsSystem:     false,
		}
		id, err := s.repo.Create(ctx, plugin)
		if err != nil {
			return fmt.Errorf("failed to create plugin: %w", err)
		}
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
		return nil
	} else {
		// 更新现有插件
		id, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			return fmt.Errorf("invalid id format: %w", err)
		}

		update := bson.M{
			"name":         req.Name,
			"module":       req.Module,
			"parameter":    req.Parameter,
			"help":         req.Help,
			"introduction": req.Introduction,
			"source":       req.Source,
			"version":      req.Version,
		}
		err = s.repo.Update(ctx, id, update)
		if err != nil {
			return err
		}

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

		return nil
	}
}

// Delete 删除插件
func (s *service) Delete(ctx *gin.Context, req *models.PluginDeleteRequest) error {
	var hashes []string
	for _, item := range req.Data {
		hashes = append(hashes, item.Hash)
		go func() {
			msg := models.Message{
				Name:    "all",
				Type:    "delete_plugin",
				Content: fmt.Sprintf(`%v_%v`, item.Hash, item.Module),
			}
			err := s.nodeService.RefreshConfig(ctx, msg)
			if err != nil {
				logger.Error("failed to refresh config", zap.Error(err))
			}
		}()
	}
	return s.repo.DeleteByHash(ctx, hashes)
}

// GetLogs 获取插件日志
func (s *service) GetLogs(ctx *gin.Context, req *models.PluginLogRequest) (string, error) {
	if req.Module == "" || req.Hash == "" {
		return "", fmt.Errorf("module and hash are required")
	}
	logKey := fmt.Sprintf("logs:plugins:%v:%v", req.Module, req.Hash)
	logs, err := s.repo.GetLogs(ctx, logKey)
	if err != nil {
		return "", err
	}
	return logs, nil
}

// CleanLogs 清理插件日志
func (s *service) CleanLogs(ctx *gin.Context, req *models.PluginLogRequest) error {
	if req.Module == "" || req.Hash == "" {
		return fmt.Errorf("module and hash are required")
	}
	logKey := fmt.Sprintf("logs:plugins:%v:%v", req.Module, req.Hash)
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
	if pluginInfo.Name == "" || pluginInfo.Module == "" {
		return fmt.Errorf("info.json 缺少 name 或 module 字段")
	}
	if pluginInfo.Hash == "" {
		pluginInfo.Hash = generatePluginHash(32)
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
	if !pluginsModules[pluginInfo.Module] {
		return fmt.Errorf("模块非法: %s", pluginInfo.Module)
	}

	// 7. 设置其他字段
	pluginInfo.Source = pluginSource
	pluginInfo.IsSystem = false

	// 8. 插入数据库

	// 9. 触发配置刷新

	p := &models.Plugin{
		Name:         pluginInfo.Name,
		Module:       pluginInfo.Module,
		Hash:         pluginInfo.Hash,
		Parameter:    pluginInfo.Parameter,
		Help:         pluginInfo.Help,
		Introduction: pluginInfo.Introduction,
		Source:       pluginInfo.Source,
		Version:      pluginInfo.Version,
		IsSystem:     false,
	}
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to create plugin: %w", err)
	}
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
	return nil
}

// Reinstall 重装插件
func (s *service) Reinstall(ctx *gin.Context, req *models.PluginReinstallRequest) error {
	if req.Node == "" {
		req.Node = "all"
	}
	msg := models.Message{
		Name:    req.Node,
		Type:    "re_install_plugin",
		Content: fmt.Sprintf(`%v_%v`, req.Hash, req.Module),
	}
	err := s.nodeService.RefreshConfig(ctx, msg)
	if err != nil {
		logger.Error("failed to refresh config", zap.Error(err))
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

const pluginSalt = "ScopeSentry"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generatePluginHash(length int) string {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机字符串
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(b)

	// 拼接盐值
	salted := randomString + pluginSalt

	// 计算 MD5
	hash := md5.Sum([]byte(salted))

	// 返回十六进制字符串
	return hex.EncodeToString(hash[:])
}
