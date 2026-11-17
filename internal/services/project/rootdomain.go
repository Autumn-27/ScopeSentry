package project

import (
	"strings"
	"unicode"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
)

// 目标前缀（区分大小写地按大写匹配，其余部分保留原样）
const (
	prefixICP   = "ICP:"
	prefixCMP   = "CMP:"
	prefixAPP   = "APP:"
	prefixAPPID = "APP-ID:"
)

// isICPText 判断一段文本是否是 ICP 备案号（可能没有 ICP: 前缀）。
// 只认带「备/证」的中文备案号，避免把 icp.example.com 这类主机名误判成备案号。
func isICPText(s string) bool {
	upper := strings.ToUpper(s)
	if strings.HasPrefix(upper, prefixICP) {
		return true
	}
	if !strings.Contains(upper, "ICP") {
		return false
	}
	return strings.Contains(s, "备") || strings.Contains(s, "证")
}

// stripICPVariant 去掉 ICP 号末尾的站点序号后缀（-1 / -2 ...），
// 同一份备案号在资产里可能以不同后缀出现，归并后才能匹配上。
func stripICPVariant(s string) string {
	if i := strings.LastIndex(s, "-"); i > 0 {
		suffix := s[i+1:]
		if suffix != "" && strings.Trim(suffix, "0123456789") == "" {
			return s[:i]
		}
	}
	return s
}

// normalizeICP 把任意形式的备案号规范成 ICP:<备案号>
func normalizeICP(raw string) string {
	body := strings.TrimSpace(raw)
	body = strings.TrimPrefix(body, prefixICP)
	body = strings.TrimSpace(body)
	body = stripICPVariant(body)
	if body == "" {
		return ""
	}
	return prefixICP + body
}

// looksLikeHostOrIP 判断字符串是否已经像一个主机/IP（而不是公司名）
func looksLikeHostOrIP(s string) bool {
	if strings.Contains(s, ".") || strings.Contains(s, ":") || strings.Contains(s, "/") {
		return true
	}
	// 纯数字（少见）或单个 label：交给上层当公司名处理更安全
	return false
}

// hasCJKOrSpace 判断是否含中日韩字符或空白：这类几乎一定是公司名而不是主机
func hasCJKOrSpace(s string) bool {
	if strings.ContainsAny(s, " \t（）()【】[]，,") {
		return true
	}
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hangul, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) {
			return true
		}
	}
	return false
}

// BuildRootDomains 把目标列表归并成项目 rootDomain 列表。
//
// 修 #228：目标里混进来的空行、以及没有前缀的公司名/备案号以前会原样落库
// （公司名缺 CMP: / 备案号缺 ICP:），导致项目与资产的关联匹配不上。
// 现在统一补前缀、规范化备案号、去掉空值并保序去重。
func BuildRootDomains(targetList []string) []string {
	out := make([]string, 0, len(targetList))
	seen := make(map[string]struct{}, len(targetList))

	add := func(v string) {
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}

	for _, raw := range targetList {
		tg := strings.TrimSpace(raw)
		if tg == "" {
			continue
		}
		upper := strings.ToUpper(tg)

		switch {
		case strings.HasPrefix(upper, prefixICP) || isICPText(tg):
			add(normalizeICP(tg))
		case strings.HasPrefix(upper, prefixCMP),
			strings.HasPrefix(upper, prefixAPP),
			strings.HasPrefix(upper, prefixAPPID):
			// 已带前缀，去掉多余空白后原样保留
			add(strings.ReplaceAll(tg, " ", ""))
		default:
			root, err := helper.GetRootDomain(tg)
			if err == nil && root != "" && looksLikeHostOrIP(root) {
				add(root)
				continue
			}
			// 解析不出根域名：含中文/空白的按公司名补 CMP:（#228）；
			// 单 label 的 ASCII（intranet、localhost）与通配符保留原样，避免把主机名弄坏
			if hasCJKOrSpace(tg) {
				add(prefixCMP + tg)
			} else {
				add(tg)
			}
		}
	}
	return out
}
