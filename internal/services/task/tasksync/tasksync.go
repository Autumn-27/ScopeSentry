// Package tasksync 处理"同步任务目标到项目"的纯逻辑：参数校验、目标去重。
//
// 单独成包以便在不连接 Mongo/Redis 的情况下做单元测试（任务服务会通过
// assets 服务间接 import bootstrap，init 阶段就要求数据库可用）。
package tasksync

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/models"
)

// 同步方式（与前端 /api/task/sync 的 option 字段对应）
const (
	OptionExisting = "existing" // 并入已有项目
	OptionNew      = "new"      // 新建项目
)

// Request 是 POST /api/task/sync 的请求体，字段名与前端 syancProjectApi 保持一致。
type Request struct {
	IDS     []string `json:"ids"`     // 要同步的任务 ID
	Option  string   `json:"option"`  // existing / new
	Project string   `json:"project"` // option=existing 时的项目 ID
	Tag     string   `json:"tag"`     // option=new 时的项目标签
	Name    string   `json:"name"`    // option=new 时的项目名称
}

// Validate 校验同步参数，返回的错误会以 400 返回给前端。
func (r Request) Validate() error {
	if len(r.IDS) == 0 {
		return errors.New("请选择要同步的任务")
	}
	switch strings.TrimSpace(r.Option) {
	case OptionExisting:
		if strings.TrimSpace(r.Project) == "" {
			return errors.New("请选择要同步到的项目")
		}
	case OptionNew:
		if strings.TrimSpace(r.Name) == "" {
			return errors.New("新建项目名称不能为空")
		}
	default:
		return fmt.Errorf("不支持的同步方式: %q", r.Option)
	}
	return nil
}

// ExtractTargets 汇总所选任务的目标：按行切分、去掉首尾空白、按出现顺序去重。
// 目标的换行/去重规则与扫描任务创建时（helper.GetTargetList）保持一致。
// 全部任务都没有目标时返回明确错误，避免把空目标写进项目。
func ExtractTargets(tasks []models.Task) ([]string, error) {
	targets := collect(nil, taskTargets(tasks))
	if len(targets) == 0 {
		return nil, errors.New("所选任务没有可同步的目标")
	}
	return targets, nil
}

func taskTargets(tasks []models.Task) []string {
	lines := make([]string, 0, len(tasks))
	for _, t := range tasks {
		lines = append(lines, strings.Split(t.Target, "\n")...)
	}
	return lines
}

// MergeExisting 把新增目标并入项目已有目标（都按行解析），返回去重后的目标列表。
// 已有目标在前，新增目标按顺序追加，重复项只保留第一次出现的。
func MergeExisting(existing string, additions []string) []string {
	targets := collect(nil, append(strings.Split(existing, "\n"), additions...))
	return targets
}

// collect 去空白、去重，保持首次出现顺序。
func collect(seen map[string]struct{}, raw []string) []string {
	if seen == nil {
		seen = make(map[string]struct{}, len(raw))
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		target := strings.TrimSpace(item)
		if target == "" {
			continue
		}
		if _, ok := seen[target]; ok {
			continue
		}
		seen[target] = struct{}{}
		result = append(result, target)
	}
	return result
}

// CountNew 返回 additions 中真正新增（不在 existing 里）的目标数量。
func CountNew(existing string, additions []string) int {
	before := len(collect(nil, strings.Split(existing, "\n")))
	after := len(MergeExisting(existing, additions))
	if after <= before {
		return 0
	}
	return after - before
}
