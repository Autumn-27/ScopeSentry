package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	assetCommon "github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/app"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/asset"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/crawler"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/dirscan"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/ip"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/mp"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/page_monitoring"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/root_domain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/sensitive"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/subdomain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/url"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/vulnerability"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/services/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	taskCommon "github.com/Autumn-27/ScopeSentry/internal/services/task/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/task"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/template"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type deps struct {
	projectService    project.Service
	taskService       task.Service
	taskCommonService taskCommon.Service
	templateService   template.Service
	assetService      asset.Service
	rootDomainService root_domain.Service
	subdomainService  subdomain.Service
	appService        app.Service
	mpService         mp.Service
	urlService        url.Service
	sensitiveService  sensitive.Service
	dirscanService    dirscan.Service
	crawlerService    crawler.Service
	vulnService       vulnerability.Service
	ipService         ip.Service
	pageMonService    page_monitoring.Service
	commonService     assetCommon.Service
	nodeService       node.Service
	pluginService     plugin.Service
}

var d = &deps{
	projectService:    project.NewService(),
	taskService:       task.NewService(),
	taskCommonService: taskCommon.NewService(),
	templateService:   template.NewService(),
	assetService:      asset.NewService(),
	rootDomainService: root_domain.NewService(),
	subdomainService:  subdomain.NewService(),
	appService:        app.NewService(),
	mpService:         mp.NewService(),
	urlService:        url.NewService(),
	sensitiveService:  sensitive.NewService(),
	dirscanService:    dirscan.NewService(),
	crawlerService:    crawler.NewService(),
	vulnService:       vulnerability.NewService(),
	ipService:         ip.NewService(),
	pageMonService:    page_monitoring.NewService(),
	commonService:     assetCommon.NewService(),
	nodeService:       node.NewService(),
	pluginService:     plugin.NewService(),
}

func registerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects",
		Description: "获取按标签分组的项目列表",
	}, listProjects)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects_data",
		Description: "分页获取项目列表，支持搜索",
	}, listProjectsData)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project",
		Description: "根据项目 ID 获取项目详情",
	}, getProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_project",
		Description: "创建新项目。必填 name 和 target，可选 tag、template、node 等",
	}, createProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tasks",
		Description: "分页获取扫描任务列表",
	}, listTasks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_task",
		Description: "根据任务 ID 获取任务详情",
	}, getTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_scan_templates",
		Description: "分页获取扫描模板列表",
	}, listScanTemplates)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_scan_template",
		Description: "根据模板 ID 获取扫描模板详情",
	}, getScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugin_modules",
		Description: "获取扫描模板的全部模块名（流水线阶段）。扫描模板由这些模块构成，每个模块下挂载若干插件。",
	}, listPluginModules)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugins",
		Description: "获取可用扫描插件，可按 module 过滤。返回每个插件的 hash、name、module、默认 parameter。创建扫描模板时，模块字段引用的就是插件 hash。",
	}, listPlugins)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_template",
		Description: createScanTemplateToolDesc,
	}, createScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_task",
		Description: "创建扫描任务。必填 name 和 node；template 必须是扫描模板的 ObjectID（用 list_scan_templates 获取，或 create_scan_template 的返回 id），不能用模板名。可用 target 直接给目标，或用 project/search 从已有资产选取目标。",
	}, createScanTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_assets",
		Description: listAssetsToolDesc,
	}, listAssets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_asset_detail",
		Description: "获取资产详情，asset_type 支持 asset 或 vulnerability",
	}, getAssetDetail)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_asset_tag",
		Description: "为资产添加标签",
	}, addAssetTag)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_nodes",
		Description: "获取扫描节点列表，online_only=true 时仅返回在线节点",
	}, listNodes)
}

type listProjectsDataInput struct {
	Search    string `json:"search,omitempty" jsonschema:"项目名称模糊搜索关键词"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"每页条数，默认 20"`
}

