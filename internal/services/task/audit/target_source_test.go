// Package audit 存放任务目标来源解析的静态契约测试。
//
// 单独成包的原因：internal/services/task/common 会通过 assets 服务 import bootstrap，
// 包 init 阶段就要求 Mongo/Redis 可用；这里只解析源码 AST，不引入后端依赖。
package audit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// TestSourceBranchesCheckTargetTp 断言 common.Insert 里每个 "xxxSource" 分支都同时判断了
// TargetTp，且 search / select 两种模式都存在。
//
// 背景：targetSource 有两种写法——裸名（asset/RootDomain/subdomain/UrlScan，等价于
// "从资产库搜索"）和 "<name>Source" + targetTp（search 搜索 / select 选中 ID）。
// 如果某个 *Source 分支漏判 targetTp，targetTp=select 的请求会掉进搜索分支，
// 变成对全库发查询（issue 分诊记录）。
func TestSourceBranchesCheckTargetTp(t *testing.T) {
	const path = "../common/common.go"
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	sourceModes := map[string]map[string]bool{} // source -> {"search":true,"select":true}
	ast.Inspect(file, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}
		sources, tps := comparedStrings(ifStmt.Cond, "TargetSource", "TargetTp")
		for _, source := range sources {
			if !strings.HasSuffix(source, "Source") {
				continue
			}
			if len(tps) == 0 {
				t.Errorf("branch on TargetSource %q does not check TargetTp", source)
				continue
			}
			if sourceModes[source] == nil {
				sourceModes[source] = map[string]bool{}
			}
			for _, tp := range tps {
				sourceModes[source][tp] = true
			}
		}
		return true
	})

	for _, source := range []string{"assetSource", "RootDomainSource", "subdomainSource", "UrlScanSource"} {
		modes := sourceModes[source]
		if modes == nil {
			t.Errorf("no branch found for TargetSource %q", source)
			continue
		}
		for _, tp := range []string{"search", "select"} {
			if !modes[tp] {
				t.Errorf("TargetSource %q has no targetTp=%q branch", source, tp)
			}
		}
	}
}

// comparedStrings 收集 cond 中以 == 比较 field 得到的字符串常量，以及以 == 比较 otherField 的常量。
func comparedStrings(cond ast.Expr, field, otherField string) (values, others []string) {
	ast.Inspect(cond, func(n ast.Node) bool {
		bin, ok := n.(*ast.BinaryExpr)
		if !ok || bin.Op != token.EQL {
			return true
		}
		if name, value, ok := selectorEqualsString(bin, field); ok && name != "" {
			values = append(values, value)
		}
		if _, value, ok := selectorEqualsString(bin, otherField); ok {
			others = append(others, value)
		}
		return true
	})
	return values, others
}

// selectorEqualsString 判断表达式是否是 x.Field == "literal"（或反过来）。
func selectorEqualsString(bin *ast.BinaryExpr, field string) (string, string, bool) {
	check := func(left, right ast.Expr) (string, string, bool) {
		sel, ok := left.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != field {
			return "", "", false
		}
		lit, ok := right.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return "", "", false
		}
		value := strings.Trim(lit.Value, "\"")
		ident := ""
		if id, ok := sel.X.(*ast.Ident); ok {
			ident = id.Name
		}
		return ident, value, true
	}
	if a, b, ok := check(bin.X, bin.Y); ok {
		return a, b, true
	}
	return check(bin.Y, bin.X)
}
