package ip

import (
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	iprepo "github.com/Autumn-27/ScopeSentry-go/internal/repositories/assets/ip"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义 IP 资产服务接口
type Service interface {
	GetIPAssets(ctx *gin.Context, query models.SearchRequest) (*models.IPAssetResponse, error)
}

type service struct {
	repo iprepo.Repository
}

// NewService 创建 IP 资产服务实例
func NewService() Service {
	return &service{
		repo: iprepo.NewRepository(),
	}
}

// GetIPAssets 获取 IP 资产列表
func (s *service) GetIPAssets(ctx *gin.Context, query models.SearchRequest) (*models.IPAssetResponse, error) {
	if query.PageIndex <= 0 {
		query.PageIndex = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	query.Index = "IPAsset"
	searchQuery, err := helper.GetSearchQuery(query)
	if err != nil {
		return nil, err
	}
	filter := bson.M(searchQuery)

	reqCtx := ctx.Request.Context()
	total, err := s.repo.Count(reqCtx, filter)
	if err != nil {
		return nil, err
	}

	projection := bson.M{
		"ip":    1,
		"ports": 1,
		"time":  1,
	}

	opts := options.Find().
		SetProjection(projection).
		SetSkip(int64((query.PageIndex - 1) * query.PageSize)).
		SetLimit(int64(query.PageSize)).
		SetSort(bson.D{{Key: "time", Value: -1}})

	assets, err := s.repo.Find(reqCtx, filter, opts)
	if err != nil {
		return nil, err
	}

	rows := flattenIPAssets(assets)

	return &models.IPAssetResponse{
		List:  rows,
		Total: total,
	}, nil
}

func flattenIPAssets(assets []models.IPAsset) []models.IPAssetFlat {
	var rows []models.IPAssetFlat

	for _, asset := range assets {
		totalRows := 0
		for _, port := range asset.Ports {
			serverCount := len(port.Server)
			if serverCount == 0 {
				serverCount = 1
			}
			totalRows += serverCount
		}
		if totalRows == 0 {
			totalRows = 1
		}

		if len(asset.Ports) == 0 {
			rows = append(rows, models.IPAssetFlat{
				IP:          asset.IP,
				IPRowSpan:   totalRows,
				PortRowSpan: 1,
				Time:        asset.Time,
			})
			continue
		}

		ipRowAssigned := false
		for _, port := range asset.Ports {
			portRowSpan := len(port.Server)
			if portRowSpan == 0 {
				portRowSpan = 1
			}
			portRowAssigned := false

			servers := port.Server
			if len(servers) == 0 {
				servers = []models.PortServer{{}}
			}

			for _, server := range servers {
				row := models.IPAssetFlat{
					IP:        asset.IP,
					Port:      port.Port,
					Domain:    server.Domain,
					Service:   server.Service,
					WebServer: server.WebServer,
					Products:  append([]string{}, server.Technologies...),
					Time:      asset.Time,
				}

				if !ipRowAssigned {
					row.IPRowSpan = totalRows
					ipRowAssigned = true
				} else {
					row.IPRowSpan = 0
				}

				if !portRowAssigned {
					row.PortRowSpan = portRowSpan
					portRowAssigned = true
				} else {
					row.PortRowSpan = 0
				}

				rows = append(rows, row)
			}
		}
	}

	return rows
}
