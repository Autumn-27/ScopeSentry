// Package models -----------------------------
// @file      : task.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/4 22:04
// -------------------------------------------
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task 任务模型
type Task struct {
	ID             primitive.ObjectID       `bson:"_id,omitempty" json:"id"`
	TaskID         string                   `bson:"id"` // 关联的任务ID
	Name           string                   `bson:"name" json:"name"`
	Node           []string                 `bson:"node" json:"node"`
	AllNode        bool                     `bson:"allNode" json:"allNode"`
	ScheduledTasks bool                     `bson:"scheduledTasks" json:"scheduledTasks"`
	Target         string                   `bson:"target" json:"target"`
	Ignore         string                   `bson:"ignore" json:"ignore"`
	Template       string                   `bson:"template" json:"template"`
	Project        []string                 `bson:"project" json:"project"`
	Search         string                   `bson:"search" json:"search"`
	CycleType      string                   `bson:"cycleType" json:"cycleType"`
	Hour           int                      `bson:"hour" json:"hour"`
	Minute         int                      `bson:"minute" json:"minute"`
	Day            int                      `bson:"day" json:"day"`
	Week           int                      `bson:"week" json:"week"`
	Status         int                      `bson:"status" json:"status"`
	Progress       int                      `bson:"progress" json:"progress"`
	CreatTime      string                   `bson:"creatTime" json:"creatTime"`
	EndTime        string                   `bson:"endTime" json:"endTime"`
	TaskNum        int                      `bson:"taskNum" json:"taskNum"`
	Duplicates     string                   `bson:"duplicates" json:"duplicates"`
	Filter         map[string][]interface{} `bson:"filter" json:"filter"`
	TargetNumber   int                      `bson:"targetNumber" json:"targetNumber"`
	TargetIds      []string                 `bson:"targetIds" json:"targetIds"`
	TargetTp       string                   `bson:"targetTp" json:"targetTp"`
	TargetSource   string                   `bson:"targetSource" json:"targetSource"`
	BindProject    string                   `bson:"bindProject" json:"bindProject"`
	Type           string                   `bson:"type" json:"type"`
	IsStart        bool                     `bson:"isStart" json:"isStart"`
	LastTime       string                   `bson:"lastTime"  json:"lastTime"`
	NextTime       string                   `bson:"nextTime" json:"nextTime"`
	Cycle          string                   `bson:"-" json:"cycle"` // 周期描述（计算字段）
}

// TaskProgress 任务进度模型
type TaskProgress struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TargetHandler       int                `bson:"targetHandler" json:"targetHandler"`             // 目标处理
	SubdomainScan       int                `bson:"subdomainScan" json:"subdomainScan"`             // 子域名扫描
	SubdomainSecurity   int                `bson:"subdomainSecurity" json:"subdomainSecurity"`     // 子域名安全
	PortScanPreparation int                `bson:"portScanPreparation" json:"portScanPreparation"` // 端口扫描准备
	PortScan            int                `bson:"portScan" json:"portScan"`                       // 端口扫描
	PortFingerprint     int                `bson:"portFingerprint" json:"portFingerprint"`         // 端口指纹
	AssetMapping        int                `bson:"assetMapping" json:"assetMapping"`               // 资产映射
	AssetHandle         int                `bson:"assetHandle" json:"assetHandle"`                 // 资产处理
	URLScan             int                `bson:"urlScan" json:"urlScan"`                         // URL扫描
	WebCrawler          int                `bson:"webCrawler" json:"webCrawler"`                   // Web爬虫
	URLSecurity         int                `bson:"urlSecurity" json:"urlSecurity"`                 // URL安全
	DirScan             int                `bson:"dirScan" json:"dirScan"`                         // 目录扫描
	VulnerabilityScan   int                `bson:"vulnerabilityScan" json:"vulnerabilityScan"`     // 漏洞扫描
	All                 int                `bson:"all" json:"all"`                                 // 总进度
	Target              string             `bson:"target" json:"target"`                           // 目标
	Node                string             `bson:"node" json:"node"`                               // 节点
	Children            []TaskProgress     `bson:"children" json:"children"`                       // 子任务进度
}

