// dirscan-------------------------------------
// @file      : dirscan.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/17 15:12
// -------------------------------------------

package dirscan

import (
	"context"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/dirscan"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	List(ctx context.Context, query models.SearchRequest) ([]models.DirScanResult, error)
}

type service struct {
	repo dirscan.Repository
}

func NewService() Service {
	return &service{repo: dirscan.NewRepository()}
}

func (s *service) List(ctx context.Context, query models.SearchRequest) ([]models.DirScanResult, error) {
	query.Index = "DirScanResult"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)
	sortBy := bson.D{{Key: "_id", Value: -1}}
	if val, ok := query.Sort["length"]; ok {
		dir := -1
		if val == "ascending" {
			dir = 1
		}
		sortBy = bson.D{{Key: "length", Value: dir}}
	}

	opts := options.Find().
		SetProjection(bson.M{
			"_id":    1,
			"url":    1,
			"status": 1,
			"msg":    1,
			"length": 1,
			"tags":   1,
		}).
		SetSort(sortBy).
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize))

	return s.repo.FindWithPagination(ctx, filter, opts)
}
