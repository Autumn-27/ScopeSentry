# ScopeSentry 上游开放 Issue 分诊

| 项 | 值 |
| --- | --- |
| 仓库 | https://github.com/Autumn-27/ScopeSentry |
| 分诊基线 | `docs/issue-triage` @ `e16bbcf`（含 PR #343/#344/#345） |
| 版本 | server `1.9.4` / scan `1.9.1`（`version.json`） |
| 样本 | 开放 issue 88 条（另有 3 条开放 PR，不计入） |
| 方法 | `gh issue list` + `/issues/<n>`、`/issues/<n>/comments` 取正文与讨论；逐条在 worktree 内只读定位代码 |
| 说明 | 指针形如 `path:line`，均以本基线文件行号为准；“scan 节点”指独立仓库 `Autumn-27/ScopeSentry-Scan` |

## 分桶总览

| 桶 | 条数 | 处置 |
| --- | --- | --- |
| 已由合并的 PR 修复 | 1 | 直接关闭 |
| 本次已修复 | 1 | 回复“已修”后关闭剩余部分另开 |
| 可修-后端bug | 10 | 本轮/下轮排期 |
| 可修-前端bug | 3 | 本轮/下轮排期 |
| 需要新功能(设计) | 39 | 进 roadmap，逐条拆需求 |
| 插件需求 | 17 | 转插件生态，不进内核 |
| 文档 | 5 | 补官网/README 即可关闭 |
| 外部仓库(scan 节点, 不在本仓库) | 7 | 转 `ScopeSentry-Scan` |
| 信息不足(标题只有一句话) | 5 | 回帖要复现信息 |
| **合计** | **88** | |

## 已由合并的 PR 修复

| Issue | 标题 | PR / commit | 证据 |
| --- | --- | --- | --- |
| #327 | 前端敏感信息管理报错（关闭规则时 400） | PR #328 / `db0fd09` | `internal/api/handlers/sensitiveRule/sensitive.go:34,41` 由 `bool` 改 `*bool`，`state:false` 不再被 `required` 拒；`sensitive_test.go:11-24` 有回归用例 |

## 本次已修复

| Issue | 已修范围 | 证据 | 遗留 |
| --- | --- | --- | --- |
| #342 | 问题 2：`targetSource=UrlScan` 触发服务端 panic | PR #343 / `4626678`：`internal/services/task/common/common.go:55-62` 补 `urlService: url.NewService()`；`common_test.go:13-22` 断言无 nil 依赖 | 问题 1（nuclei `parameter` 不生效）、问题 3（模板目录 `/apps/poc` 写死）在 scan 节点仓库，见深度分析 |

## 可修-后端bug

