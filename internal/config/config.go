package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils"

	"github.com/Autumn-27/ScopeSentry/internal/utils/random"
	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	System  SystemConfig  `mapstructure:"system"`
	Server  ServerConfig  `mapstructure:"server"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	Redis   RedisConfig   `mapstructure:"redis"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Log     LogConfig     `mapstructure:"logs"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Timezone      string `mapstructure:"timezone"`
	PluginKey     string `mapstructure:"plugin_key"`
	SecretKey     string `mapstructure:"secret_key"`
	Debug         bool   `mapstructure:"debug"`
	ExeDir        string `yaml:"-"`
	IconDir       string `yaml:"-"`
	ScreenshotDir string `yaml:"-"`
	ImgDir        string `yaml:"-"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	IP       string `mapstructure:"ip"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string        `mapstructure:"secret"`
	Expire time.Duration `mapstructure:"expire"`
}

// LogConfig 日志配置
type LogConfig struct {
	TotalLogs int    `mapstructure:"total_logs"`
	Filename  string `mapstructure:"file_name"`
	Level     string `mapstructure:"leval"`
}

// GlobalConfig 全局配置实例
var GlobalConfig Config

// init 加载配置
func init() {
	// 设置默认值
	setDefaults()

	// 设置配置文件路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	exeDir := filepath.Dir(exePath)
	GlobalConfig.System.ExeDir = exeDir

	err = utils.EnsureDir(filepath.Join(exeDir, "files", "export"))
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	// 创建上传目录
	err = utils.EnsureDir(filepath.Join(exeDir, "uploads"))
	if err != nil {
		fmt.Printf("Error creating uploads directory: %s\n", err)
		return
	}
	// 创建图片目录
	GlobalConfig.System.ImgDir = filepath.Join(exeDir, "images")
	err = utils.EnsureDir(GlobalConfig.System.ImgDir)
	if err != nil {
		fmt.Printf("Error creating images directory: %s\n", err)
		return
	}
	GlobalConfig.System.IconDir = filepath.Join(GlobalConfig.System.ImgDir, "icon")
	err = utils.EnsureDir(GlobalConfig.System.IconDir)
	if err != nil {
		fmt.Printf("Error creating icon directory: %s\n", err)
		return
	}
	GlobalConfig.System.ScreenshotDir = filepath.Join(GlobalConfig.System.ImgDir, "screenshots")
	err = utils.EnsureDir(GlobalConfig.System.ScreenshotDir)
	if err != nil {
		fmt.Printf("Error creating screenshots directory: %s\n", err)
		return
	}

	viper.AddConfigPath(exeDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 设置环境变量前缀和分隔符
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			fmt.Printf("error reading config file: %v\n", err)
			return
		}
	}
	// 从环境变量加载配置
	loadFromEnv()

	// 解析配置到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		fmt.Printf("error unmarshaling config: %v\n", err)
		return
	}

	// 处理插件密钥
	if err := handlePluginKey(); err != nil {
		fmt.Println(err)
		return
	}
	if err := createConfigFile(); err != nil {
		fmt.Printf("error creating config file: %v\n", err)
		return
	}

	return
}

// setDefaults 设置默认配置
func setDefaults() {
	// 系统配置
	viper.SetDefault("system.timezone", "Asia/Shanghai")
	viper.SetDefault("system.debug", false)

	viper.SetDefault("jwt.secret", random.GenerateString(32))
	viper.SetDefault("jwt.expire", time.Duration(24)*time.Hour)
	// 服务器配置
	viper.SetDefault("server.port", 8082)
	viper.SetDefault("server.mode", "info")
	viper.SetDefault("server.read_timeout", 60*time.Second)
	viper.SetDefault("server.write_timeout", 60*time.Second)

	// MongoDB配置
	viper.SetDefault("mongodb.ip", "127.0.0.1")
	viper.SetDefault("mongodb.port", 27017)
	viper.SetDefault("mongodb.database", "ScopeSentry")
	viper.SetDefault("mongodb.username", "")
	viper.SetDefault("mongodb.password", "")

	// Redis配置
	viper.SetDefault("redis.ip", "127.0.0.1")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")

	// 日志配置
	viper.SetDefault("logs.total_logs", 1000)
	viper.SetDefault("logs.level", "info")
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv() {
	// 时区配置
	if timezone := os.Getenv("TIMEZONE"); timezone != "" {
		viper.Set("system.timezone", timezone)
	}

	// MongoDB 配置
	if ip := os.Getenv("MONGODB_IP"); ip != "" {
		viper.Set("mongodb.ip", ip)
	}
	if port := os.Getenv("MONGODB_PORT"); port != "" {
		viper.Set("mongodb.port", port)
	}
	if database := os.Getenv("MONGODB_DATABASE"); database != "" {
		viper.Set("mongodb.database", database)
	}
	if username := os.Getenv("MONGODB_USER"); username != "" {
		viper.Set("mongodb.username", username)
	}
	if password := os.Getenv("MONGODB_PASSWORD"); password != "" {
		viper.Set("mongodb.password", password)
	}

	// Redis 配置
	if ip := os.Getenv("REDIS_IP"); ip != "" {
		viper.Set("redis.ip", ip)
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		viper.Set("redis.port", port)
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		viper.Set("redis.password", password)
	}
}

// createConfigFile 创建配置文件
func createConfigFile() error {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// 设置配置文件路径
	configPath := filepath.Join(exeDir, "config.yaml")

	// 确保配置目录存在
	if err := os.MkdirAll(exeDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	// 写入配置文件
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// handlePluginKey 处理插件密钥
func handlePluginKey() error {
	pluginKeyFile := "PLUGINKEY"
	if _, err := os.Stat(pluginKeyFile); err == nil {
		// 文件存在，读取密钥
		data, err := os.ReadFile(pluginKeyFile)
		if err != nil {
			return err
		}
		GlobalConfig.System.PluginKey = string(data)
	} else {
		// 文件不存在，生成新密钥
		key := random.GeneratePluginKey()
		if err := os.WriteFile(pluginKeyFile, []byte(key), 0644); err != nil {
			return err
		}
		GlobalConfig.System.PluginKey = key
	}
	return nil
}

// GetMongoURI 获取MongoDB连接URI
func GetMongoURI() string {
	uri := GlobalConfig.MongoDB.IP + ":" + string(GlobalConfig.MongoDB.Port)
	if GlobalConfig.MongoDB.Username != "" && GlobalConfig.MongoDB.Password != "" {
		uri = "mongodb://" + GlobalConfig.MongoDB.Username + ":" + GlobalConfig.MongoDB.Password + "@" + uri
	}
	return uri
}

// GetRedisAddr 获取Redis连接地址
func GetRedisAddr() string {
	return GlobalConfig.Redis.IP + ":" + GlobalConfig.Redis.Port
}

// GetPluginKey 获取插件密钥
func GetPluginKey() string {
	return GlobalConfig.System.PluginKey
}

// GetSecretKey 获取密钥
func GetSecretKey() string {
	return GlobalConfig.System.SecretKey
}

// GetTimezone 获取时区
func GetTimezone() string {
	return GlobalConfig.System.Timezone
}

// GetTotalLogs 获取日志总数限制
func GetTotalLogs() int {
	return GlobalConfig.Log.TotalLogs
}
