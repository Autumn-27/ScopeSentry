---
name: scopesentry-mcp
description: 通过 ScopeSentry MCP 管理安全扫描平台（项目、任务、模板、资产、节点）。在用户提到 ScopeSentry、MCP、API Key、扫描任务、资产查询时使用。
---

# ScopeSentry MCP 使用指南

面向**已部署 ScopeSentry 实例**的用户。通过 Cursor（或其他 MCP 客户端）连接平台，无需本地源码。

## 1. 准备工作

### 1.1 确认服务可访问

- 默认 Web 界面：`http://<主机>:8082`
- MCP 端点：`http://<主机>:8082/mcp`（若前面有反向代理或前端代理，以实际 `/mcp` 地址为准）

### 1.2 创建 API Key

1. 浏览器登录 ScopeSentry Web 界面
2. 进入 **API Key** 管理页创建密钥（或通过管理员提供的接口创建）
3. 保存返回的 `ssk_...` 字符串（**仅显示一次**）

### 1.3 配置 Cursor MCP

Cursor → Settings → MCP → 添加服务器：

```json
{
  "mcpServers": {
    "scopesentry": {
      "url": "http://<你的主机>:8082/mcp",
      "headers": {
        "X-API-Key": "ssk_你的密钥"
      }
    }
  }
}
```

也可使用：`Authorization: Bearer ssk_你的密钥`

配置完成后重启 MCP 或重载 Cursor，确认工具列表中出现 `list_projects`、`list_assets` 等。

---

## 2. 工具一览

| 工具 | 用途 |
|------|------|
| `list_projects` | 按标签分组的项目树（含项目 ID） |
| `list_projects_data` | 分页项目列表，可按名称搜索 |
| `get_project` | 项目详情 |
| `create_project` | 新建项目 |
| `list_tasks` | 扫描任务列表 |
| `get_task` | 任务详情 |
| `list_scan_templates` | 扫描模板列表 |
| `get_scan_template` | 模板详情 |
| `list_plugin_modules` | 扫描流水线模块名 |
| `list_plugins` | 可用插件（含 hash、默认参数） |
| `create_scan_template` | 创建扫描模板 |
| `create_scan_task` | 创建扫描任务 |
| `list_assets` | 查询各类资产 |
| `get_asset_detail` | 资产或漏洞详情 |
| `add_asset_tag` | 为资产添加标签 |
| `list_nodes` | 扫描节点列表 |

各工具参数以 MCP 工具描述（schema）为准；`list_assets` 的字段说明最完整，查询资产前可先阅读该工具 description。

---

## 3. 常用工作流

### 3.1 按项目查资产

1. `list_projects` 或 `list_projects_data` 获取目标项目的 **ObjectID**（`id` / `children[].value`）
2. `list_assets` 传入 `filter.project`（**必须是 ID，不能写项目中文名**）

```json
{
  "asset_type": "asset",
  "pageIndex": 1,
  "pageSize": 20,
  "search": "domain=example",
  "filter": {
    "project": ["<项目ObjectID>"]
  }
}
```

### 3.2 创建扫描任务

1. `list_nodes` 获取在线节点名称
2. `list_scan_templates` 或 `create_scan_template` 获取模板 **ObjectID**
3. `create_scan_task`：`name`、`node` 必填，`template` 填模板 ID（不能填模板名）

### 3.3 创建扫描模板

1. `list_plugin_modules` → 模块名列表
2. `list_plugins`（可按 `module` 过滤）→ 各插件 `hash` 与默认 `parameter`
3. `create_scan_template`：用 `modules` 指定「模块 → 插件 hash 数组」

---

## 4. 资产查询（`list_assets`）

### 4.1 资产类型 `asset_type`

`asset`、`RootDomain`、`subdomain`、`app`、`mp`、`UrlScan`、`SensitiveResult`、`DirScanResult`、`crawler`、`vulnerability`、`PageMonitoring`、`IPAsset`、`SubdomainTakerResult`

别名示例：`web`→asset、`vuln`→vulnerability、`ip`→IPAsset、`url`→UrlScan

### 4.2 参数说明

| 参数 | 说明 |
|------|------|
| `pageIndex` / `pageSize` | 分页，默认 1 / 20 |
| `search` | 搜索表达式（见下节） |
| `filter` | 精确过滤 JSON（见下节） |
| `sort` | 仅 UrlScan、DirScanResult 支持按 `length` 排序 |
| `sid` | 仅 SensitiveResult：敏感规则名称 |

`search` 与 `filter` **可同时使用**。

### 4.3 search 搜索表达式

自定义 DSL（**不是 SQL**）：

