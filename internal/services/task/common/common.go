// common-------------------------------------
// @file      : common.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/6/9 21:03
// -------------------------------------------

package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/task/task"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/asset"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/root_domain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/subdomain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/url"
	"github.com/Autumn-27/ScopeSentry/internal/services/dictionary"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
)

// Service 定义任务服务接口
type Service interface {
	Insert(ctx *gin.Context, task *models.Task) (string, error)
	CreateTaskScan(ctx context.Context, task models.Task, id string, stopToStart bool) (int64, error)
	GetScanTemplate(ctx context.Context, task *models.Task) (*models.ScanTemplate, error)
	ParameterParser(ctx context.Context, parameters models.Parameters) (models.Parameters, error)
}

type service struct {
	taskRepo          task.Repository
	projectService    project.Service
	assetService      asset.Service
	rootDomainService root_domain.Service
	subdomainService  subdomain.Service
	urlService        url.Service
	nodeService       node.Service
	dictionaryService dictionary.ManageService
	portsService      dictionary.PortService
}

// NewService 创建任务服务实例
func NewService() Service {
	return &service{
		taskRepo:          task.NewRepository(),
		projectService:    project.NewService(),
		assetService:      asset.NewService(),
		rootDomainService: root_domain.NewService(),
		subdomainService:  subdomain.NewService(),
		nodeService:       node.NewService(),
		dictionaryService: dictionary.NewManageService(),
		portsService:      dictionary.NewPortService(),
	}
}