| # | 标题 | 影响面 | 代码指针 | 难度 | 建议修法 |
| --- | --- | --- | --- | --- | --- |
| 317 | 同步到项目提示 request failed | 高：任务结果无法进项目，核心工作流断；日志里 `404 POST /api/task/sync` | `frontend/src/api/task/index.ts:314`；`internal/api/routes/task/task.go:16-72`（无 sync 路由） | M | 补 `POST /api/task/sync` handler（把目标/根域名/ICP/公司名写入项目），或前端改调已有接口并同步文档 |
| 242 | 同步任务到项目后 rootdomain/app/mp 全空 | 高：项目资产关联被误清空 | `internal/repositories/project/project.go:367-401`（先 `$set project:""` 再回填）、`429-459` | M | 清空条件限制在本项目且改为精确集合差，先算“应保留集合”再批量 `$in`，避免 regex 误命中 |
| 228 | 同步到项目时换行/ICP 公司名未加前缀 | 中：项目归属匹配不上 | `internal/services/project/project.go:236-264`、`376-390`（`getBeforeLastDash`） | S | 解析前 `TrimSpace` 去空行，ICP 名统一补 `ICP:` 前缀后再入库并参与匹配 |
| 232 | 资产过大导致拼接 rootdomain 正则过长 | 中：大项目更新超时/失败 | `internal/repositories/project/project.go:429-434,462-467,495-525` | L | 改为 `domain→project` 哈希表或 `$in` 精确匹配；过渡期按 500 条分片查询 |
| 340 | 指纹管理处一键更新后每个指纹重复两份 | 中：所有做过一键更新的实例必现 | `internal/services/fingerprint/fingerprint.go:174-181`（仅按 `fingerprint_id` upsert）、`internal/api/handlers/fingerprint/fingerprint.go:181-192`（`Add` 写随机 `SAA_` id） | S | upsert 前按 `name` 去重/回填 `fingerprint_id`，并给 `FingerprintRules.name` 建唯一索引（`initdb.go:243` 处补） |
| 262 | 开启被动扫描后任务提前显示结束 | 中：用户以为扫完，实际被动扫描仍在跑 | `internal/services/task/task/task.go:528-568`（`progress==100` 即 `status=3`） | S | 完成判定加入 PassiveScan 未结束条件，或为被动扫描单列进度位 |
| 295 | 删除未认证的接口 | 高：未鉴权即可删除扫描模板 | `internal/api/routes/task/template.go:25-29`（唯一缺 `WithAuth()` 的业务路由） | S | 补 `Middlewares: common.WithAuth()`；顺带核对其余路由的权限位 |
| 254 | 分布式节点总是离线 | 高：节点反复掉线需重启 | `internal/services/node/service.go:24-26,59-83,102-110` | M | 心跳改用 Redis TTL 而非“节点自报时间 vs 服务器时间”，容差可配，收到心跳立即恢复 state=1 |
| 297 | 获取本地缓存任务时仍用旧扫描模板 | 中：模板改了不生效，复扫用旧参数 | `internal/services/task/common/common.go:213-220`（入队即快照）、`internal/repositories/task/task/task.go:179-190,107-120` | M | 队列 payload 带模板版本，节点消费时回源校验；或模板保存时清理未消费的 `NodeTask:*` |
| 169 | 漏洞增加 LastScanTime 字段 | 中：无法判断漏洞是否仍存在 | `internal/models/vulnerability.go:8-24`（只有 `time`） | S | 加 `lastScanTime`，漏洞 upsert 命中时刷新该字段并在列表暴露 |

## 可修-前端bug

| # | 标题 | 影响面 | 代码指针 | 难度 | 建议修法 |
| --- | --- | --- | --- | --- | --- |
| 225 | csv 导出格式优化 | 低但显眼：导出的其实是 xlsx，后缀/文案却是 csv；不改单选项还会 400 | `frontend/src/views/Asset/export/exportData.vue:221,255`；`internal/api/handlers/export/export.go:55-57` | S | 默认值与单选文案统一成 `xlsx`，后端错误文案同步改为 “xlsx or json” |
| 288 | 定时任务运行状态不对（未开始也显示运行中） | 中：误判任务在跑 | `frontend/src/views/Task/Task.vue:90-114`（只映射 1/2/3）；`internal/services/task/common/common.go:148-151`（Insert 一律 `Status=1`） | S/M | 增加 `status=0 待开始` 与前端标签；计划任务列表用 `lastTime/nextTime` 推导状态 |
| 216 | 资产搜索界面优化 | 中：直接输裸域名 500 | `internal/utils/helper/search.go:423+`（`SearchToMongoDB` 解析失败即报错）；资产搜索页提交处 | S | 前端提交前把裸 token 归一成 `domain=xxx`，后端解析失败降级为模糊匹配而非 500 |

## 需要新功能(设计)

