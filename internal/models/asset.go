package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Asset 表示一个资产的基本信息
type Asset struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Time             string             `bson:"time" json:"time"`
	LastScanTime     string             `bson:"lastScanTime" json:"lastScanTime"`
	Host             string             `bson:"host" json:"domain"`
	IP               string             `bson:"ip" json:"ip"`
	Port             string             `bson:"port" json:"port"`
	Service          string             `bson:"service" json:"service"`
	Type             string             `bson:"type" json:"type"`
	URL              string             `bson:"url,omitempty" json:"url"`
	Title            string             `bson:"title,omitempty" json:"title"`
	StatusCode       int                `bson:"statuscode,omitempty" json:"status"`
	ContentLength    int64              `bson:"contentlength,omitempty" json:"contentlength,omitempty"`
	ResponseBodyHash string             `bson:"bodyhash" csv:"bodyhash"`
	Technologies     []string           `bson:"technologies,omitempty" json:"products"`
	RawHeaders       string             `bson:"rawheaders,omitempty" json:"banner"`
	Banner           string             `bson:"metadata" json:"metadata"`
	Body             string             `bson:"body,omitempty" json:"body,omitempty"`
	Screenshot       string             `bson:"screenshot,omitempty" json:"screenshot,omitempty"`
	FaviconMMH3      string             `bson:"faviconmmh3,omitempty" json:"faviconmmh3,omitempty"`
	JARM             string             `bson:"jarm,omitempty" json:"jarm,omitempty"`
	TLSData          map[string]string  `bson:"tlsdata,omitempty" json:"tlsdata,omitempty"`
	Hashes           map[string]string  `bson:"hashes,omitempty" json:"hashes,omitempty"`
	Project          string             `bson:"project,omitempty" json:"project,omitempty"`
	TaskName         []string           `bson:"taskName,omitempty" json:"taskName,omitempty"`
	Tags             []string           `bson:"tags,omitempty" json:"tags"`
	RootDomain       string             `bson:"rootDomain,omitempty" json:"rootDomain,omitempty"`
	CDNName          string             `bson:"cdnname,omitempty" json:"cdnname,omitempty"`
	Error            string             `bson:"error,omitempty" json:"error,omitempty"`
	FaviconPath      string             `bson:"faviconpath,omitempty" json:"faviconpath,omitempty"`
	IconContent      string             `bson:"iconcontent,omitempty" json:"icon"`
	Webcheck         bool               `bson:"webcheck,omitempty" json:"webcheck,omitempty"`
	WebServer        string             `bson:"webServer,omitempty" json:"webServer,omitempty"`
}

type HttpBody struct {
	BodyHash string `bson:"hash,omitempty" json:"hash,omitempty"`
	Content  string `bson:"content,omitempty" json:"content,omitempty"`
}

// AssetChangeLog 表示资产变更日志
type FieldChange struct {
	FieldName string `bson:"fieldname" json:"fieldname"`
	Old       string `bson:"old" json:"old"`
	New       string `bson:"new" json:"new"`
}

type AssetChangeLog struct {
	AssetID   string        `bson:"assetid" json:"assetid"`
	Timestamp string        `bson:"timestamp" json:"timestamp"`
	Change    []FieldChange `bson:"change" json:"change"`
}

// AssetSearchQuery 表示资产搜索查询参数
type AssetSearchQuery struct {
	PageIndex  int                    `json:"pageIndex"`
	PageSize   int                    `json:"pageSize"`
	Filters    map[string]interface{} `json:"filters"`
	SortField  string                 `json:"sortField"`
	SortOrder  int                    `json:"sortOrder"`
	SearchText string                 `json:"searchText"`
}

// AssetResponse 表示资产API响应
type AssetResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type IPAssetTmp struct {
	ID           primitive.ObjectID `bson:"_id"`
	IP           string             `json:"ip" bson:"ip"`
	Domain       string             `json:"domain" bson:"domain"`
	Port         string             `json:"port" bson:"port"`
	Service      string             `json:"service" bson:"service"`
	WebServer    string             `json:"webServer" bson:"webServer"`
	Technologies []string           `json:"technologies" bson:"technologies"`
	Project      string             `json:"project" bson:"project"`
	TaskName     []string           `json:"taskName" bson:"taskName"`
	RootDomain   string             `bson:"rootDomain,omitempty" json:"rootDomain,omitempty"`
}

type IPAsset struct {
	IP         string    `json:"ip" bson:"ip"`
	Ports      []IPPorts `json:"ports" bson:"ports"`
	Project    []string  `json:"project" bson:"project"`
	TaskName   []string  `json:"taskName" bson:"taskName"`
	RootDomain []string  `bson:"rootDomain,omitempty" json:"rootDomain,omitempty"`
	Time       string    `bson:"time" csv:"time"`
}

type IPAssetFlat struct {
	IP          string   `json:"ip"`
	IPRowSpan   int      `json:"ipRowSpan"`
	Port        string   `json:"port"`
	PortRowSpan int      `json:"portRowSpan"`
	Domain      string   `json:"domain"`
	Service     string   `json:"service"`
	WebServer   string   `json:"webServer"`
	Products    []string `json:"products"`
	Time        string   `json:"time"`
}

type IPAssetResponse struct {
	List  []IPAssetFlat `json:"list"`
	Total int64         `json:"total"`
}

type IPPorts struct {
	Port   string       `json:"port" bson:"port"`
	Server []PortServer `json:"server" bson:"server"`
}

type PortServer struct {
	Domain       string   `json:"domain" bson:"domain"`
	Service      string   `json:"service" bson:"service"`
	WebServer    string   `json:"webServer" bson:"webServer"`
	Technologies []string `json:"products" bson:"technologies"`
}
