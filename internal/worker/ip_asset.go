// worker-------------------------------------
// @file      : ip.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/11/12 22:29
// -------------------------------------------

package worker

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"

	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ipAssetBatchLimit    int64 = 500
	ipAssetIdleSleep           = 60 * time.Second
	ipAssetErrorSleep          = 10 * time.Second
	ipAssetBatchInterval       = 3 * time.Second
)

type ipAssetTmpRecord struct {
	ID                primitive.ObjectID `bson:"_id"`
	models.IPAssetTmp `bson:",inline"`
}

type ipAggregation struct {
	ip          string
	asset       models.IPAsset
	taskNames   map[string]struct{}
	rootDomains map[string]struct{}
	projects    map[string]struct{}
	portMap     map[string]*models.IPPorts
}

func newIPAggregation(ip string) *ipAggregation {
	return &ipAggregation{
		ip:          ip,
		asset:       models.IPAsset{IP: ip},
		taskNames:   make(map[string]struct{}),
		rootDomains: make(map[string]struct{}),
		projects:    make(map[string]struct{}),
		portMap:     make(map[string]*models.IPPorts),
	}
}

func (agg *ipAggregation) append(tmp models.IPAssetTmp) {
	if tmp.Project != "" {
		agg.projects[tmp.Project] = struct{}{}
	}

	for _, task := range tmp.TaskName {
		if task == "" {
			continue
		}
		agg.taskNames[task] = struct{}{}
	}

	if tmp.RootDomain != "" {
		agg.rootDomains[tmp.RootDomain] = struct{}{}
	}

	if tmp.Port == "" {
		return
	}

	portInfo, ok := agg.portMap[tmp.Port]
	if !ok {
		portInfo = &models.IPPorts{Port: tmp.Port}
		agg.portMap[tmp.Port] = portInfo
	}

	portInfo.Server = append(portInfo.Server, models.PortServer{
		Domain:       tmp.Domain,
		Service:      tmp.Service,
		WebServer:    tmp.WebServer,
		Technologies: append([]string{}, tmp.Technologies...),
	})
}

func (agg *ipAggregation) finalize() models.IPAsset {
	asset := agg.asset

	if len(agg.projects) > 0 {
		asset.Project = toSortedSlice(agg.projects)
	} else {
		asset.Project = []string{}
	}

	if len(agg.taskNames) > 0 {
		asset.TaskName = toSortedSlice(agg.taskNames)
	} else {
		asset.TaskName = []string{}
	}

	if len(agg.rootDomains) > 0 {
		asset.RootDomain = toSortedSlice(agg.rootDomains)
	} else {
		asset.RootDomain = []string{}
	}

	if len(agg.portMap) > 0 {
		ports := make([]models.IPPorts, 0, len(agg.portMap))
		for _, p := range agg.portMap {
			serverCopy := make([]models.PortServer, len(p.Server))
			copy(serverCopy, p.Server)
			ports = append(ports, models.IPPorts{
				Port:   p.Port,
				Server: serverCopy,
			})
		}
		sort.Slice(ports, func(i, j int) bool {
			return ports[i].Port < ports[j].Port
		})
		asset.Ports = ports
	} else {
		asset.Ports = []models.IPPorts{}
	}

	return asset
}

