---
name: scopesentry-mcp
description: 通过 ScopeSentry MCP 服务管理安全扫描平台：项目、任务、扫描模板、各类资产与节点。在用户提到 ScopeSentry、MCP、API Key、扫描任务、资产查询时使用。
---

# ScopeSentry MCP

## 前置条件

1. ScopeSentry 服务已启动（默认 `http://localhost:8082`）
2. 已通过 JWT 登录后创建 API Key：`POST /api/apikey/create`，请求体 `{"name":"cursor"}`
3. 保存返回的 `key`（格式 `ssk_...`，仅显示一次）

## Cursor MCP 配置

在 Cursor Settings → MCP 中添加：

```json
{
  "mcpServers": {
    "scopesentry": {
      "url": "http://localhost:8082/mcp",
      "headers": {
        "X-API-Key": "ssk_你的密钥"
      }
    }
  }
}
```

也可使用 `Authorization: Bearer ssk_你的密钥`。

## 可用工具

| 工具 | 说明 |
|------|------|
| `list_projects` | 按标签分组的项目树 |
| `list_projects_data` | 分页项目列表 |
| `get_project` | 项目详情（需 `id`） |
| `create_project` | 创建项目（需 `name`、`target`） |
| `list_tasks` | 扫描任务列表 |
| `get_task` | 任务详情 |
| `list_scan_templates` | 扫描模板列表 |
| `get_scan_template` | 模板详情 |
| `list_plugin_modules` | 模板的全部模块名 |
| `list_plugins` | 可用扫描插件（含 hash、默认参数），可按 `module` 过滤 |
| `create_scan_template` | 创建模板（按 模块→插件hash 组装） |
| `create_scan_task` | 创建扫描任务（需 `name`、`node`，`template` 为模板 ID） |
| `list_assets` | 查询资产（需 `asset_type`） |
| `get_asset_detail` | 资产/漏洞详情 |
| `add_asset_tag` | 添加资产标签 |
| `list_nodes` | 扫描节点列表 |

## 资产类型 (asset_type)

`asset`、`RootDomain`、`subdomain`、`app`、`mp`、`UrlScan`、`SensitiveResult`、`DirScanResult`、`crawler`、`vulnerability`、`PageMonitoring`、`IPAsset`、`SubdomainTakerResult`

别名：`web`、`root_domain`、`url`、`sensitive`、`dirscan`、`vuln`、`ip` 等。

## 资产查询语法

### search 搜索表达式

| 运算符 | 含义 | 示例 |
|--------|------|------|
| `=` | 模糊匹配 | `domain=example` |
| `==` | 精确匹配 | `project=="某项目"` |
| `!=` | 排除 | `port!="80"` |
| `&&` | 与 | `project=="A" && port==443` |
| `\|\|` | 或 | `title=admin \|\| body=login` |

通用字段：`project`、`tag`、`task`、`rootDomain`

asset 类型常用字段：`domain`、`ip`、`port`、`service`、`app`、`title`、`statuscode`、`icon`

### filter 精确过滤

JSON 对象，同 key 多值为 OR。支持：`project`、`port`、`service`、`app`、`icon`、`statuscode`、`level`、`type`、`color`、`status`、`tags`、`task`、`sname`

```json
{"project": ["项目A"], "port": ["80", "443"]}
```

### sort 排序

```json
{"length": "-1"}
```

value 为 `1`（升序）或 `-1`（降序）。

## 扫描模板与插件

扫描模板由多个**模块**（流水线阶段）组成，每个模块下挂载若干**插件**。
模板的模块字段引用的是插件 **hash**（不是插件 id，也不是插件名）。

模块固定列表（`list_plugin_modules`）：`TargetHandler`、`SubdomainScan`、`SubdomainSecurity`、`PortScanPreparation`、`PortScan`、`PortFingerprint`、`AssetMapping`、`AssetHandle`、`URLScan`、`WebCrawler`、`URLSecurity`、`DirScan`、`VulnerabilityScan`、`PassiveScan`

### 创建模板工作流

1. `list_plugin_modules` 获取模块名
2. `list_plugins`（可传 `module` 过滤）获取每个模块下插件的 `hash` 和默认 `parameter`
3. `create_scan_template` 用 `modules` 指定「模块 → 插件 hash 列表」，工具会自动回填默认参数

```json
{
  "name": "我的模板",
  "modules": {
    "SubdomainScan": ["d60ba73c70aac430a0a54e796e7e19b8"],
    "PortScan": ["<rustscan的hash>"],
    "AssetMapping": ["<httpx的hash>"]
  }
}
```

如需自定义插件参数，用 `parameters`：`{"SubdomainScan": {"<hash>": "-t 20"}}`

### 创建任务

`create_scan_task` 的 `template` 必须是模板 **ObjectID**（来自 `list_scan_templates` 或 `create_scan_template` 返回的 `id`），不能用模板名。

## API Key 管理

- 列表：`GET /api/apikey/list`（需 JWT）
- 创建：`POST /api/apikey/create` `{"name":"名称"}`
- 删除：`POST /api/apikey/delete` `{"id":"..."}`
