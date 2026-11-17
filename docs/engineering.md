# 工程化约定

## 本地命令（Makefile）

```bash
make help            # 列出目标
make build           # go build ./...
make binary          # 产出 dist/ScopeSentry
make test-hermetic   # 不需要 Mongo/Redis 的单元测试
make vet             # 全仓 go vet（CI 同款阻塞闸）
make fmt             # gofmt -w .
make fmt-check       # CI 同款：有未格式化文件就失败
make mod-check       # go mod tidy 后检查 go.mod/go.sum 是否漂移
make frontend-test   # 前端单元测试（vitest）
make frontend-check  # 语言包一致性 + 主题三态自检 + 单测
make ci              # fmt-check + mod-check + test-hermetic + frontend-check
```

前端另有：

```bash
cd frontend
pnpm install --no-frozen-lockfile
pnpm build:pro       # dist-pro
node scripts/check-locales.mjs
npx esno ./scripts/check-theme.ts
```

## CI

- 运行在自建 runner `dev-tx-sh-app-01-ss`（标签 `self-hosted,linux,x64,scopesentry-ci`，systemd 服务）。
- 触发：PR / push main / 手动 `workflow_dispatch`；同一 ref 的旧运行会被取消（concurrency）。
- 门禁：
  - **gofmt**（阻塞）
  - **go mod tidy 漂移**（阻塞）
  - `go build ./...`（阻塞）
  - hermetic 单测（阻塞）
  - 全仓 `go vet ./...`（阻塞；历史 144 条告警已清零）
  - 前端 `build:pro` + 语言包 + 主题自检（阻塞）
  - `vue-tsc`（非阻塞，历史类型错误较多）
  - **Integration job**（阻塞）：临时起 Mongo + Redis（随机端口、run_id 命名、跑完清理），
    先跑一个包把空库初始化一次，再 `go test ./...`；CI 上等价于 `make test-integration`
- Go 单测产出 `coverage.out` 并上传 artifact（`go-coverage`）。
- Go 模块走 `goproxy.cn`（这台机器到 proxy.golang.org 的 IPv6 无路由）；`setup-go` 的 Actions 缓存关闭，模块缓存在 runner 的 `~/go/pkg/mod` 常驻。

## 节点存活判定

- 节点心跳写在 Redis `node:<name>` 哈希里（`state` / `updateTime`），服务端只读。
- 判定用 `ClassifyLiveness`：同一时间字符串按「服务器时区」和「UTC」各解释一次，取与服务器时间最接近者；
  因此节点与服务器时区不同（例如节点跑 UTC）不会再被误判离线，节点时钟略快也只标记 ClockSkew。
- 超时窗口默认 300s，可用 `NODE_TIMEOUT_SECONDS` 覆盖（节点多、网络慢时可放宽）。
- 时间字符串解析失败时保留节点并打 Warn，而不是把它从列表里抹掉。

## 探活

- `GET /healthz`：ping Mongo + Redis，正常 200 `{status:"ok",version,time,checks:{mongo,redis}}`；任一依赖不可用返回 **503**。
- K8s 的 readiness / liveness 已指向该端点（`k8s/pods/scope-sentry-deployment.yml`）。

## 测试分层

| 层 | 范围 | 运行方式 |
| --- | --- | --- |
| hermetic | 纯逻辑/中间件，不连库 | `make test-hermetic` |
| integration | 需要 Mongo/Redis（包 init 就会连库） | 本地起容器后 `go test ./...` |
| frontend | 语言包/主题自检 + 单元测试（vitest + jsdom） | `make frontend-check` |

新增测试优先放叶子包（不 import `internal/bootstrap`），否则包级 init 会直接连库、CI 无法运行。

## 本地 Git hook（可选）

```bash
git config core.hooksPath .githooks   # 启用 pre-commit：暂存的 .go 文件必须已 gofmt
```

## 仓库规范

- `.editorconfig`：Go 用 tab，其余 2 空格；LF。
- `.gitattributes`：统一 LF，图片/字体标记为 binary。
- PR 模板在 `.github/pull_request_template.md`，要求写清验证项。
- Dependabot 每周检查 gomod / npm / GitHub Actions。

## 分支与发布

- 功能分支 → PR → main；CI 绿后合并。
- 版本由 `.goreleaser.yaml` 在 tag 时注入 `internal/constants.Version`。
