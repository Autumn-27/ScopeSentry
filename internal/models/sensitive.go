package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SensitiveInfo 表示一条敏感信息记录
type SensitiveChild struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Match  []string           `bson:"match" json:"match"`
	Time   string             `bson:"time" json:"time"`
	Url    string             `bson:"url" json:"url"`
	Tags   []string           `bson:"tags" json:"tags"`
	Status int                `bson:"status" json:"status"`
	Name   string             `bson:"name" json:"name"`
	Color  string             `bson:"color" json:"color"`
}

type Sensitive struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Match  []string           `bson:"match" json:"match"`
	Time   string             `bson:"time" json:"time"`
	Url    string             `bson:"url" json:"url"`
	Tags   []string           `bson:"tags" json:"tags"`
	Status int                `bson:"status" json:"status"`
	Name   string             `bson:"sid" json:"name"`
	Color  string             `bson:"color" json:"color"`
	BodyId string             `bson:"md5" json:"body_id"`
}

type SensitiveInfo struct {
	ID       string           `bson:"_id" json:"id"`
	Time     string           `bson:"time" json:"time"`
	Url      string           `bson:"url" json:"url"`
	BodyID   string           `bson:"body_id" json:"body_id"`
	Children []SensitiveChild `bson:"children" json:"children"`
}

// SensitiveRule 表示一条敏感信息检测规则
type SensitiveRule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`               // 规则名称
	Type        string             `bson:"type" json:"type"`               // 规则类型
	Pattern     string             `bson:"pattern" json:"pattern"`         // 匹配模式（正则表达式）
	Description string             `bson:"description" json:"description"` // 规则描述
	Level       string             `bson:"level" json:"level"`             // 敏感等级（high/medium/low）
	State       string             `bson:"state" json:"state"`             // 规则状态（enabled/disabled）
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`   // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`   // 更新时间
}

// SensitiveRuleItem 用于前端规则管理（映射 Python 的 SensitiveRule 集合）
type SensitiveRuleItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Regular string             `bson:"regular" json:"regular"`
	Color   string             `bson:"color" json:"color"`
	State   bool               `bson:"state" json:"state"`
}

type SensitiveRuleListItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Regular string `json:"regular"`
	Color   string `json:"color"`
	State   bool   `json:"state"`
}

type GetSensitiveBodyRequest struct {
	ID string `json:"id" binding:"required" example:"d41d8cd98f00b204e9800998ecf8427e"`
}

type GetSensitiveBodyResponse struct {
	Body string `json:"body" example:"<html>content</html>"`
}

type SensitiveSIDStat struct {
	Name  string `json:"name" bson:"name"`
	Count int    `json:"count" bson:"count"`
	Color string `json:"color" bson:"color"`
}

type GetSensitiveMatchInfoRequest struct {
	Sid string `json:"sid" binding:"required" example:"1234"`
}

type MatchListResponse struct {
	List []string `json:"list"`
}
