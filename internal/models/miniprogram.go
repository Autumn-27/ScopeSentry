package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MiniProgram 表示一个小程序
type MiniProgram struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`               // 小程序名称
	AppID       string             `bson:"appid" json:"appid"`             // 小程序AppID
	Type        string             `bson:"type" json:"type"`               // 小程序类型（微信/支付宝/百度等）
	Version     string             `bson:"version" json:"version"`         // 版本号
	Icon        string             `bson:"icon" json:"icon"`               // 小程序图标URL
	Description string             `bson:"description" json:"description"` // 小程序描述
	Developer   string             `bson:"developer" json:"developer"`     // 开发者
	Category    string             `bson:"category" json:"category"`       // 分类
	Tags        []string           `bson:"tags" json:"tags"`               // 标签
	Status      string             `bson:"status" json:"status"`           // 状态（上线/下线/审核中等）
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`   // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`   // 更新时间
}
