package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApiKey API Key 模型，用于 MCP 等外部集成认证
type ApiKey struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	KeyHash    string             `bson:"keyHash" json:"-"`
	KeyPrefix  string             `bson:"keyPrefix" json:"keyPrefix"`
	Enabled    bool               `bson:"enabled" json:"enabled"`
	CreatedBy  string             `bson:"createdBy" json:"createdBy"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	LastUsedAt *time.Time         `bson:"lastUsedAt,omitempty" json:"lastUsedAt,omitempty"`
}

// CreateApiKeyRequest 创建 API Key 请求
type CreateApiKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateApiKeyResponse 创建 API Key 响应（明文 key 仅返回一次）
type CreateApiKeyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	KeyPrefix string `json:"keyPrefix"`
	CreatedAt string `json:"createdAt"`
}

// DeleteApiKeyRequest 删除 API Key 请求
type DeleteApiKeyRequest struct {
	ID string `json:"id" binding:"required"`
}
