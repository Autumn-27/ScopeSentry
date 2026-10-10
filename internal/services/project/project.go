package project

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"sort"
	"strings"
	"sync"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"golang.org/x/sync/errgroup"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	repo "github.com/Autumn-27/ScopeSentry/internal/repositories/project"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetProjectsByTag(ctx *gin.Context) ([]models.TagGroup, error)
	GetTarget(ctx *gin.Context, id string) (string, error)
	GetTargets(ctx context.Context, ids []string) (string, error)

	GetProjectsData(ctx *gin.Context, search string, pageIndex, pageSize int) (models.ProjectListResponse, error)
	GetProjectContent(ctx *gin.Context, id string) (*models.ProjectContentResponse, error)
	AddProject(ctx *gin.Context, p *models.Project) error
	DeleteProjects(ctx *gin.Context, ids []string, delA bool) error
	UpdateProject(ctx *gin.Context, p *models.UpdateProject) error
	UpdateAssetsProject(ctx *gin.Context, id string) error
	RefreshProjectAssociations(ctx context.Context) error
}

type service struct {
	projectRepo repo.Repository
}

func NewService() Service {
	return &service{
		projectRepo: repo.NewRepository(),
	}
}

func (s *service) GetProjectsByTag(ctx *gin.Context) ([]models.TagGroup, error) {
	// 在service层定义需要的字段
	projection := bson.M{
		"_id":  1,
		"name": 1,
		"tag":  1,
	}

	projects, err := s.projectRepo.GetProjectsByTag(ctx, projection)
	if err != nil {
		return nil, err
	}

	// 按tag分组
	tagMap := make(map[string][]models.TagProject)
	for _, p := range projects {
		tagProject := models.TagProject{
			Value: p.ID.Hex(),
			Label: p.Name,
		}
		tagMap[p.Tag] = append(tagMap[p.Tag], tagProject)
	}

	// 转换为响应格式并按tag排序
	var result []models.TagGroup
	for tag, projects := range tagMap {
		tagGroup := models.TagGroup{
			Label:    tag,
			Value:    "",
			Children: projects,
		}
		result = append(result, tagGroup)
	}

	// 按tag排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})

	return result, nil
}

func (s *service) GetTarget(ctx *gin.Context, id string) (string, error) {
	return s.projectRepo.GetTarget(ctx, id)
}

func (s *service) GetTargets(ctx context.Context, ids []string) (string, error) {
	return s.projectRepo.GetTargets(ctx, ids)
}

