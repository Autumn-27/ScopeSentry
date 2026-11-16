// Package audit 只包含针对路由注册表的静态契约测试。
//
// 单独成包的原因：路由包及其 handlers 会 import bootstrap / mongodb，
// 包 init 阶段就必须连上 MongoDB，这里只做 AST 解析，不引入任何后端依赖，
// 保证测试可以在没有数据库的环境下运行。
package audit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// anonymousRoutes 是允许匿名访问的路由白名单，key 为 "<相对 internal/api/routes 的路径>|<METHOD>|<Path>"。
// 只有登录接口必须匿名可访问；任何新增的公开路由都必须在这里登记并写明原因。
// 见 issue #295：POST /api/task/template/delete 曾经漏挂认证中间件，匿名即可删除扫描模板。
var anonymousRoutes = map[string]string{
	"user/user.go|POST|/login": "登录接口必须匿名可访问",
}

// TestEveryRouteRequiresAuth 静态扫描 internal/api/routes 下所有路由注册，
// 断言除白名单外每条路由都挂了认证中间件（WithAuth / WithAuthAndPermission）。
func TestEveryRouteRequiresAuth(t *testing.T) {
	seen := map[string]bool{}
	violations := make([]string, 0)

	root := ".."
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		rel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(path, root), "/"))
		ast.Inspect(file, func(n ast.Node) bool {
			kv, ok := n.(*ast.KeyValueExpr)
			if !ok {
				return true
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Routes" {
				return true
			}
			list, ok := kv.Value.(*ast.CompositeLit)
			if !ok {
				return true
			}
			for _, elt := range list.Elts {
				route, ok := elt.(*ast.CompositeLit)
				if !ok {
					continue
				}
				method, routePath, hasAuth := parseRouteLiteral(route)
				if method == "" || routePath == "" {
					continue
				}
				id := rel + "|" + method + "|" + routePath
				seen[id] = true
				if hasAuth {
					continue
				}
				if _, allowed := anonymousRoutes[id]; allowed {
					continue
				}
				violations = append(violations, id)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk route files: %v", err)
	}

	for _, v := range violations {
		t.Errorf("route %s is registered without auth middleware", v)
	}
	if len(seen) == 0 {
		t.Fatal("no routes were parsed; the audit walk is broken")
	}
	for id, reason := range anonymousRoutes {
		if !seen[id] {
			t.Errorf("anonymous route %q (%s) no longer exists; update the allowlist", id, reason)
		}
	}
}

// parseRouteLiteral 从 models.Route{...} 字面量里取出 method、path 以及是否挂了认证中间件。
func parseRouteLiteral(route *ast.CompositeLit) (method, path string, hasAuth bool) {
	for _, elt := range route.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		switch key.Name {
		case "Method":
			method = stringLiteral(kv.Value)
		case "Path":
			path = stringLiteral(kv.Value)
		case "Middlewares":
			hasAuth = containsAuthMiddleware(kv.Value)
		}
	}
	return method, path, hasAuth
}

func stringLiteral(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return ""
	}
	return value
}

// containsAuthMiddleware 判断中间件表达式里是否调用了 WithAuth / WithAuthAndPermission。
func containsAuthMiddleware(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok && (sel.Sel.Name == "WithAuth" || sel.Sel.Name == "WithAuthAndPermission") {
			found = true
			return false
		}
		return true
	})
	return found
}
