package mcp

// 资产查询相关工具说明（写入 Tool.Description，供 MCP 客户端展示）
const listAssetsToolDesc = `查询 ScopeSentry 资产列表。

【接口与请求体】
与各资产 POST 接口一致（如 /api/assets/asset），body 为 models.SearchRequest:
- pageIndex, pageSize: 分页
- search: 搜索表达式字符串（对应前端 Csearch 搜索框 searchParams）
- filter: 精确过滤对象（对应前端项目/任务下拉、统计侧栏点击、表格列筛选，合并后传入）
- sort: 排序（多数类型无效）
- sid: 仅 SensitiveResult 敏感规则名

MCP 的 asset_type 会映射为后端 Index（如 asset→"asset"），再调用 helper.GetSearchQuery 生成 MongoDB 查询。

【前端参数分工（Csearch.vue + 各资产页）】
- search: 用户在搜索框输入的 DSL，如 domain=baidu && port==443
- filter.project: ElTreeSelect 选中值，为项目 ObjectID 数组（非项目名）
- filter.task: ElSelect 选中值，为任务名称数组（非任务 ID）；也可用动态标签 task=名称 写入 filter
- filter 还可来自: 统计侧栏点击(port/service/app/icon)、表格列筛选(statuscode/level/type/status 等)
- 注意: 各资产页 searchKeywordsData 里虽有 project 提示，但后端 SearchToMongoDB 未注册 project 关键字，project 只能走 filter

asset_type 可选值:
- asset: Web/端口资产
- RootDomain: 根域名
- subdomain: 子域名
- app: 移动应用
- mp: 小程序
- UrlScan: URL
- SensitiveResult: 敏感信息
- DirScanResult: 目录扫描
- crawler: 爬虫
- vulnerability: 漏洞
- PageMonitoring: 页面监控
- IPAsset: IP 聚合资产
- SubdomainTakerResult: 子域接管

【search 搜索表达式】（非 SQL，为自定义 DSL，由 helper.SearchToMongoDB 解析）
- field=value : 模糊匹配（不区分大小写）
- field=="值" : 精确匹配，含空格的值需加双引号
- field!="值" : 排除匹配
- expr1 && expr2 : 与
- expr1 || expr2 : 或
- (expr) : 分组
- 通用 search 字段(所有类型，SearchToMongoDB 注入): tag→tags、task→taskName（任务名）、rootDomain
- project 不支持 search，只能通过 filter.project（值为项目 ObjectID，list_projects / list_projects_data 获取）

各 asset_type 可用 search 关键字 → MongoDB 字段:
- asset: domain→host, ip, port, service, app→technologies, title, statuscode, icon→faviconmmh3, banner→metadata, type, body, header→rawheaders
- RootDomain: domain, icp, company
- subdomain: domain→host, ip, type, value
- app: name, icp, company, category, description, url, apk
- mp: name, icp, company, category, description, url
- UrlScan: url→output, input, source, resultId, type→outputtype（无 statuscode search；HTTP 状态码仅 filter.status）
- SensitiveResult: url, sname→sid, body, info→match, md5
- DirScanResult: url, statuscode→status, redirect→msg, length
- vulnerability: url, vulname, matched, request, response, level
- crawler: url, method, body, resultId
- PageMonitoring: url, hash, diff, response
- IPAsset: ip, domain→ports.server.domain, port→ports.port, service→ports.server.service, webServer→ports.server.webServer, app→ports.server.technologies
- SubdomainTakerResult: domain→input, value, type→cname, response

search 示例:
- domain=baidu && port==443
- port==443 && service=nginx
- title="后台登录" || body=admin

【filter 精确过滤】（与 search 可组合；由 helper.GetSearchQuery 追加 $and 条件）
JSON 对象，key 为筛选维度，value 为字符串数组（同 key 多值为 OR，不同 key 之间为 AND）。
仅 filterKeyCache 中定义的 key 生效，其余 key 会被忽略。

filter 要点:
- project: 仅 filter，值为项目 ObjectID 数组（数据库存 ID；展示时部分接口会转成项目名）
- task: 可用 search（task=="任务名"）或 filter.task（任务名数组）；值为 list_tasks 的 name，非任务 ID
- 同 key 多值 OR，不同 key AND；未在 filterKeyCache 的 key 被忽略

全局 filter key → MongoDB 字段:
- project→project, port→port, service→service, app→technologies
- icon→faviconmmh3, statuscode→statuscode, status→status
- level→level, type→type, color→color, tags→tags
- task→taskName, sname→sid

各 asset_type 可用 filter（未列出的 key 对该集合无效或字段不存在）:
- asset: project, port, service, app, icon, statuscode, type, task, tags
- RootDomain: project, tags
- subdomain: project, type, task, tags
- app / mp: project, tags
- UrlScan: status（HTTP 状态码，非 statuscode）, tags
- SensitiveResult: status（1未处理/2处理中/3忽略/4疑似/5确认/6已处理）, color, sname, tags
- DirScanResult: status（HTTP 状态码）, tags
- crawler: project, task, tags
- vulnerability: project, level（critical/high/medium/low/info/unknown）, status（1-6）, task, tags
- PageMonitoring: tags
- IPAsset: project, port, service, app（嵌套 ports 字段，走聚合查询）
- SubdomainTakerResult: tags

filter 示例（project 必须为 ObjectID，不可用项目名）:
- asset: {"project":["<项目ObjectID>"],"port":["443"]}
- subdomain: {"project":["<项目ObjectID>"],"type":["A"]}
- vulnerability: {"project":["<项目ObjectID>"],"level":["high"]}

search + filter 组合示例:
- search: domain=baidu && port==443，filter: {"project":["<项目ObjectID>"]}

注意:
- 勿在 search 中写 project=... 或 project=="..."（无效或与 && 组合报错）
- DirScanResult 的 HTTP 状态码仅在 search 中用 statuscode（如 statuscode==200）
- UrlScan 无 statuscode search；状态筛选用 filter.status
- SensitiveResult 规则名筛选：search 用 sname=规则名，或 filter.sname

【sort 排序】
- UrlScan / DirScanResult: 支持 {"length":"ascending"} 升序，其他值（含 "descending"、"-1"）为降序
- 其余类型: 服务端固定按 time 或 _id 排序，sort 参数通常无效

【sid 参数】
仅 SensitiveResult: 传敏感规则名称（sid），用于按规则展开匹配详情（对应前端点击规则名）`

