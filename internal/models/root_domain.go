package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RootDomain 根域名数据结构
type RootDomain struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain  string             `bson:"domain" json:"domain"`
	ICP     string             `bson:"icp" json:"icp"`
	Company string             `bson:"company" json:"company"`
	Project string             `bson:"project" json:"project"`
	Time    string             `bson:"time" json:"time"`
	Tags    []string           `bson:"tags" json:"tags"`
}

// RootDomainResponse 根域名响应数据结构
type RootDomainResponse struct {
	List  []RootDomain `json:"list"`
	Total int64        `json:"total"`
}
