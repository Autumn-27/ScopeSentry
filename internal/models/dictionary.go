// models-------------------------------------
// @file      : dictionary.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/28 20:29
// -------------------------------------------

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Port 仅用于基础端口键值
type Port struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// DictionaryMeta 字典元数据
type DictionaryMeta struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Category string             `json:"category" bson:"category"`
	Size     string             `json:"size" bson:"size"`
}

// PortDoc 带有ID的端口字典文档
type PortDoc struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Value string             `json:"value" bson:"value"`
}

type UpdatePort struct {
	ID    string `json:"id"`
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}