func toSortedSlice(set map[string]struct{}) []string {
	vals := make([]string, 0, len(set))
	for v := range set {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	return vals
}

func aggregateIPAssets(records []ipAssetTmpRecord) map[string]models.IPAsset {
	result := make(map[string]models.IPAsset)
	aggMap := make(map[string]*ipAggregation)

	for _, record := range records {
		if record.IP == "" {
			continue
		}
		if _, ok := aggMap[record.IP]; !ok {
			aggMap[record.IP] = newIPAggregation(record.IP)
		}
		aggMap[record.IP].append(record.IPAssetTmp)
	}

	for ip, agg := range aggMap {
		result[ip] = agg.finalize()
	}

	return result
}

func mergeIPAsset(existing, incoming models.IPAsset) models.IPAsset {
	merged := models.IPAsset{
		IP: incoming.IP,
	}

	merged.Project = mergeStringSlices(existing.Project, incoming.Project)
	merged.TaskName = mergeStringSlices(existing.TaskName, incoming.TaskName)
	merged.RootDomain = mergeStringSlices(existing.RootDomain, incoming.RootDomain)
	merged.Ports = mergeIPPorts(existing.Ports, incoming.Ports)

	return merged
}

func mergeStringSlices(a, b []string) []string {
	set := make(map[string]struct{})

	for _, v := range a {
		if v == "" {
			continue
		}
		set[v] = struct{}{}
	}

	for _, v := range b {
		if v == "" {
			continue
		}
		set[v] = struct{}{}
	}

	if len(set) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(set))
	for v := range set {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

func mergeIPPorts(existing, incoming []models.IPPorts) []models.IPPorts {
	portMap := make(map[string][]models.PortServer)

	addPorts := func(ports []models.IPPorts) {
		for _, p := range ports {
			if p.Port == "" {
				continue
			}
			portMap[p.Port] = mergeServers(portMap[p.Port], p.Server)
		}
	}

	addPorts(existing)
	addPorts(incoming)

	if len(portMap) == 0 {
		return []models.IPPorts{}
	}

	ports := make([]models.IPPorts, 0, len(portMap))
	for port, servers := range portMap {
		serversCopy := make([]models.PortServer, len(servers))
		copy(serversCopy, servers)
		ports = append(ports, models.IPPorts{
			Port:   port,
			Server: serversCopy,
		})
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})

	return ports
}

func mergeServers(existing, incoming []models.PortServer) []models.PortServer {
	type serverAggregation struct {
		domain       string
		service      string
		webServer    string
		technologies map[string]struct{}
	}

	aggregations := make(map[string]*serverAggregation)

	addServer := func(server models.PortServer) {
		key := fmt.Sprintf("%s|%s", server.Domain, server.Service)
		agg, ok := aggregations[key]
		if !ok {
			agg = &serverAggregation{
				domain:       server.Domain,
				service:      server.Service,
				technologies: make(map[string]struct{}),
			}
			aggregations[key] = agg
		}

		if agg.domain == "" && server.Domain != "" {
			agg.domain = server.Domain
		}
		if agg.service == "" && server.Service != "" {
			agg.service = server.Service
		}

		if server.WebServer != "" {
			agg.webServer = server.WebServer
		}

		for _, tech := range server.Technologies {
			if tech == "" {
				continue
			}
			agg.technologies[tech] = struct{}{}
		}
	}

	for _, server := range existing {
		addServer(server)
	}
	for _, server := range incoming {
		addServer(server)
	}

	result := make([]models.PortServer, 0, len(aggregations))
	for _, agg := range aggregations {
		technologies := make([]string, 0, len(agg.technologies))
		for tech := range agg.technologies {
			technologies = append(technologies, tech)
		}
		sort.Strings(technologies)

		result = append(result, models.PortServer{
			Domain:       agg.domain,
			Service:      agg.service,
			WebServer:    agg.webServer,
			Technologies: technologies,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Domain == result[j].Domain {
			return result[i].Service < result[j].Service
		}
		return result[i].Domain < result[j].Domain
	})

	return result
}

func IPAssetHandle(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("任务发生错误，已恢复:", r)
			select {
			case <-ctx.Done():
				fmt.Println("任务已取消，IPAssetHandle不再重启")
			default:
				go IPAssetHandle(ctx)
			}
		}
	}()

	repo := common.NewRepository()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("任务被取消，IPAssetHandle停止执行")
			return
		default:
			opts := options.Find().
				SetLimit(ipAssetBatchLimit).
				SetSort(bson.D{{Key: "_id", Value: 1}})

			var tmpRecords []ipAssetTmpRecord
			err := repo.Find(ctx, "IPAssetTmp", bson.M{}, opts, &tmpRecords)
			if err != nil {
				logger.Error(fmt.Sprintf("查询 IPAssetTmp 数据失败: %v\n", err))
				time.Sleep(ipAssetErrorSleep)
				continue
			}

			if len(tmpRecords) == 0 {
				time.Sleep(ipAssetIdleSleep)
				continue
			}

			ipAssets := aggregateIPAssets(tmpRecords)
			for ip, asset := range ipAssets {
				var existing models.IPAsset
				findErr := repo.Find(ctx, "IPAsset", bson.M{"ip": ip}, options.Find().SetLimit(1), &existing)
				var merged models.IPAsset
				if findErr != nil {
					if errors.Is(findErr, mongo.ErrNoDocuments) {
						merged = asset
					} else {
						logger.Error(fmt.Sprintf("查询 IPAsset 失败 ip=%s: %v\n", ip, findErr))
						continue
					}
				} else {
					merged = mergeIPAsset(existing, asset)
				}
				merged.Time = helper.GetNowTimeString()
				if err := repo.Upsert(ctx, "IPAsset", bson.M{"ip": ip}, merged); err != nil {
					logger.Error(fmt.Sprintf("写入 IPAsset 失败 ip=%s: %v\n", ip, err))
				}
			}

			var ids []primitive.ObjectID
			for _, record := range tmpRecords {
				if record.ID != primitive.NilObjectID {
					ids = append(ids, record.ID)
				}
			}

			if len(ids) > 0 {
				_, err := repo.DeleteMany(ctx, "IPAssetTmp", bson.M{"_id": bson.M{"$in": ids}})
				if err != nil {
					logger.Error(fmt.Sprintf("删除 IPAssetTmp 数据失败: %v\n", err))
				}
			}

			logger.Info(fmt.Sprintf("IPAssetHandle 完成一次处理，输入 %d 条，输出 %d 个 IP 资产", len(tmpRecords), len(ipAssets)))

			time.Sleep(ipAssetBatchInterval)
		}
	}
}