func listProjects(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsByTag(c)
	if err != nil {
		return errorResult("获取项目列表失败", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listProjectsData(ctx context.Context, _ *mcp.CallToolRequest, input listProjectsDataInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsData(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("获取项目数据失败", err)
	}
	return jsonToolResult(result)
}

type getProjectInput struct {
	ID string `json:"id" jsonschema:"项目 MongoDB ObjectID"`
}

func getProject(ctx context.Context, _ *mcp.CallToolRequest, input getProjectInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id 不能为空", nil)
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectContent(c, input.ID)
	if err != nil {
		return errorResult("获取项目详情失败", err)
	}
	if result == nil {
		return errorResult("项目不存在", nil)
	}
	return jsonToolResult(result)
}

type createProjectInput struct {
	Name           string   `json:"name" jsonschema:"项目名称，必填"`
	Tag            string   `json:"tag,omitempty" jsonschema:"项目标签，用于分组"`
	Target         string   `json:"target" jsonschema:"扫描目标，必填。多行或逗号分隔，支持域名/IP/URL 等"`
	Template       string   `json:"template,omitempty" jsonschema:"关联扫描模板 ID"`
	Node           []string `json:"node,omitempty" jsonschema:"指定扫描节点名称列表"`
	AllNode        bool     `json:"allNode,omitempty" jsonschema:"是否使用全部节点"`
	Ignore         string   `json:"ignore,omitempty" jsonschema:"忽略目标列表，格式同 target"`
	Duplicates     string   `json:"duplicates,omitempty" jsonschema:"去重策略"`
	ScheduledTasks bool     `json:"scheduledTasks,omitempty" jsonschema:"是否启用定时扫描"`
	Hour           int      `json:"hour,omitempty" jsonschema:"定时扫描间隔（小时），仅在启用 scheduledTasks 时有效"`
}

func createProject(ctx context.Context, _ *mcp.CallToolRequest, input createProjectInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || input.Target == "" {
		return errorResult("name 和 target 不能为空", nil)
	}
	p := &models.Project{
		Name:           input.Name,
		Tag:            input.Tag,
		Target:         input.Target,
		Template:       input.Template,
		Node:           input.Node,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Tp:             "project",
	}
	c := ginContext(ctx)
	if err := d.projectService.AddProject(c, p); err != nil {
		return errorResult("创建项目失败", err)
	}
	return jsonToolResult(ginH{"success": true, "message": "项目创建成功"})
}

type listTasksInput struct {
	Search    string `json:"search,omitempty" jsonschema:"任务名称模糊搜索"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"每页条数，默认 20"`
}

func listTasks(ctx context.Context, _ *mcp.CallToolRequest, input listTasksInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	tasks, total, err := d.taskService.List(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("获取任务列表失败", err)
	}
	return jsonToolResult(ginH{"list": tasks, "total": total})
}

type getTaskInput struct {
	ID string `json:"id" jsonschema:"任务 MongoDB ObjectID"`
}

func getTask(ctx context.Context, _ *mcp.CallToolRequest, input getTaskInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id 不能为空", nil)
	}
	c := ginContext(ctx)
	result, err := d.taskService.GetTaskDetail(c, input.ID)
	if err != nil {
		return errorResult("获取任务详情失败", err)
	}
	if result == nil {
		return errorResult("任务不存在", nil)
	}
	return jsonToolResult(result)
}

type listScanTemplatesInput struct {
	Query     string `json:"query,omitempty" jsonschema:"模板名称模糊搜索"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"每页条数，默认 20"`
}

func listScanTemplates(ctx context.Context, _ *mcp.CallToolRequest, input listScanTemplatesInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	result, err := d.templateService.List(ctx, input.PageIndex, input.PageSize, input.Query)
	if err != nil {
		return errorResult("获取模板列表失败", err)
	}
	return jsonToolResult(result)
}

type getScanTemplateInput struct {
	ID string `json:"id" jsonschema:"扫描模板 MongoDB ObjectID"`
}

func getScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input getScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id 不能为空", nil)
	}
	c := ginContext(ctx)
	result, err := d.templateService.Detail(c, input.ID)
	if err != nil {
		return errorResult("获取模板详情失败", err)
	}
	return jsonToolResult(result)
}

