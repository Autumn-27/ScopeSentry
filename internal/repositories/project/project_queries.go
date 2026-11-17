package project

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// maxICPPerRegex 单个正则里最多塞多少个备案号，避免拼出超长正则（#232）
const maxICPPerRegex = 50

// icpValues 取出需要按前缀匹配的备案号（去掉 ICP: 前缀、保序去重）
func icpValues(rootDomains []string) []string {
	out := make([]string, 0, len(rootDomains))
	seen := make(map[string]struct{}, len(rootDomains))
	for _, d := range rootDomains {
		v := strings.TrimSpace(d)
		if len(v) < len("ICP:") || !strings.EqualFold(v[:len("ICP:")], "ICP:") {
			continue
		}
		v = strings.TrimSpace(v[len("ICP:"):])
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// icpMatchClauses 生成 {icp: {$regex: ^(?:a|b|c)}} 子句，按 maxICPPerRegex 分片。
// 没有备案号时返回 nil —— 大多数项目不会走正则匹配，也就不会拼出巨串。
func icpMatchClauses(rootDomains []string) []bson.M {
	values := icpValues(rootDomains)
	if len(values) == 0 {
		return nil
	}
	clauses := make([]bson.M, 0, (len(values)+maxICPPerRegex-1)/maxICPPerRegex)
	for i := 0; i < len(values); i += maxICPPerRegex {
		end := i + maxICPPerRegex
		if end > len(values) {
			end = len(values)
		}
		parts := make([]string, 0, end-i)
		for _, v := range values[i:end] {
			parts = append(parts, "^"+regexp.QuoteMeta(v))
		}
		clauses = append(clauses, bson.M{
			"icp": bson.M{"$regex": "(?:" + strings.Join(parts, "|") + ")", "$options": "i"},
		})
	}
	return clauses
}

// buildRootDomainQuery 构建 RootDomain 集合的查询条件。
//
// isUpdate 分支只清理「本项目」的关联（#242：以前用 project 字段存在就清，
// 会把其它项目已经关联好的资产一起清掉）。
// 匹配分支用 $in 精确匹配 + 分片的 ICP 前缀正则（#232）。
func (r *repository) buildRootDomainQuery(rootDomains []string, projectID string, isUpdate bool) bson.M {
	icpClauses := icpMatchClauses(rootDomains)
	if isUpdate {
		and := []bson.M{
			{"project": projectID},
			{"domain": bson.M{"$nin": rootDomains}},
			{"company": bson.M{"$nin": rootDomains}},
		}
		if len(icpClauses) > 0 {
			and = append(and, bson.M{"$nor": icpClauses})
		}
		return bson.M{"$and": and}
	}
	or := []bson.M{
		{"domain": bson.M{"$in": rootDomains}},
		{"company": bson.M{"$in": rootDomains}},
	}
	or = append(or, icpClauses...)
	return bson.M{"$or": or}
}

// buildAppQuery 构建 app 集合的查询条件（字段：name/bundleID/company/icp）
func (r *repository) buildAppQuery(rootDomains []string, projectID string, isUpdate bool) bson.M {
	icpClauses := icpMatchClauses(rootDomains)
	if isUpdate {
		and := []bson.M{
			{"project": projectID},
			{"company": bson.M{"$nin": rootDomains}},
		}
		if len(icpClauses) > 0 {
			and = append(and, bson.M{"$nor": icpClauses})
		}
		return bson.M{"$and": and}
	}
	or := []bson.M{
		{"name": bson.M{"$in": rootDomains}},
		{"bundleID": bson.M{"$in": rootDomains}},
		{"company": bson.M{"$in": rootDomains}},
	}
	or = append(or, icpClauses...)
	return bson.M{"$or": or}
}

// buildMpQuery 构建 mp 集合的查询条件（字段：name/company/bundleID/icp）
func (r *repository) buildMpQuery(rootDomains []string, projectID string, isUpdate bool) bson.M {
	icpClauses := icpMatchClauses(rootDomains)
	if isUpdate {
		and := []bson.M{
			{"project": projectID},
			{"name": bson.M{"$nin": rootDomains}},
			{"company": bson.M{"$nin": rootDomains}},
			{"bundleID": bson.M{"$nin": rootDomains}},
		}
		if len(icpClauses) > 0 {
			and = append(and, bson.M{"$nor": icpClauses})
		}
		return bson.M{"$and": and}
	}
	or := []bson.M{
		{"company": bson.M{"$in": rootDomains}},
	}
	or = append(or, icpClauses...)
	return bson.M{"$or": or}
}