| # | 标题 | 现状与建议 |
| --- | --- | --- |
| 321 | 分布式优化建议 | 现需手工 compose；已有 `k8s/` 清单（PR #344）可复用，建议出一键安装脚本/agent |
| 335 | 资产建议新增备注字段 | `internal/models/asset.go:8-40` 无 remark；加字段 + 列表编辑入口 |
| 334 | 资产类别新增快应用 | ICP 分类只有 域名/小程序/APP；需新增 quickapp 集合与统计口径 |
| 337 | 如何配置全局代理 | 与 #266 同一需求；需配置项 + 节点注入插件进程环境变量 |
| 266 | 全局代理 | 扫描流量走代理需节点侧支持；SYN 端口扫描无法走 HTTP 代理（#266 评论） |
| 278 | 进度详情 100% 时迁入 mongodb | 完成即删 `TaskInfo:tmp/time`（`task.go:556-560`），详情随之丢失；建议完成时快照落库 |
| 277 | 可以选择目标暂停 | 现仅整任务 stop（`task.go:243-257`），需 target 粒度控制 |
| 273 | 插件参数设为全局变量 | 参数只存在 `ScanTemplate.Parameters`（`models/template.go:26,45`），建议加全局默认层 |
| 270 | 域名泛解析全局过滤 | 无 wildcard 过滤逻辑（grep 仅命中 CT 模块），需在子域/资产入库前判定 |
| 261 | 增加时间段导出 | 导出仅按筛选条件（`handlers/export/export.go:65-75`），缺时间区间与“已读”标记 |
| 259 | 扫描器重启总是重新加载字典 | 建议加字典 hash 接口，节点比对一致则跳过下载 |
| 258 | 删除全部（17w 条只能 1000 条删） | `handlers/assets/common/common.go:34-54` 只接受 ids；需“按筛选条件删除”后台任务 |
| 239 | 新增通知（webhook robot / 端口告警） | 现有通知配置仅基础项（`internal/config/notification.go`），需 webhook 通道 + 端口事件 |
| 230 | 动态解析页面与跳转 URL 扫描 | 跳转链依赖 katana 上游（作者已改 katana，#230 评论），服务端只能扩范围配置 |
| 224 | 通过 IP 取 ASN 再取网段 | 无 ASN 代码（grep 0 命中），需接第三方 API + CIDR 展开 |
| 223 | 可自动化实现的 case | 参数收集入 Redis + 缓存投毒检测，属新检测能力 |
| 217 | 标签管理（颜色） | `models/asset.go:33` 有 tags 无颜色；可复用敏感规则的 color 设计 |
| 203 | 资产管理（手工录入/任务创建即同步） | 服务端无资产写入接口（仅 initdb 写），需新增 add 接口 + 前端表单 |
| 200 | redis 存储优化 | task 相关 key 无 TTL（`repositories/task/task/task.go:107-120` 只删不设过期），进程态建议落 mongo |
| 191 | 一键删除功能 | 同 #258，误报清理需要批量删除入口 |
| 188 | 扫描包 header 修改 | 作者已答“可以支持”（#188 评论）；在模板/插件参数层加自定义 header |
| 182 | 通知设置默认发测试消息 | save 时发一条测试消息；notification 刷新通道已存在（`configuration/service.go:167-208`） |
| 179 | 被动扫描改开关 | `PassiveScan` 已在模块列表（`constants/defaults.go:352-367`），缺独立开关与端口配置 |
| 175 | POC 批量删除/设置 Tag | 批量删除已有（`Poc.vue:255`），缺批量打 Tag |
| 168 | 扫描过程中过滤不扫描域名 | `task.Ignore` 仅入口过滤（`helper/task.go:107-149`），派生资产未过滤 |
| 165 | 导入目标数据 | 任务 add 无文件导入（`handlers/task` 无 Import）；插件开发文档作者已给 demo 链接 |
| 159 | 自带 nuclei POC 缺变量不可用 | 需 POC 导入时校验 variables 并在 UI 标记不可用 |
| 157 | 线程运行逻辑优化 | 建议按模块批量投递（收 N 个目标再跑单进程），降内存与进程数 |
| 140 | 云主机自动创建节点 | 需云厂商 API + 节点注册自动化 |
| 137 | mongodb 认证加密参数 | `config.go:46-52,276` 只有 user/pass；缺 authSource/authMechanism/TLS |
| 132 | 内网域名不能扫描 | 作者定位为 DNS 校验（默认外网 DNS 导致内网域名丢失，#132 评论）；建议加“跳过 DNS 校验/自定义 resolver” |
| 117 | 增加工具、反弹 shell | 内置工具箱与 shell 通道，属新能力，需先定安全边界 |
| 116 | 权限划分 | RBAC 骨架已在（`routes/role`、`constants/permissions.go`），user 路由被注释（`routes/user/user.go:29-56`） |
| 115 | quake/fofa/hunter 页面调用、去重、导入 | 现仅配置 key（configuration），缺页面查询/去重/导入流程 |
| 113 | 代理池 | 需代理池管理与插件参数注入 |
| 109 | 一键更新 | 现仅启动时 DB 迁移（`internal/update`），无 UI 升级；文档仍要求手工 docker pull |
| 91 | 项目树状分组 | 现有 project 为平铺，需层级模型 |
| 86 | 节点自动选择/可达性探测 | 建任务前探测节点可达并自动分配；与 #285“节点主动拉取”可一并解决 |
| 285 | redis 相关问题（节点分发/日志） | 见深度分析：改节点主动拉取 + Redis Stream 日志 |