func (s *service) Insert(ctx *gin.Context, task *models.Task) (string, error) {
	// task.TargetSource 分为普通、从资产选择数据进行创建（支持资产、根域名、子域名、url 这里修改原始逻辑 统一为后边的情况）、从项目、从资产、从根域名、从子域名、从url
	var query models.SearchRequest
	var target string
	var err error
	query.PageSize = task.TargetNumber
	query.SearchExpression = task.Search
	query.Filter = task.Filter
	if task.TargetSource == "general" {
		//普通的情况不需要做改变
	} else if task.TargetSource == "project" {
		// 从项目中获取数据，读取projet的目标
		target, err = s.projectService.GetTargets(ctx, task.Project)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "asset" || (task.TargetSource == "assetSource" && task.TargetTp == "search") {
		// 从资产处创建
		target, err = s.assetService.GetTaskTarget(ctx, query)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "RootDomain" || (task.TargetSource == "RootDomainSource" && task.TargetTp == "search") {
		// 从根域名创建
		target, err = s.rootDomainService.GetTaskTarget(ctx, query)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "subdomain" || (task.TargetSource == "subdomainSource" && task.TargetTp == "search") {
		// 从子域名处创建
		target, err = s.subdomainService.GetTaskTarget(ctx, query)
		if err != nil {
		}
	} else if task.TargetSource == "UrlScan" || (task.TargetSource == "UrlScanSource" && task.TargetTp == "search") {
		// 从url扫描处创建
		target, err = s.urlService.GetTaskTarget(ctx, query)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "assetSource" && task.TargetTp == "select" {
		target, err = s.assetService.GetTaskTargetByIDs(ctx, task.TargetIds)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "RootDomainSource" && task.TargetTp == "select" {
		target, err = s.rootDomainService.GetTaskTargetByIDs(ctx, task.TargetIds)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "subdomainSource" && task.TargetTp == "select" {
		target, err = s.subdomainService.GetTaskTargetByIDs(ctx, task.TargetIds)
		if err != nil {
			return "", err
		}
	} else if task.TargetSource == "UrlScanSource" && task.TargetTp == "select" {
		target, err = s.urlService.GetTaskTargetByIDs(ctx, task.TargetIds)
		if err != nil {
			return "", err
		}
	}
	if task.Target == "" {
		task.Target = target
	}
	var result []string
	result, _ = helper.GetTargetList(task.Target, task.Ignore)
	task.Target = strings.Join(result, "\n")
	task.TaskNum = len(result)
	if task.TaskNum == 0 {
		return "", errors.New("Target is null")
	}
	task.Progress = 0
	task.CreatTime = helper.GetNowTimeString()
	task.EndTime = ""
	task.Status = 1
	task.Type = task.TargetSource
	insertId, err := s.taskRepo.Insert(ctx, task)
	if err != nil {
		return "", err
	}
	go func() {
		_, err = s.CreateTaskScan(context.Background(), *task, insertId, false)
		if err != nil {
			logger.Error(fmt.Sprintf("CreateTaskScan %v", err))
		}
	}()
	hex, err := primitive.ObjectIDFromHex(insertId)
	if err != nil {
		return "", err
	}
	return hex.Hex(), nil
}

func (s *service) CreateTaskScan(ctx context.Context, task models.Task, id string, stopToStart bool) (int64, error) {
	logger.Info(fmt.Sprintf("CreateTaskScan task:%v %v", id, task.Name))
	task.ID, _ = primitive.ObjectIDFromHex(id)
	task.IsStart = stopToStart
	if task.AllNode {
		// 选择自动加入后 获取所有节点
		allNode, err := s.nodeService.GetNodeData(ctx, false)
		if err != nil {
			return 0, err
		}
		var nameMap = make(map[string]struct{})
		for _, n := range task.Node {
			nameMap[n] = struct{}{}
		}
		for _, n := range allNode {
			if _, exists := nameMap[n.Name]; !exists {
				nameMap[n.Name] = struct{}{}
				task.Node = append(task.Node, n.Name)
			}
		}
	}

	// 如果是暂停之后再开始的，则不需要删除缓存和填入目标
	if !stopToStart {
		// 删除缓存
		err := s.taskRepo.ClearTaskCache(ctx, id)
		if err != nil {
			return 0, err
		}
		// 生成原始target
		targets, err := helper.GetTargetList(task.Target, task.Ignore)
		if err != nil {
			return 0, err
		}
		err = s.taskRepo.PushTaskInfoList(ctx, id, targets)
		if err != nil {
			return 0, err
		}
	}
	template, err := s.GetScanTemplate(ctx, &task)
	if err != nil {
		return 0, err
	}
	for _, n := range task.Node {
		template.Target = ""
		err := s.taskRepo.RPushNodeTask(ctx, n, template)
		if err != nil {
			logger.Error(fmt.Sprintf("CreateTaskScan task:%v error:%v", n, err))
			continue
		}
	}
	logger.Info(fmt.Sprintf("CreateTaskScan task:%v success", task.ID))
	return 0, nil
}

func (s *service) GetScanTemplate(ctx context.Context, task *models.Task) (*models.ScanTemplate, error) {
	scanTemplate, err := s.taskRepo.FindTemplateByID(ctx, task.Template)
	if err != nil {
		return &models.ScanTemplate{}, err
	}
	if len(scanTemplate.VulList) != 0 {
		if helper.StringInSlice("ed93b8af6b72fe54a60efdb932cf6fbc", scanTemplate.VulnerabilityScan) {
			vulTmp := ""
			if helper.StringInSlice("All Poc", scanTemplate.VulList) {
				vulTmp = "*"
			} else {
				for _, vul := range scanTemplate.VulList {
					vulTmp += vul + ".yaml" + ","
				}
				vulTmp = strings.TrimRight(vulTmp, ",")
			}
			if _, exists := scanTemplate.Parameters.VulnerabilityScan["ed93b8af6b72fe54a60efdb932cf6fbc"]; exists {
				scanTemplate.Parameters.VulnerabilityScan["ed93b8af6b72fe54a60efdb932cf6fbc"] = scanTemplate.Parameters.VulnerabilityScan["ed93b8af6b72fe54a60efdb932cf6fbc"] + " -t " + vulTmp
			} else {
				scanTemplate.Parameters.VulnerabilityScan["ed93b8af6b72fe54a60efdb932cf6fbc"] = "-t " + vulTmp
			}
		}
	}
	scanTemplate.Parameters, err = s.ParameterParser(ctx, scanTemplate.Parameters)
	if err != nil {
		return &models.ScanTemplate{}, err
	}
	scanTemplate.VulList = []string{}
	scanTemplate.ID = task.ID.Hex()
	scanTemplate.Ignore = task.Ignore
	scanTemplate.Type = task.Type
	scanTemplate.Duplicates = task.Duplicates
	scanTemplate.IsStart = task.IsStart
	scanTemplate.TaskName = task.Name
	return scanTemplate, err
}

var placeholderRegex = regexp.MustCompile(`\{(.*?)\}`)

// ParameterParser 解析参数中的字典和端口引用
func (s *service) ParameterParser(ctx context.Context, parameters models.Parameters) (models.Parameters, error) {
	// 1. 获取字典数据
	dictList := make(map[string]string)
	dicts, err := s.dictionaryService.List(ctx)
	if err != nil {
		return models.Parameters{}, err
	}
	for _, dict := range dicts {
		dictList[strings.ToLower(dict.Category+"."+dict.Name)] = dict.ID.Hex()
	}

	// 2. 获取端口数据
	portList := make(map[string]string)
	ports, _, err := s.portsService.Get(ctx, "", 1, 1000)
	if err != nil {
		return models.Parameters{}, err
	}
	for _, port := range ports {
		portList[strings.ToLower(port.Name)] = port.Value
	}

	// 公共处理函数
	processMap := func(m map[string]string) {
		for k, v := range m {
			newValue := placeholderRegex.ReplaceAllStringFunc(v, func(match string) string {
				// match 就是 {dict.xxx} 或 {port.xxx}
				inner := strings.Trim(match, "{}")
				parts := strings.SplitN(inner, ".", 2)
				if len(parts) != 2 {
					return match
				}

				tp := strings.ToLower(parts[0])
				val := strings.ToLower(parts[1])

				switch tp {
				case "dict":
					if real, ok := dictList[val]; ok {
						return real
					}
					logger.Error(fmt.Sprintf("未找到字典参数: %s", val))
				case "port":
					if real, ok := portList[val]; ok {
						return real
					}
					logger.Error(fmt.Sprintf("未找到端口参数: %s", val))
				}
				return match
			})
			m[k] = newValue
		}
	}

	// 3. 自动批量处理所有字段（集中在一个 slice 中，方便以后加字段）
	allMaps := []*map[string]string{
		&parameters.TargetHandler,
		&parameters.SubdomainScan,
		&parameters.SubdomainSecurity,
		&parameters.PortScanPreparation,
		&parameters.PortScan,
		&parameters.PortFingerprint,
		&parameters.AssetMapping,
		&parameters.AssetHandle,
		&parameters.URLScan,
		&parameters.WebCrawler,
		&parameters.URLSecurity,
		&parameters.DirScan,
		&parameters.VulnerabilityScan,
	}

	for _, mp := range allMaps {
		processMap(*mp)
	}

	return parameters, nil
}
