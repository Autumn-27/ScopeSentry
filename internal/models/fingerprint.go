// models-------------------------------------
// @file      : fingerprint.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/28 20:40
// -------------------------------------------

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// FingerprintRule 对应 FingerprintRules 集合
type FingerprintRule struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Rule           string             `json:"rule" bson:"rule"`
	Express        []string           `json:"express" bson:"express"`
	Category       string             `json:"category" bson:"category"`
	ParentCategory string             `json:"parent_category" bson:"parent_category"`
	Amount         int                `json:"amount" bson:"amount"`
	State          bool               `json:"state" bson:"state"`
	Company        string             `json:"company" bson:"company"`
}

// FingerprintQuery 查询参数
type FingerprintQuery struct {
	Search    string `json:"search"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}
