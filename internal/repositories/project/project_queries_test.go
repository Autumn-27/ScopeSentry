package project

import (
	"encoding/json"
	"strings"
	"testing"
)

func queryJSON(q interface{}) string {
	b, _ := json.Marshal(q)
	return string(b)
}

// #242：清理步骤必须只作用于本项目，否则会把其它项目的资产关联清空
func TestClearQueryIsScopedToProject(t *testing.T) {
	r := &repository{}
	q := r.buildRootDomainQuery([]string{"example.com"}, "PROJ_A", true)
	out := queryJSON(q)
	if !strings.Contains(out, "PROJ_A") {
		t.Fatalf("clear query must restrict to the project, got %s", out)
	}
	if strings.Contains(out, "$exists") {
		t.Fatalf("clear query must not use a bare $exists on project, got %s", out)
	}
}

// #232：不再把所有域名拼进一条正则
func TestQueryKeepsDomainsExact(t *testing.T) {
	r := &repository{}
	domains := []string{"a.com", "b.com", "ICP:京ICP备1号"}
	q := r.buildRootDomainQuery(domains, "P", false)
	out := queryJSON(q)
	if !strings.Contains(out, "\"$in\"") {
		t.Fatalf("domains should be matched with $in, got %s", out)
	}
	if strings.Contains(out, "^a") || strings.Contains(out, "^b") {
		t.Fatalf("plain domains must not go into the ICP regex, got %s", out)
	}
	if !strings.Contains(out, "京ICP备1号") {
		t.Fatalf("ICP value should be matched by prefix regex, got %s", out)
	}
}

func TestICPOnlyProjectsProduceNoRegex(t *testing.T) {
	r := &repository{}
	q := r.buildRootDomainQuery([]string{"a.com"}, "P", false)
	if strings.Contains(queryJSON(q), "$regex") {
		t.Fatalf("no ICP entries should mean no regex clause, got %s", queryJSON(q))
	}
}

func TestICPMatchClausesAreChunked(t *testing.T) {
	domains := make([]string, 0, 130)
	for i := 0; i < 130; i++ {
		domains = append(domains, "ICP:京ICP备"+string(rune('A'+i%26))+"-"+itoa(i))
	}
	clauses := icpMatchClauses(domains)
	if len(clauses) != 3 { // 50 + 50 + 30
		t.Fatalf("expected 3 chunked clauses, got %d", len(clauses))
	}
	for _, c := range clauses {
		body := queryJSON(c)
		if len(body) > 4000 {
			t.Fatalf("chunk regex too long: %d bytes", len(body))
		}
	}
}

func TestICPMatchClausesEscapes(t *testing.T) {
	clauses := icpMatchClauses([]string{"ICP:a.b*c"})
	if len(clauses) != 1 {
		t.Fatalf("expected one clause, got %d", len(clauses))
	}
	out := queryJSON(clauses[0])
	if strings.Contains(out, "a.b*c") {
		t.Fatalf("metacharacters must be escaped, got %s", out)
	}
}

func TestICPMatchClausesDedupes(t *testing.T) {
	clauses := icpMatchClauses([]string{"ICP:same", "ICP:same"})
	if len(clauses) != 1 {
		t.Fatalf("expected one clause, got %d", len(clauses))
	}
	if n := strings.Count(queryJSON(clauses[0]), "same"); n != 1 {
		t.Fatalf("duplicate ICP values should be deduped, got %d occurrences", n)
	}
}

func TestAppAndMpQueriesScopedToProject(t *testing.T) {
	r := &repository{}
	for name, q := range map[string]interface{}{
		"app": r.buildAppQuery([]string{"a.com"}, "PROJ_B", true),
		"mp":  r.buildMpQuery([]string{"a.com"}, "PROJ_B", true),
	} {
		out := queryJSON(q)
		if !strings.Contains(out, "PROJ_B") || strings.Contains(out, "$exists") {
			t.Fatalf("%s clear query must be project-scoped, got %s", name, out)
		}
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
