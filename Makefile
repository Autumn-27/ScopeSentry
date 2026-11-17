# ScopeSentry 开发者入口
# 用法：make help
SHELL := /bin/bash
GO ?= go
PNPM ?= pnpm
BINARY ?= dist/ScopeSentry

# 不依赖 Mongo/Redis 的包（其余包在 init 阶段就会连库）
HERMETIC_PKGS := \
	./internal/services/task/tasksync/ \
	./internal/services/task/taskstate/ \
	./internal/services/task/audit/ \
	./internal/api/routes/audit/ \
	./internal/repositories/plugin/ \
	./internal/services/project/ \
	./internal/api/handlers/sensitiveRule/ \
	./internal/api/middleware/compress/

.PHONY: help build binary run test-hermetic test-integration vet fmt fmt-check mod-check \
	frontend-install frontend-build frontend-test frontend-check ci clean

help: ## 列出所有目标
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

build: ## 编译全部 Go 包（不产出二进制）
	$(GO) build ./...

binary: ## 产出服务器二进制到 dist/ScopeSentry
	mkdir -p dist
	$(GO) build -o $(BINARY) ./cmd/main

run: ## 本地启动服务（需先配置 Mongo/Redis 环境变量）
	$(GO) run ./cmd/main

test-hermetic: ## 跑不依赖 Mongo/Redis 的单元测试
	$(GO) test -count=1 $(HERMETIC_PKGS)

test-integration: ## 全量测试（需要可用的 Mongo/Redis，读 MONGODB_*/REDIS_* 环境变量）
	$(GO) test -tags integration -count=1 ./...

vet: ## 对本次维护的包做 go vet
	$(GO) vet ./internal/api/handlers/task/ ./internal/api/routes/task/ \
		./internal/services/task/... ./internal/repositories/plugin/

fmt: ## 格式化所有 Go 代码
	gofmt -w .

fmt-check: ## 检查是否有未格式化的 Go 文件
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then \
		echo "以下文件未格式化（运行 make fmt）:"; echo "$$out"; exit 1; fi

mod-check: ## 检查 go.mod / go.sum 是否与代码一致
	$(GO) mod tidy
	@git diff --exit-code go.mod go.sum || (echo "go.mod/go.sum 有漂移，请提交 tidy 结果"; exit 1)

frontend-install: ## 安装前端依赖
	cd frontend && $(PNPM) install --no-frozen-lockfile

frontend-build: ## 前端生产构建（dist-pro）
	cd frontend && $(PNPM) build:pro

frontend-test: ## 前端单元测试（vitest）
	cd frontend && $(PNPM) test

frontend-check: ## 前端轻量自检：语言包一致性 + 主题三态 + 单测
	cd frontend && node scripts/check-locales.mjs && npx esno ./scripts/check-theme.ts
	cd frontend && $(PNPM) test

ci: fmt-check mod-check test-hermetic frontend-check ## 本地复现 CI 的核心门禁

clean: ## 清理构建产物
	rm -rf dist frontend/dist-pro
