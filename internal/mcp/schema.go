package mcp

// 资产查询相关工具说明（写入 Tool.Description，供 MCP 客户端展示）
const listAssetsToolDesc = `查询 ScopeSentry 资产列表。

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

【search 搜索表达式语法】
- field=value : 模糊匹配（不区分大小写）
- field=="值" : 精确匹配，含空格的值需加双引号
- field!="值" : 排除匹配
- expr1 && expr2 : 与
- expr1 || expr2 : 或
- (expr) : 分组
- 通用字段(所有类型): project、tag、task、rootDomain

各 asset_type 可用 search 字段:
- asset: domain, ip, port, service, app, title, statuscode, icon, banner, type, body, header
- RootDomain: domain, icp, company
- subdomain: domain, ip, type, value
- app/mp: name, icp, company, category, description, url (app 另有 apk)
- UrlScan: url, input, source, resultId, type
- SensitiveResult: url, sname, body, info, md5
- DirScanResult: url, statuscode, redirect, length
- vulnerability: url, vulname, matched, request, response, level
- crawler: url, method, body, resultId
- PageMonitoring: url, hash, diff, response
- IPAsset: ip, domain, port, service, webServer, app
- SubdomainTakerResult: domain, value, type, response

示例:
- project=="某项目" && domain=example
- port==443 && service=nginx
- title="后台登录" || body=admin

【filter 精确过滤】
JSON 对象，key 为维度，value 为字符串数组（同维度多值为 OR）。
支持 key: project, port, service, app, icon, statuscode, level, type, color, status, tags, task, sname
示例: {"project":["项目A"],"port":["80","443"]}

【sort 排序】
JSON 对象，value 为 "1"(升序) 或 "-1"(降序)。部分类型支持 length 等字段。
示例: {"length":"-1"}`

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
