package models

// DeleteRequest 删除请求模型
type DeleteRequest struct {
	IDs           []string `json:"ids" binding:"required"`
	Index         string   `json:"index"`
	DeleateAssets bool     `json:"delA"`
}

// TagRequest 标签请求模型
type TagRequest struct {
	Type string `json:"tp" binding:"required"`
	ID   string `json:"id" binding:"required"`
	Tag  string `json:"tag" binding:"required"`
}

// StatusRequest 状态更新请求模型
type StatusRequest struct {
	Type   string `json:"tp" binding:"required"`
	ID     string `json:"id" binding:"required"`
	Status int    `json:"status" binding:"required"`
}

// TotalRequest 总数请求模型
type TotalRequest struct {
	Index string `json:"index" binding:"required"`
}

// Document 文档模型
type Document struct {
	ID     string   `bson:"_id,omitempty" json:"id,omitempty"`
	Tags   []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Status bool     `bson:"status,omitempty" json:"status,omitempty"`
}

type IdRequest struct {
	ID string `json:"id" `
}