type createScanTemplateInput struct {
	Name         string              `json:"name" jsonschema:"模板名称，必填"`
	Modules      map[string][]string `json:"modules,omitempty" jsonschema:"模块到插件 hash 列表的映射。key 为模块名(用 list_plugin_modules 获取)，value 为该模块下要启用的插件 hash 数组(用 list_plugins 获取，按数组顺序执行)。例 {\"SubdomainScan\":[\"d60ba73c...\"]}"`
	Parameters   map[string]map[string]string `json:"parameters,omitempty" jsonschema:"可选，覆盖插件运行参数。结构为 模块名->插件hash->参数字符串。不提供时自动使用插件默认参数"`
	VulList      []string            `json:"vullist,omitempty" jsonschema:"可选，nuclei POC 模板 ID 列表，仅 VulnerabilityScan 使用 nuclei 时有效"`
	TemplateJSON string              `json:"template_json,omitempty" jsonschema:"可选，完整 ScanTemplate JSON。提供时优先于 modules，用于高级自定义"`
}

func createScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input createScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" && input.TemplateJSON == "" {
		return errorResult("name 或 template_json 至少提供一个", nil)
	}

	var tmpl *models.ScanTemplate

	switch {
	case input.TemplateJSON != "":
		tmpl = &models.ScanTemplate{}
		if err := json.Unmarshal([]byte(input.TemplateJSON), tmpl); err != nil {
			return errorResult("template_json 格式无效", err)
		}
		if tmpl.Name == "" {
			tmpl.Name = input.Name
		}
	case len(input.Modules) > 0:
		built, err := buildTemplateFromModules(ctx, input)
		if err != nil {
			return errorResult(err.Error(), nil)
		}
		tmpl = built
	default:
		tmpl = &models.ScanTemplate{Name: input.Name}
	}

	if tmpl.Name == "" {
		return errorResult("模板名称不能为空", nil)
	}

	id, err := d.templateService.Save(ctx, "", tmpl)
	if err != nil {
		return errorResult("创建模板失败", err)
	}
	return jsonToolResult(ginH{"success": true, "id": id, "message": "模板创建成功，可用于 create_scan_task 的 template 参数"})
}

// buildTemplateFromModules 根据「模块->插件hash列表」组装扫描模板，并自动回填插件默认参数
func buildTemplateFromModules(ctx context.Context, input createScanTemplateInput) (*models.ScanTemplate, error) {
	validModule := make(map[string]bool, len(constants.PLUGINSMODULES))
	for _, m := range constants.PLUGINSMODULES {
		validModule[m] = true
	}

	templateMap := map[string]any{
		"name": input.Name,
	}
	paramMap := map[string]map[string]string{}

	for module, hashes := range input.Modules {
		if !validModule[module] {
			return nil, fmt.Errorf("无效的模块名: %s（用 list_plugin_modules 获取合法模块）", module)
		}
		if len(hashes) == 0 {
			continue
		}
		templateMap[module] = hashes

		modParams := map[string]string{}
		for _, hash := range hashes {
			// 优先使用调用方提供的参数覆盖
			if input.Parameters != nil {
				if mp, ok := input.Parameters[module]; ok {
					if v, ok := mp[hash]; ok {
						modParams[hash] = v
						continue
					}
				}
			}
			// 否则回填插件默认参数
			plg, err := d.pluginService.GetPluginByHash(ctx, hash)
			if err != nil || plg == nil {
				return nil, fmt.Errorf("插件 hash 不存在: %s（模块 %s）", hash, module)
			}
			if plg.Module != module {
				return nil, fmt.Errorf("插件 %s(hash=%s) 属于模块 %s，不能放入模块 %s", plg.Name, hash, plg.Module, module)
			}
			modParams[hash] = plg.Parameter
		}
		if len(modParams) > 0 {
			paramMap[module] = modParams
		}
	}

	if len(paramMap) > 0 {
		templateMap["Parameters"] = paramMap
	}
	if len(input.VulList) > 0 {
		templateMap["vullist"] = input.VulList
	}

	data, err := json.Marshal(templateMap)
	if err != nil {
		return nil, fmt.Errorf("组装模板失败: %w", err)
	}
	tmpl := &models.ScanTemplate{}
	if err := json.Unmarshal(data, tmpl); err != nil {
		return nil, fmt.Errorf("组装模板失败: %w", err)
	}
	return tmpl, nil
}