| 运算符 | 含义 | 示例 |
|--------|------|------|
| `=` | 模糊匹配 | `domain=example` |
| `==` | 精确匹配 | `port==443` |
| `!=` | 排除 | `port!="80"` |
| `&&` | 与 | `domain=baidu && port==443` |
| `\|\|` | 或 | `title=admin \|\| body=login` |

**所有类型通用 search 字段：** `tag`、`task`（任务名称）、`rootDomain`

**project 不能写在 search 里**（无效或与 `&&` 组合时报错）。筛项目请用 `filter.project`。

**各类型常用 search 字段：**

| asset_type | 字段 |
|------------|------|
| asset | domain, ip, port, service, app, title, statuscode, icon, banner, type, body, header |
| RootDomain | domain, icp, company |
| subdomain | domain, ip, type, value |
| app | name, icp, company, category, description, url, apk |
| mp | name, icp, company, category, description, url |
| UrlScan | url, input, source, resultId, type |
| SensitiveResult | url, sname, body, info, md5 |
| DirScanResult | url, statuscode, redirect, length |
| vulnerability | url, vulname, matched, request, response, level |
| crawler | url, method, body, resultId |
| PageMonitoring | url, hash, diff, response |
| IPAsset | ip, domain, port, service, webServer, app |
| SubdomainTakerResult | domain, value, type, response |

**search 示例：**

- `domain=baidu && port==443`
- `task=="某任务名"`
- `level==high`（vulnerability）
- `statuscode==200`（DirScanResult）

### 4.4 filter 精确过滤

JSON 对象：同 key 多个值为 **OR**，不同 key 为 **AND**。

| filter key | 含义 | 取值说明 |
|------------|------|----------|
| `project` | 所属项目 | **ObjectID**，用 `list_projects` / `list_projects_data` 获取 |
| `task` | 来源任务 | **任务名称**，用 `list_tasks` 的 `name` |
| `port` | 端口 | 如 `"443"` |
| `service` | 服务/协议 | 如 `"https"` |
| `app` | 应用指纹 | 如 `"Nginx"` |
| `icon` | 图标 hash | |
| `statuscode` | HTTP 状态码 | 主要用于 asset |
| `status` | 状态 | UrlScan/DirScan HTTP 码；漏洞/敏感信息处理状态 |
| `level` | 漏洞等级 | critical / high / medium / low / info |
| `type` | 类型 | 如子域名记录类型 A、CNAME |
| `color` | 敏感规则颜色 | SensitiveResult |
| `sname` | 敏感规则名 | SensitiveResult |
| `tags` | 标签 | |

**各类型可用 filter key：**

| asset_type | filter key |
|------------|------------|
| asset | project, port, service, app, icon, statuscode, type, task, tags |
| RootDomain | project, tags |
| subdomain | project, type, task, tags |
| app / mp | project, tags |
| UrlScan | status, tags |
| DirScanResult | status, tags |
| SensitiveResult | status, color, sname, tags |
| crawler | project, task, tags |
| vulnerability | project, level, status, task, tags |
| PageMonitoring / SubdomainTakerResult | tags |
| IPAsset | project, port, service, app |

**filter 示例：**

```json
{"project": ["<项目ObjectID>"], "port": ["443"]}
```

**组合查询示例：**

```json
{
  "asset_type": "asset",
  "search": "domain=baidu && port==443",
  "filter": {"project": ["<项目ObjectID>"]},
  "pageIndex": 1,
  "pageSize": 10
}
```

**注意：**

- `filter.project` 勿填项目显示名称
- UrlScan 的 HTTP 状态用 `filter.status`；DirScanResult 可在 search 中用 `statuscode==200`
- SensitiveResult 按规则名：`search` 用 `sname=规则名`，或 `filter.sname`

### 4.5 排序 sort

仅 **UrlScan**、**DirScanResult** 支持：

```json
{"length": "ascending"}
```

其他类型忽略 `sort`，按时间默认排序。

---

## 5. 扫描模板模块名

`TargetHandler`、`SubdomainScan`、`SubdomainSecurity`、`PortScanPreparation`、`PortScan`、`PortFingerprint`、`AssetMapping`、`AssetHandle`、`URLScan`、`WebCrawler`、`URLSecurity`、`DirScan`、`VulnerabilityScan`、`PassiveScan`

---

## 6. 故障排查

| 现象 | 处理 |
|------|------|
| MCP 无工具 | 检查 URL、API Key、ScopeSentry 是否运行 |
| 401 / 403 | 重新创建或更换 API Key |
| 资产查不到 | 确认 `filter.project` 为 ObjectID；勿在 search 写 project |
| 模板/任务创建失败 | `template` 必须是模板 ObjectID；`node` 填在线节点名 |
| 查询很慢 | 缩小 `pageSize`，增加 search/filter 条件 |

---