## 插件需求

| # | 插件/能力 | 说明 |
| --- | --- | --- |
| 284 | 新端口扫描方式（TXPortMap） | 属 scan 节点插件；建议先对比 rustscan 全端口实测 |
| 209 | 403 bypass 插件 | 属 URLSecurity 类插件 |
| 181 | HeartK（JS 分析） | 属 URLSecurity/工具类插件 |
| 150 | pocsuite3 兼容 | 需在节点侧封装 pocsuite3 调用 |
| 147 | 基于 urlscan 的目录扫描 | DirScan 插件，输入改为 UrlScan 结果 |
| 144 | 资产分拣插件 | 按资产特征自动分类/打标 |
| 131 | ffuf | DirScan 插件 |
| 122 | awvs | 外部扫描器对接 |
| 121 | afrog | 漏洞扫描插件 |
| 120 | xpoc | 漏洞扫描插件 |
| 114 | oneforall | 子域名插件 |
| 111 | host 碰撞插件 | 同 #301 |
| 110 | 多协议账号密码爆破 | 与 #61 同族 |
| 101 | 按标题清理误报的插件 | 作者答“目前只有扫描端插件，服务端后边考虑”（#101 评论）→ 需先支持服务端插件 |
| 301 | host 碰撞功能 | 同 #111，可加大收集面 |
| 304 | ssrf 插件（插件市场里消失） | 需确认是否下架，恢复或说明替代 |
| 61 | 常见服务弱口令爆破 + 自定义字典 | 已有 spray 插件（#254 评论），缺字典自定义与多协议模板 |

## 文档

| # | 主题 | 建议 |
| --- | --- | --- |
| 283 | 国内服务器拉镜像超时 | 补镜像加速/离线包部署章节（作者已回复排查方向） |
| 280 | 加载漏洞扫描插件时无 templates、无漏洞产出 | 说明“新版本不再内置 POC，需自行导入”，并补 POC 导入路径与 `/apps/poc` 约定 |
| 187 | 官网文档不全（页面监控/并发/分发） | 作者已有大段答复（#187 评论），整理进官网“使用教程/并发说明” |
| 45 | 常见问题（部署/进度/范围） | 把已沉淀的问答（compose 版本、进度算法、tlsprobe 串域）编成 FAQ |
| 247 | xray+rad+katana+wayback 链路排障 | 补“插件流量是否走代理/如何看插件日志”的排障章节 |

## 外部仓库(scan 节点, 不在本仓库)

| # | 现象 | 判定依据 |
| --- | --- | --- |
| 330 | `logger.SlogErrorLocal(fmt.Sprintf("Failed to get CPU usage:", err))` 多传参数 | 本仓库 grep `SlogErrorLocal`/`Failed to get CPU usage` = 0 命中，`git log -S` 亦无；在 `ScopeSentry-Scan` |
| 306 | Ignore 目标仍被 httpx/目录/漏洞扫描请求 | 服务端只下发 `scanTemplate.Ignore`（`common.go:254`），派生资产过滤在执行端 |
| 298 | httpx 跟随跳转时记录跳转后 URL | httpx 调用与结果解析在节点插件 |
| 196 | scan 容器 `apps` 目录大量 core 文件占盘 | 容器内进程崩溃产物，属节点镜像/运行时 |
| 219 | 8H12G 跑默认任务内存打满 | 作者归因 fasthttp/katana 结果 JSON 解析（#254 评论），节点侧 |
| 177 | POC 批量导入/删除后 scan 端不同步 | 服务端已发 `RefreshConfig`（`services/poc/poc.go:132,175,203,335`），消费逻辑在节点 |
| 197 | 资产 `time` 字段为空（`lastScanTime` 有值） | 服务端不写资产集合（仅 initdb 写系统集合），资产由节点直写 MongoDB |

## 信息不足(标题只有一句话)