type createScanTaskInput struct {
	Name           string   `json:"name" jsonschema:"任务名称，必填，不可重复"`
	Target         string   `json:"target" jsonschema:"扫描目标。多行或逗号分隔的域名/IP/URL 等"`
	Node           []string `json:"node" jsonschema:"执行扫描的节点名称列表，必填"`
	Template       string   `json:"template,omitempty" jsonschema:"扫描模板 ID 或名称"`
	AllNode        bool     `json:"allNode,omitempty" jsonschema:"是否使用全部在线节点"`
	Ignore         string   `json:"ignore,omitempty" jsonschema:"忽略目标列表"`
	Duplicates     string   `json:"duplicates,omitempty" jsonschema:"去重策略"`
	Project        []string `json:"project,omitempty" jsonschema:"关联项目 ID 列表"`
	Search         string   `json:"search,omitempty" jsonschema:"从已有资产中选择目标的搜索表达式，语法同 list_assets 的 search"`
	ScheduledTasks bool     `json:"scheduledTasks,omitempty" jsonschema:"是否创建为计划任务"`
	Hour           int      `json:"hour,omitempty" jsonschema:"计划任务间隔-小时"`
	Minute         int      `json:"minute,omitempty" jsonschema:"计划任务间隔-分钟"`
	Day            int      `json:"day,omitempty" jsonschema:"计划任务间隔-天"`
	CycleType      string   `json:"cycleType,omitempty" jsonschema:"周期类型，如 hourly/daily/weekly"`
}

func createScanTask(ctx context.Context, _ *mcp.CallToolRequest, input createScanTaskInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || len(input.Node) == 0 {
		return errorResult("name 和 node 不能为空", nil)
	}
	c := ginContext(ctx)
	exists, err := d.taskService.CheckTaskNameExists(c, input.Name)
	if err != nil {
		return errorResult("检查任务名失败", err)
	}
	if exists {
		return errorResult("任务名已存在", nil)
	}

	taskModel := &models.Task{
		Name:           input.Name,
		Target:         input.Target,
		Node:           input.Node,
		Template:       input.Template,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		Project:        input.Project,
		Search:         input.Search,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Minute:         input.Minute,
		Day:            input.Day,
		CycleType:      input.CycleType,
	}
	taskID, err := d.taskCommonService.Insert(ctx, taskModel)
	if err != nil {
		return errorResult("创建扫描任务失败", err)
	}
	return jsonToolResult(ginH{"success": true, "id": taskID})
}

