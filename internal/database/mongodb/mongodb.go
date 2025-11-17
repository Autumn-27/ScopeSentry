package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/logger"

	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func init() {
	// 设置连接超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	username := url.QueryEscape(config.GlobalConfig.MongoDB.Username)
	password := url.QueryEscape(config.GlobalConfig.MongoDB.Password)

	uri := config.GlobalConfig.MongoDB.IP + ":" + strconv.Itoa(config.GlobalConfig.MongoDB.Port)
	if username != "" && password != "" {
		uri = "mongodb://" + username + ":" + password + "@" + uri
	}

	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(uri)

	// 连接到 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect %v to MongoDB: %v", uri, err))
		return
	}

	// 检查连接
	if err = client.Ping(ctx, nil); err != nil {
		logger.Error(fmt.Sprintf("failed to ping MongoDB: %v", err))
		return
	}

	Client = client
	DB = client.Database(config.GlobalConfig.MongoDB.Database)

	logger.Info("MongoDB connected successfully")
	return
}

func Close() error {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := Client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB: %v", err)
		}
	}
	return nil
}