// TaskProgressDetail 任务进度详情（用于API响应）
type TaskProgressDetail struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	TargetHandler       []string             `json:"TargetHandler"`       // [start, end]
	SubdomainScan       []string             `json:"SubdomainScan"`       // [start, end]
	SubdomainSecurity   []string             `json:"SubdomainSecurity"`   // [start, end]
	PortScanPreparation []string             `json:"PortScanPreparation"` // [start, end]
	PortScan            []string             `json:"PortScan"`            // [start, end]
	PortFingerprint     []string             `json:"PortFingerprint"`     // [start, end]
	AssetMapping        []string             `json:"AssetMapping"`        // [start, end]
	AssetHandle         []string             `json:"AssetHandle"`         // [start, end]
	URLScan             []string             `json:"URLScan"`             // [start, end]
	WebCrawler          []string             `json:"WebCrawler"`          // [start, end]
	URLSecurity         []string             `json:"URLSecurity"`         // [start, end]
	DirScan             []string             `json:"DirScan"`             // [start, end]
	VulnerabilityScan   []string             `json:"VulnerabilityScan"`   // [start, end]
	All                 []string             `json:"All"`                 // [start, end]
	Target              string               `json:"target"`
	Node                string               `json:"node"`
	Children            []TaskProgressDetail `json:"children,omitempty"`
}

// ListRequest 任务列表分页查询请求
type ListRequest struct {
	Search    string `json:"search" binding:"omitempty" example:"测试"`
	PageIndex int    `json:"pageIndex" binding:"required,min=1" example:"1"`
	PageSize  int    `json:"pageSize" binding:"required,min=1,max=100" example:"10"`
}

// ScheduledTask 计划任务模型
type ScheduledTask struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	TaskID         string             `bson:"id" json:"id"`                         // 关联的任务ID
	Name           string             `bson:"name" json:"name"`                     // 任务名称
	Type           string             `bson:"type" json:"type"`                     // 任务类型
	LastTime       string             `bson:"lastTime" json:"lastTime"`             // 上次执行时间
	NextTime       string             `bson:"nextTime" json:"nextTime"`             // 下次执行时间
	State          bool               `bson:"scheduledTasks" json:"state"`          // 是否启用
	Node           []string           `bson:"node" json:"node"`                     // 节点列表
	Cycle          string             `bson:"-" json:"cycle"`                       // 周期描述（计算字段）
	AllNode        bool               `bson:"allNode" json:"allNode"`               // 是否所有节点
	RunnerID       string             `bson:"runner_id" json:"runner_id"`           // 运行器ID
	Project        []string           `bson:"project" json:"project"`               // 项目列表
	TargetSource   string             `bson:"targetSource" json:"targetSource"`     // 目标来源
	Day            int                `bson:"day" json:"day"`                       // 天数
	Minute         int                `bson:"minute" json:"minute"`                 // 分钟
	Hour           int                `bson:"hour" json:"hour"`                     // 小时
	Search         string             `bson:"search" json:"search"`                 // 搜索条件
	CycleType      string             `bson:"cycleType" json:"cycleType"`           // 周期类型
	ScheduledTasks bool               `bson:"scheduledTasks" json:"scheduledTasks"` // 是否计划任务
	Target         string             `bson:"target" json:"target"`                 // 目标
	Ignore         string             `bson:"ignore" json:"ignore"`                 // 忽略规则
	Template       string             `bson:"template" json:"template"`             // 模板
	Duplicates     string             `bson:"duplicates" json:"duplicates"`         // 去重规则
	Week           int                `bson:"week" json:"week"`                     // 星期
}

// ScheduledTaskListRequest 计划任务列表请求
type ScheduledTaskListRequest struct {
	Search    string `json:"search" binding:"omitempty" example:"测试"`
	PageIndex int    `json:"pageIndex" binding:"omitempty,min=1" example:"1"`
	PageSize  int    `json:"pageSize" binding:"omitempty,min=1,max=100" example:"10"`
}