func (s *service) GetProjectsData(ctx *gin.Context, search string, pageIndex, pageSize int) (models.ProjectListResponse, error) {
	// tag 聚合统计
	tagCounts, err := s.projectRepo.AggregateTagCounts(ctx, search)
	if err != nil {
		return models.ProjectListResponse{}, err
	}
	// 计算 All 总数
	all := 0
	for _, v := range tagCounts {
		all += v
	}
	tagCounts["All"] = all

	result := make(map[string][]models.ProjectBrief)
	var mu sync.Mutex
	base := ctx.Request.Context()
	g, _ := errgroup.WithContext(base)
	for tag := range tagCounts {
		t := tag
		g.Go(func() error {
			items, err := s.projectRepo.ListProjectsByTag(ctx, search, t, pageIndex, pageSize)
			if err != nil {
				return err
			}
			mu.Lock()
			result[t] = items
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return models.ProjectListResponse{}, err
	}
	return models.ProjectListResponse{Result: result, Tag: tagCounts}, nil
}

func (s *service) GetProjectContent(ctx *gin.Context, id string) (*models.ProjectContentResponse, error) {
	p, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	target, err := s.projectRepo.GetTarget(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := &models.ProjectContentResponse{
		Name:           p.Name,
		Tag:            p.Tag,
		Target:         target,
		Node:           p.Node,
		Logo:           p.Logo,
		ScheduledTasks: p.ScheduledTasks,
		Hour:           p.Hour,
		AllNode:        p.AllNode,
		Duplicates:     p.Duplicates,
		Template:       p.Template,
		Ignore:         p.Ignore,
	}
	return resp, nil
}

func (s *service) AddProject(ctx *gin.Context, p *models.Project) error {
	// 名称唯一
	exists, err := s.projectRepo.ExistsByName(ctx, p.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("name already exists")
	}
	p.Tp = "project"
	p.Target = strings.TrimSpace(p.Target)
	targetList, err := helper.GetTargetList(p.Target, p.Ignore)
	if err != nil {
		return err
	}
	rootDomains := []string{}
	for _, tg := range targetList {
		var rootDomain string
		if strings.Contains(tg, "CMP:") || strings.Contains(tg, "ICP:") ||
			strings.Contains(tg, "APP:") || strings.Contains(tg, "APP-ID:") {

			if strings.Contains(tg, "ICP:") {
				rootDomain = getBeforeLastDash(strings.ReplaceAll(tg, "ICP:", ""))
				rootDomain = "ICP:" + rootDomain
			} else {
				rootDomain = tg
			}
		} else {
			rootDomain, _ = helper.GetRootDomain(tg)
		}

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}
	p.RootDomains = rootDomains
	rawTarget := p.Target
	p.Target = ""
	projectId, err := s.projectRepo.InsertProject(ctx, p)
	if err != nil {
		return err
	}
	err = s.projectRepo.UpsertProjectTarget(ctx, projectId, rawTarget)
	if err != nil {
		return err
	}
	go func() {
		err = s.projectRepo.UpdateAssetsProject(context.Background(), rootDomains, projectId, false)
		if err != nil {
			logger.Error(err.Error())
		}
		if err := s.recalcAssetCount(context.Background(), projectId); err != nil {
			logger.Error(err.Error())
		}
	}()
	//id, err := s.projectRepo.InsertProject(ctx, p)
	//if err != nil {
	//	return err
	//}
	//if err := s.projectRepo.UpsertProjectTarget(ctx, id, p.Target); err != nil {
	//	return err
	//}
	// 调度
	//if p.ScheduledTasks {
	//	if err := s.projectRepo.CreateOrUpdateProjectSchedule(ctx, id, p.Hour, true, p.Name); err != nil {
	//		return err
	//	}
	//}
	// runNow 可后续接入任务系统
	return nil
}

func (s *service) UpdateProject(ctx *gin.Context, p *models.UpdateProject) error {
	p.Tp = "project"
	p.Target = strings.TrimSpace(p.Target)
	targetList, err := helper.GetTargetList(p.Target, p.Ignore)
	if err != nil {
		return err
	}
	rootDomains := []string{}
	for _, tg := range targetList {
		var rootDomain string
		if strings.Contains(tg, "CMP:") || strings.Contains(tg, "ICP:") ||
			strings.Contains(tg, "APP:") || strings.Contains(tg, "APP-ID:") {

			if strings.Contains(tg, "ICP:") {
				rootDomain = getBeforeLastDash(strings.ReplaceAll(tg, "ICP:", ""))
				rootDomain = "ICP:" + rootDomain
			} else {
				rootDomain = tg
			}
		} else {
			rootDomain, _ = helper.GetRootDomain(tg)
		}

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}
	p.RootDomains = rootDomains
	rawTarget := p.Target
	p.Target = ""
	if err := s.projectRepo.UpdateProject(ctx, p); err != nil {
		return err
	}

	if err = s.projectRepo.UpsertProjectTarget(ctx, p.ID, rawTarget); err != nil {
		return err
	}
	go func() {
		err = s.projectRepo.UpdateAssetsProject(context.Background(), rootDomains, p.ID, true)
		if err != nil {
			logger.Error(err.Error())
		}
		if err := s.recalcAssetCount(context.Background(), p.ID); err != nil {
			logger.Error(err.Error())
		}
	}()
	return nil
}

func (s *service) DeleteProjects(ctx *gin.Context, ids []string, delA bool) error {

	// 删除项目
	if err := s.projectRepo.DeleteProjects(ctx, ids); err != nil {
		return err
	}
	if err := s.projectRepo.DeleteProjectTargets(ctx, ids); err != nil {
		return err
	}
	if delA {
		go func() {
			s.projectRepo.DeleteProjectAsset(ctx, ids)
		}()

	}
	return nil
}

func getBeforeLastDash(s string) string {
	index := strings.LastIndex(s, "-") // 查找最后一个 '-'
	if index != -1 {
		return s[:index]
	}
	return s
}

func (s *service) UpdateAssetsProject(ctx *gin.Context, id string) error {
	// 获取项目信息
	project, err := s.projectRepo.FindByID(ctx.Request.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	// 获取项目目标
	target, err := s.projectRepo.GetTarget(ctx.Request.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to get project target: %w", err)
	}

	// 解析目标获取根域名列表
	targetList, err := helper.GetTargetList(target, project.Ignore)
	if err != nil {
		return fmt.Errorf("failed to parse target list: %w", err)
	}

	rootDomains := []string{}
	for _, tg := range targetList {
		var rootDomain string
		if strings.Contains(tg, "CMP:") || strings.Contains(tg, "ICP:") ||
			strings.Contains(tg, "APP:") || strings.Contains(tg, "APP-ID:") {

			if strings.Contains(tg, "ICP:") {
				rootDomain = getBeforeLastDash(strings.ReplaceAll(tg, "ICP:", ""))
				rootDomain = "ICP:" + rootDomain
			} else {
				rootDomain = tg
			}
		} else {
			rootDomain, _ = helper.GetRootDomain(tg)
		}

		// 检查是否已经存在于 rootDomains
		exists := false
		for _, rd := range rootDomains {
			if rd == rootDomain {
				exists = true
				break
			}
		}
		if !exists {
			rootDomains = append(rootDomains, rootDomain)
		}
	}

	// 更新资产项目关联
	err = s.projectRepo.UpdateAssetsProject(ctx.Request.Context(), rootDomains, id, true)
	if err != nil {
		return fmt.Errorf("failed to update assets project: %w", err)
	}

	return s.recalcAssetCount(ctx.Request.Context(), id)
}

// recalcAssetCount 重新统计项目的资产总数并写回 AssetCount 字段。
func (s *service) recalcAssetCount(ctx context.Context, projectID string) error {
	count, err := s.projectRepo.CountProjectAssets(ctx, projectID)
	if err != nil {
		return err
	}
	return s.projectRepo.SetAssetCount(ctx, projectID, count)
}

// RefreshProjectAssociations 重新盖章所有项目的资产关联并回填 AssetCount。
//
// 背景(见 issue #338):扫描结果入库不携带 project 字段,资产关联只在项目
// 增/改时通过 UpdateAssetsProject 回填,任务结束后无人触发——导致 MCP/API
// 按项目过滤(list_assets filter[project])拿不到数据、AssetCount 恒 0。
// 本方法在任务结束事件处调用作为自愈:对每个项目执行 add/remove 两个方向
// 的盖章,再统计回填。幂等,可重复调用。
func (s *service) RefreshProjectAssociations(ctx context.Context) error {
	projects, err := s.projectRepo.GetProjectsByTag(ctx, bson.M{"_id": 1, "root_domains": 1, "AssetCount": 1})
	if err != nil {
		return err
	}
	for _, p := range projects {
		if len(p.RootDomains) == 0 {
			continue
		}
		id := p.ID.Hex()
		if err := s.projectRepo.UpdateAssetsProject(ctx, p.RootDomains, id, true); err != nil {
			logger.Error(fmt.Sprintf("refresh project %s association failed: %v", id, err))
			continue
		}
		if err := s.recalcAssetCount(ctx, id); err != nil {
			logger.Error(fmt.Sprintf("refresh project %s recalc count failed: %v", id, err))
		}
	}
	return nil
}