| # | 标题 | 缺什么 |
| --- | --- | --- |
| 341 | 插件市场问题 | 仅一张截图，无报错文本/复现步骤 |
| 290 | 查看插件日志报错 | 仅截图；需 `logs:plugins:<module>:<hash>`（`services/plugin/plugin.go:259-271`）对应的 module/hash 与错误 |
| 289 | 敏感信息规则 pagesize 限制了大小 | 无正文；需说明是规则列表还是结果列表、期望上限 |
| 244 | 参数收集封装函数 | 无正文；需说明是要重构接口还是新功能 |
| 153 | 目录扫描结构优化 | 无正文；需给出期望的目录结构与收益 |

## 重点问题深度分析

### #342 nuclei parameter 被忽略 + 模板目录写死 /apps/poc

**结论：问题 2 已在本仓库修复；问题 1、3 的根因在 `Autumn-27/ScopeSentry-Scan`，不在本仓库。**

- 服务端只做“透传 + 追加 `-t`”：`internal/services/task/common/common.go:225-260` `GetScanTemplate` 在 `VulList` 非空时把 `-t <poc>.yaml,...` 追加到 `Parameters.VulnerabilityScan[hash]`，其余参数原样下发；`ParameterParser`（`common.go:262-339`）只替换 `{dict.*}`/`{port.*}`，不会吞掉 `-s`。
- 默认值本身就是 `-s high,critical`：`internal/constants/defaults.go:269-278`（Help 明确写“支持 t, s, es, tags, etags, rl…”，插件 id `ed93b8af6b72fe54a60efdb932cf6fbc`）。因此“info/medium 模板仍被执行”只能发生在执行端解析参数时。
- 全仓 grep 证据：`nuclei` 仅命中常量表、`plugin-template/nuclei_monitor.go` 示例插件与 `-t` 拼接处；`/apps/poc` **0 命中**（`grep -rn "/apps/poc" .`）。#280 的堆栈直接给出 `ScopeSentry-Scan/modules/vulnerabilityscan/nuclei/nuclei.go:270`，即 `(*Plugin).Execute` 里的 `No templates available`。
- 所以问题 1（`-s` 不生效）与问题 3（模板目录写死 `/apps/poc`）要在 scan 节点仓库改：解析 `parameter` 时把 `-s/-tags/-etags/-rl/-c` 映射到内嵌引擎 options，并把 `-t` 的 base dir 由常量改为可配置。
- 本仓库仍有一个可修点（藏在同一 issue 里）：`targetSource=UrlScan` + `targetTp=select` 时，`common.go:111` 的分支不判 `TargetTp`，会走 search 路径且 `PageSize=0`（`Query.PageSize=task.TargetNumber`，0 表示不限，见 `services/assets/url/url.go:112-129`）——即便 panic 已修，也会把整库 URL 当目标。建议补 `UrlScanSource && TargetTp=="select"` 分支的等价判断。
- 修复顺序建议：节点侧先支持 `-severity` 透传（收益最大、改动最小），再谈 `templatesDir`/`pocProfile` 字段。

### #285 redis 相关问题（节点多了分发失败）

- 任务分发是“服务端推、节点被动消费”：`CreateTaskScan` 遍历 `task.Node` 逐节点 `RPush`（`internal/services/task/common/common.go:213-220` → `internal/repositories/task/task/task.go:179-190`，key `NodeTask:<node>`）。
- `AllNode` 只在**建任务那一刻**展开节点列表（`common.go:178-190` `nodeService.GetNodeData`），之后加入的节点拿不到该任务；节点离线期间消息在 list 里无限堆积，恢复后一次性消费 → 与 #297（旧模板）、#254（越重启越乱）同源。
- 日志侧：`GetLogs` 用 `logs:plugins:<module>:<hash>` 的 Redis list（`internal/services/plugin/plugin.go:259-271`），无长度上限；作者也在 #285 承认要“限制插件日志存储条数”。
- 建议：节点改为主动拉取（服务端暴露按节点名 claim/pop 的接口，或共享 Stream + 消费组），日志改 `XADD ... MAXLEN ~ N`；节点侧任务列表按 target 粒度分片，避免单节点堆积全量任务。

### #254 分布式节点总是离线 + 资源打满