// createScanTemplateToolDesc 扫描模板创建说明
const createScanTemplateToolDesc = `创建扫描模板。

扫描模板由多个【模块】(流水线阶段)组成，每个模块下挂载若干【插件】，
模块字段引用的是插件的 hash（不是插件 id，也不是插件名）。

推荐创建流程:
1. list_plugin_modules 获取全部模块名
2. list_plugins (可按 module 过滤) 获取每个模块下可用插件的 hash 与默认 parameter
3. 用 modules 参数指定 模块->插件hash列表 (同模块内按数组顺序执行)
4. 本工具会自动按插件默认参数回填 Parameters；如需自定义参数用 parameters 覆盖

参数说明:
- name: 模板名称，必填
- modules: {模块名: [插件hash,...]}，如 {"SubdomainScan":["d60ba73c..."],"PortScan":["..."]}
- parameters: 可选，{模块名: {插件hash: 参数字符串}}，覆盖默认参数
- vullist: 可选，nuclei POC 模板 ID 列表 (VulnerabilityScan 用 nuclei 时)
- template_json: 可选，完整 ScanTemplate JSON，优先级最高，用于高级自定义

创建成功返回模板 id，可直接用于 create_scan_task 的 template 参数。

常见模块: TargetHandler, SubdomainScan, SubdomainSecurity, PortScanPreparation,
PortScan, PortFingerprint, AssetMapping, AssetHandle, URLScan, WebCrawler,
URLSecurity, DirScan, VulnerabilityScan, PassiveScan`
