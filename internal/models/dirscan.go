// models-------------------------------------
// @file      : dirscan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/17 15:07
// -------------------------------------------

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DirScanResult struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL    string             `json:"url"`
	Status int                `json:"status"`
	Msg    string             `json:"msg"`
	Length int                `json:"length"`
	Tags   []string           `json:"tags"`
}