- 离线判定阈值写死 300s：`internal/services/node/service.go:24-26` `NodeTimeout = 300`。
- 判定方式是“节点自报 `UpdateTime` vs 服务器 `time.Now()`”：`service.go:59-83`，差值 >300s 就异步 `updateNodeState(...,"3")`（`service.go:102-110`）。日志里 `Node time difference {"time_diff": 9.12...}` 说明这条链路对时钟/时区敏感；容器化部署时节点与服务器时区不一致（`config.GlobalConfig.System.Timezone`）会直接误判。
- 一旦置 3，只有节点下次心跳写回 `state=1` 才恢复，前端表现为“总是离线、重启好几分钟又掉”。
- 内存侧不是本仓库问题：作者归因 fasthttp/katana 结果 JSON 解析（#254 评论），并建议 `sentrydir` 不选只留 spray——属 scan 节点。
- 建议：心跳改 Redis TTL（`SET node:<name> EX 90`，服务端读 TTL 判活），容差可配；state=3 增加“疑似离线/恢复中”过渡态；把 `time_diff` 从 Info 降噪或加阈值告警。

### #288 / #297 定时任务状态与“本地缓存的旧模板”

- #288 根因有两层：任务表的状态枚举只有 1/2/3（`frontend/src/views/Task/Task.vue:90-114`），而插入任务时无条件 `task.Status = 1`（`internal/services/task/common/common.go:148-151`）；计划任务每轮生成的子任务在首次执行前就显示“运行中”。模型层也没有“待开始”位：`internal/models/task.go:102-127` 把 `scheduledTasks` 直接映射成 `state`（启用位）。计划任务列表页（`ScheduledTask.vue:46-120`）只显示启用/停用与 `lastTime/nextTime`，两边口径不一致。
- #288 另有产品诉求：作者已承诺“计划任务页加一个执行一次的选项”（#288 评论），属新功能，别混进 bug 修。
- #297 的“本地缓存任务用旧模板”不是节点 bug，而是服务端**入队即快照**：`GetScanTemplate` 把解析后的模板（含 `-t`、字典/端口占位符）序列化后 `RPush`（`common.go:213-220`），而 `ClearTaskCache` 只清 `TaskInfo:tmp/time/progress:*`（`repositories/task/task/task.go:107-120`），从不触碰 `NodeTask:*`。
- 因此：模板改完→队列里未消费的条目仍是旧参数；节点重启/离线恢复后消费到旧 payload，复现“用的还是以前的模板”。
- 建议修法（按成本从低到高）：模板保存时按 taskId 清理未消费的 `NodeTask:*` 条目；payload 增加 `templateVersion` 并在消费时回源比对；或把模板解析从入队时机挪到节点 claim 之后（代价是节点要回源，需要接口）。

## 建议下一轮处理顺序

| Issue | 理由（影响/成本） |
| --- | --- |
| #295 | 未鉴权即可删模板，纯安全洞，一行中间件（S） |
| #340 | 一键更新必现重复，所有实例中招（S） |
| #317 | “同步到项目”整体 404，核心工作流断（M） |
| #225 | csv/xlsx 文案与默认值不一致，用户可见（S） |
| #262 | 被动扫描未完成即显示完成，误判扫描结束（S） |
| #288 | 计划任务状态误导，前端映射 + 一个状态位（S/M） |
| #228 | ICP/换行前缀导致项目匹配失败，且是 #242 的前置（S） |
| #242 | 同步会误清项目资产，影响大（M） |
| #297 | 旧模板复扫，修完顺带缓解 #285 的堆积（M） |
| #254 | 节点误判离线，容器/时区场景高频（M） |
| #232 | 巨型正则，直接卡大项目更新（L） |
| #169 | 漏洞缺 lastScanTime，先补字段再谈复测（S） |

> #342 的问题 1/3 不在本仓库（见深度分析），下一轮以“给 scan 节点开 issue/PR”方式推进，不计入本表。

## 复核命令

```bash
gh issue list -R Autumn-27/ScopeSentry -L 100 --state open --json number,title
gh api repos/Autumn-27/ScopeSentry/issues/342/comments
grep -rn "/apps/poc" . ; grep -rn "SlogErrorLocal" .
git show 4626678 --stat   # #342 panic 修复（PR #343）
git show db0fd09 --stat   # #327 修复（PR #328）
grep -n "Path:" internal/api/routes/task/task.go        # 确认无 /sync
sed -n '90,114p' frontend/src/views/Task/Task.vue        # 状态枚举缺“待开始”
```