type listAssetsInput struct {
	AssetType        string              `json:"asset_type" jsonschema:"资产类型，必填。如 asset、RootDomain、subdomain、app、mp、UrlScan、SensitiveResult、DirScanResult、crawler、vulnerability、PageMonitoring、IPAsset、SubdomainTakerResult"`
	PageIndex        int                 `json:"pageIndex,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize         int                 `json:"pageSize,omitempty" jsonschema:"每页条数，默认 20"`
	SearchExpression string              `json:"search,omitempty" jsonschema:"搜索表达式（非 SQL）。= 模糊、== 精确、!= 排除，&&/|| 组合。通用字段 tag/task/rootDomain（project 不支持 search）；各类型专有字段见工具 description"`
	Filter           map[string][]string `json:"filter,omitempty" jsonschema:"精确过滤 JSON，与 search 可组合。filter.project 为项目 ObjectID（list_projects 获取，非项目名）；filter.task 为任务名称。详见工具 description"`
	Sort             map[string]string   `json:"sort,omitempty" jsonschema:"排序。仅 UrlScan/DirScanResult 支持 length：ascending 升序，其他值降序"`
	Sid              string              `json:"sid,omitempty" jsonschema:"敏感信息规则名称，仅 asset_type 为 SensitiveResult 时有效"`
}

func listAssets(ctx context.Context, _ *mcp.CallToolRequest, input listAssetsInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	query := models.SearchRequest{
		PageIndex:        input.PageIndex,
		PageSize:         input.PageSize,
		Index:            index,
		SearchExpression: input.SearchExpression,
		Sort:             input.Sort,
		Sid:              input.Sid,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		query.Filter = filter
	}

	c := ginContext(ctx)
	data, err := queryAssets(ctx, c, index, query)
	if err != nil {
		return errorResult("查询资产失败", err)
	}
	return jsonToolResult(data)
}

type getAssetDetailInput struct {
	AssetType string `json:"asset_type" jsonschema:"资产类型。详情查询支持 asset 或 vulnerability"`
	ID        string `json:"id" jsonschema:"资产 MongoDB ObjectID；vulnerability 类型传 hash 值"`
}

func getAssetDetail(ctx context.Context, _ *mcp.CallToolRequest, input getAssetDetailInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id 不能为空", nil)
	}
	c := ginContext(ctx)
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}

	switch index {
	case "asset":
		result, err := d.assetService.GetAssetByID(c, input.ID)
		if err != nil {
			return errorResult("获取资产详情失败", err)
		}
		if result == nil {
			return errorResult("资产不存在", nil)
		}
		return jsonToolResult(result)
	case "vulnerability":
		result, err := d.vulnService.GetVulnerabilityDetailByHash(c, input.ID)
		if err != nil {
			return errorResult("获取漏洞详情失败", err)
		}
		return jsonToolResult(result)
	default:
		return errorResult("asset_type 仅支持 asset 或 vulnerability 的详情查询", nil)
	}
}

type addAssetTagInput struct {
	AssetType string `json:"asset_type" jsonschema:"资产类型，同 list_assets 的 asset_type"`
	ID        string `json:"id" jsonschema:"资产 MongoDB ObjectID"`
	Tag       string `json:"tag" jsonschema:"要添加的标签名称"`
}

func addAssetTag(ctx context.Context, _ *mcp.CallToolRequest, input addAssetTagInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.ID == "" || input.Tag == "" {
		return errorResult("id 和 tag 不能为空", nil)
	}
	c := ginContext(ctx)
	req := &models.TagRequest{Type: index, ID: input.ID, Tag: input.Tag}
	if err := d.commonService.AddTag(c, req); err != nil {
		return errorResult("添加标签失败", err)
	}
	return jsonToolResult(ginH{"success": true})
}

type listNodesInput struct {
	OnlineOnly bool `json:"online_only,omitempty" jsonschema:"true 时仅返回在线节点，默认 false 返回全部"`
}

func listNodes(ctx context.Context, _ *mcp.CallToolRequest, input listNodesInput) (*mcp.CallToolResult, any, error) {
	result, err := d.nodeService.GetNodeData(ctx, input.OnlineOnly)
	if err != nil {
		return errorResult("获取节点列表失败", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listPluginModules(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	return jsonToolResult(ginH{"modules": constants.PLUGINSMODULES})
}

type listPluginsInput struct {
	Module string `json:"module,omitempty" jsonschema:"按模块名过滤，留空返回全部扫描插件。模块名用 list_plugin_modules 获取"`
	Search string `json:"search,omitempty" jsonschema:"按插件名称模糊搜索（仅在未指定 module 时生效）"`
}

type pluginBrief struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Module    string `json:"module"`
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	Status    bool   `json:"status"`
}

func toPluginBriefs(plugins []models.Plugin) []pluginBrief {
	briefs := make([]pluginBrief, 0, len(plugins))
	for _, p := range plugins {
		// 跳过服务端插件，它们不参与扫描流水线
		if p.Type == "server" {
			continue
		}
		briefs = append(briefs, pluginBrief{
			Hash:      p.Hash,
			Name:      p.Name,
			Module:    p.Module,
			Parameter: p.Parameter,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	return briefs
}

func listPlugins(ctx context.Context, _ *mcp.CallToolRequest, input listPluginsInput) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)

	if input.Module != "" {
		plugins, err := d.pluginService.ListByModule(c, input.Module)
		if err != nil {
			return errorResult("获取插件列表失败", err)
		}
		return jsonToolResult(ginH{"list": toPluginBriefs(plugins)})
	}

	resp, err := d.pluginService.List(c, &models.PluginListRequest{
		PageIndex: 1,
		PageSize:  500,
		Search:    input.Search,
	})
	if err != nil {
		return errorResult("获取插件列表失败", err)
	}
	return jsonToolResult(ginH{"list": toPluginBriefs(resp.List), "total": resp.Total})
}

func queryAssets(ctx context.Context, c *gin.Context, index string, query models.SearchRequest) (any, error) {
	switch index {
	case "asset":
		list, err := d.assetService.GetAssets(c, query)
		return ginH{"list": list}, err
	case "RootDomain":
		return d.rootDomainService.GetRootDomainData(c, query)
	case "subdomain":
		list, err := d.subdomainService.GetSubdomains(ctx, query)
		return ginH{"list": list}, err
	case "app":
		return d.appService.GetAppData(c, query)
	case "mp":
		return d.mpService.GetMPData(c, query)
	case "UrlScan":
		list, err := d.urlService.GetURLs(ctx, query)
		return ginH{"list": list}, err
	case "SensitiveResult":
		list, err := d.sensitiveService.GetSensitiveInfo(ctx, query)
		return ginH{"list": list}, err
	case "DirScanResult":
		list, err := d.dirscanService.List(ctx, query)
		return ginH{"list": list}, err
	case "crawler":
		list, err := d.crawlerService.GetCrawlers(ctx, query)
		return ginH{"list": list}, err
	case "vulnerability":
		return d.vulnService.GetVulnerabilities(ctx, query)
	case "PageMonitoring":
		return d.pageMonService.GetResult(c, query)
	case "IPAsset":
		return d.ipService.GetIPAssets(c, query)
	case "SubdomainTakerResult":
		return querySubdomainTaker(ctx, query)
	default:
		return nil, fmt.Errorf("不支持的资产类型: %s", index)
	}
}

func querySubdomainTaker(ctx context.Context, query models.SearchRequest) (any, error) {
	takerService := subdomain.NewTakerService()
	c := ginContext(ctx)
	return takerService.GetSubdomainTakerData(c, query)
}

func normalizeAssetIndex(assetType string) (string, error) {
	assetType = strings.TrimSpace(assetType)
	if assetType == "" {
		return "", fmt.Errorf("asset_type 不能为空")
	}
	aliases := map[string]string{
		"asset":                "asset",
		"web":                  "asset",
		"rootdomain":           "RootDomain",
		"root_domain":          "RootDomain",
		"root-domain":          "RootDomain",
		"subdomain":            "subdomain",
		"app":                  "app",
		"mp":                   "mp",
		"miniprogram":          "mp",
		"mini_program":         "mp",
		"url":                  "UrlScan",
		"urlscan":              "UrlScan",
		"sensitive":            "SensitiveResult",
		"sensitiveresult":      "SensitiveResult",
		"sensitive_result":     "SensitiveResult",
		"dirscan":              "DirScanResult",
		"dir_scan":             "DirScanResult",
		"dirscanresult":        "DirScanResult",
		"directory":            "DirScanResult",
		"crawler":              "crawler",
		"vulnerability":        "vulnerability",
		"vuln":                 "vulnerability",
		"pagemonitoring":       "PageMonitoring",
		"page_monitoring":      "PageMonitoring",
		"ip":                   "IPAsset",
		"ipasset":              "IPAsset",
		"subdomaintaker":       "SubdomainTakerResult",
		"subdomain_taker":      "SubdomainTakerResult",
		"subdomaintakerresult": "SubdomainTakerResult",
	}
	key := strings.ToLower(assetType)
	if v, ok := aliases[key]; ok {
		return v, nil
	}
	if _, ok := aliases[strings.ReplaceAll(key, "-", "_")]; ok {
		return aliases[strings.ReplaceAll(key, "-", "_")], nil
	}
	// 允许直接传 MongoDB 集合名
	valid := []string{"asset", "RootDomain", "subdomain", "app", "mp", "UrlScan",
		"SensitiveResult", "DirScanResult", "crawler", "vulnerability",
		"PageMonitoring", "IPAsset", "SubdomainTakerResult"}
	for _, v := range valid {
		if v == assetType {
			return v, nil
		}
	}
	return "", fmt.Errorf("不支持的 asset_type: %s", assetType)
}

type ginH map[string]any
