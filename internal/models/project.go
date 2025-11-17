package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Project 项目模型
type Project struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Tag            string             `bson:"tag" json:"tag"`
	Logo           string             `bson:"logo" json:"logo"`
	Target         string             `bson:"target" json:"target"`
	ScheduledTasks bool               `bson:"scheduledTasks" json:"scheduledTasks"`
	Hour           int                `bson:"hour" json:"hour"`
	AllNode        bool               `bson:"allNode" json:"allNode"`
	Node           []string           `bson:"node" json:"node"`
	Duplicates     string             `bson:"duplicates" json:"duplicates"`
	Ignore         string             `bson:"ignore" json:"ignore"`
	Template       string             `bson:"template" json:"template"`
	Tp             string             `bson:"tp" json:"tp"`
	RootDomains    []string           `bson:"root_domains" json:"rootDomains"`
	AssetCount     int                `bson:"AssetCount" json:"assetCount"`
}

type UpdateProject struct {
	ID             string   `bson:"id,omitempty" json:"id"`
	Name           string   `bson:"name" json:"name"`
	Tag            string   `bson:"tag" json:"tag"`
	Logo           string   `bson:"logo" json:"logo"`
	Target         string   `bson:"target" json:"target"`
	ScheduledTasks bool     `bson:"scheduledTasks" json:"scheduledTasks"`
	Hour           int      `bson:"hour" json:"hour"`
	AllNode        bool     `bson:"allNode" json:"allNode"`
	Node           []string `bson:"node" json:"node"`
	Duplicates     string   `bson:"duplicates" json:"duplicates"`
	Ignore         string   `bson:"ignore" json:"ignore"`
	Template       string   `bson:"template" json:"template"`
	Tp             string   `bson:"tp" json:"tp"`
	RootDomains    []string `bson:"root_domains" json:"rootDomains"`
	AssetCount     int      `bson:"AssetCount" json:"assetCount"`
}

// TagGroup 标签分组响应
type TagGroup struct {
	Label    string       `json:"label"`
	Value    string       `json:"value"`
	Children []TagProject `json:"children"`
}

// TagProject 标签下的项目
type TagProject struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ProjectListRequest 请求（定义在 handler 也可，此处作为复用）
type ProjectListRequest struct {
	Search    string `json:"search"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

// ProjectBrief 列表项概要
type ProjectBrief struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Logo       string `json:"logo"`
	AssetCount int    `json:"AssetCount"`
	Tag        string `json:"tag"`
}

// ProjectListResponse 列表响应数据
type ProjectListResponse struct {
	Result map[string][]ProjectBrief `json:"result"`
	Tag    map[string]int            `json:"tag"`
}

// ProjectContentRequest 请求
type ProjectContentRequest struct {
	ID string `json:"id"`
}

// ProjectContentResponse 内容响应
type ProjectContentResponse struct {
	Name           string   `json:"name"`
	Tag            string   `json:"tag"`
	Target         string   `json:"target"`
	Node           []string `json:"node"`
	Logo           string   `json:"logo"`
	ScheduledTasks bool     `json:"scheduledTasks"`
	Hour           int      `json:"hour"`
	AllNode        bool     `json:"allNode"`
	Duplicates     string   `json:"duplicates"`
	Template       string   `json:"template"`
	Ignore         string   `json:"ignore"`
}
