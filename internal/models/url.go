package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL 表示一个URL记录
type URL struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Input  string             `bson:"input" json:"input"`
	Source string             `bson:"source" json:"source"`
	Status int                `bson:"status" json:"status"`
	Length int                `bson:"length" json:"length"`
	Type   string             `bson:"outputtype" json:"type"`
	Output string             `bson:"output" json:"url"`
	Time   string             `bson:"time" json:"time"`
	Tags   []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
