package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subdomain 表示一个子域名
type Subdomain struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Host       string             `bson:"host" json:"host"`             // 主机名
	Type       string             `bson:"type" json:"type"`             // 类型
	Value      []string           `bson:"value" json:"value"`           // 值列表
	IP         []string           `bson:"ip" json:"ip"`                 // IP地址列表
	Time       string             `bson:"time" json:"time"`             // 时间
	Tags       []string           `bson:"tags" json:"tags"`             // 标签
	Project    string             `bson:"project" json:"project"`       // 项目
	TaskName   string             `bson:"taskName" json:"taskName"`     // 任务名称
	RootDomain string             `bson:"rootDomain" json:"rootDomain"` // 根域名
}

// SubdomainResponse 子域名响应数据结构
type SubdomainResponse struct {
	List  []Subdomain `json:"list"`
	Total int64       `json:"total"`
}

type SubdomainTakerResult struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Input    string             `bson:"input" json:"host"`
	Value    string             `bson:"value" json:"value"`
	Cname    string             `bson:"cname" json:"type"`
	Response string             `bson:"response" json:"response"`
	Tags     []string           `bson:"tags" json:"tags"`
}

type SubdomainTakerResponse struct {
	List  []SubdomainTakerResult `json:"list"`
	Total int64                  `json:"total"`
}