// ScheduledTaskDetailRequest 计划任务详情请求
type ScheduledTaskDetailRequest struct {
	ID string `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
}

// ScheduledTaskDeleteRequest 计划任务删除请求
type ScheduledTaskDeleteRequest struct {
	IDs []string `json:"ids" binding:"required" example:"[\"507f1f77bcf86cd799439011\"]"`
}

// ScheduledTaskAddRequest 计划任务添加请求
type ScheduledTaskAddRequest struct {
	Name           string   `json:"name" binding:"required" example:"定时扫描任务"`
	Type           string   `json:"type" binding:"omitempty" example:"scan"`
	Node           []string `json:"node" binding:"required" example:"[\"node1\", \"node2\"]"`
	AllNode        bool     `json:"allNode" binding:"omitempty" example:"true"`
	Project        []string `json:"project" binding:"omitempty" example:"[\"project1\"]"`
	TargetSource   string   `json:"targetSource" binding:"omitempty" example:"general"`
	Day            int      `json:"day" binding:"omitempty" example:"1"`
	Minute         int      `json:"minute" binding:"omitempty" example:"0"`
	Hour           int      `json:"hour" binding:"omitempty" example:"1"`
	Search         string   `json:"search" binding:"omitempty" example:""`
	CycleType      string   `json:"cycleType" binding:"omitempty" example:"nhours"`
	ScheduledTasks bool     `json:"scheduledTasks" binding:"omitempty" example:"true"`
	Target         string   `json:"target" binding:"omitempty" example:"example.com"`
	Ignore         string   `json:"ignore" binding:"omitempty" example:""`
	Template       string   `json:"template" binding:"omitempty" example:"default"`
	Duplicates     string   `json:"duplicates" binding:"omitempty" example:"true"`
	Week           int      `json:"week" binding:"omitempty" example:"1"`
}

// ScheduledTaskUpdateRequest 计划任务更新请求
type ScheduledTaskUpdateRequest struct {
	ID             string   `json:"id" binding:"required" example:"507f1f77bcf86cd799439011"`
	Name           string   `json:"name" binding:"omitempty" example:"定时扫描任务"`
	Type           string   `json:"type" binding:"omitempty" example:"scan"`
	Node           []string `json:"node" binding:"omitempty" example:"[\"node1\", \"node2\"]"`
	AllNode        bool     `json:"allNode" binding:"omitempty" example:"true"`
	Project        []string `json:"project" binding:"omitempty" example:"[\"project1\"]"`
	TargetSource   string   `json:"targetSource" binding:"omitempty" example:"general"`
	Day            int      `json:"day" binding:"omitempty" example:"1"`
	Minute         int      `json:"minute" binding:"omitempty" example:"0"`
	Hour           int      `json:"hour" binding:"omitempty" example:"1"`
	Search         string   `json:"search" binding:"omitempty" example:""`
	CycleType      string   `json:"cycleType" binding:"omitempty" example:"nhours"`
	ScheduledTasks bool     `json:"scheduledTasks" binding:"omitempty" example:"true"`
	Target         string   `json:"target" binding:"omitempty" example:"example.com"`
	Ignore         string   `json:"ignore" binding:"omitempty" example:""`
	Template       string   `json:"template" binding:"omitempty" example:"default"`
	Duplicates     string   `json:"duplicates" binding:"omitempty" example:"true"`
	Week           int      `json:"week" binding:"omitempty" example:"1"`
}

// PageMonitoringTask 页面监控任务模型
type PageMonitoringTask struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL     string             `bson:"url" json:"url"`
	Hash    []string           `bson:"hash" json:"hash"`
	MD5     string             `bson:"md5" json:"md5"`
	State   int                `bson:"state" json:"state"`
	Project string             `bson:"project" json:"project"`
	Time    string             `bson:"time" json:"time"`
}

// PageMonitoringListRequest 页面监控列表请求
type PageMonitoringListRequest struct {
	Search    string `json:"search" binding:"omitempty" example:"example.com"`
	PageIndex int    `json:"pageIndex" binding:"omitempty,min=1" example:"1"`
	PageSize  int    `json:"pageSize" binding:"omitempty,min=1,max=100" example:"10"`
}

// PageMonitoringAddRequest 页面监控添加请求
type PageMonitoringAddRequest struct {
	URL string `json:"url" binding:"required" example:"https://example.com"`
}

// PageMonitoringDeleteRequest 页面监控删除请求
type PageMonitoringDeleteRequest struct {
	IDs []string `json:"ids" binding:"required" example:"[\"507f1f77bcf86cd799439011\"]"`
}

// PageMonitoringUpdateRequest 页面监控更新请求
type PageMonitoringUpdateRequest struct {
	Hour    int      `json:"hour" binding:"omitempty" example:"24"`
	Node    []string `json:"node" binding:"omitempty" example:"[\"node1\"]"`
	AllNode bool     `json:"allNode" binding:"omitempty" example:"true"`
	State   bool     `json:"state" binding:"omitempty" example:"true"`
}
