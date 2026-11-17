package helper

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func assetKeyword() map[string]string {
	// 复制一份，避免测试之间互相污染（SearchToMongoDB 会往 map 里写 task/rootDomain/tag）
	kw := make(map[string]string, len(SearchKey["asset"]))
	for k, v := range SearchKey["asset"] {
		kw[k] = v
	}
	return kw
}

func asJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// 裸词（直接输域名/标题）必须展开成跨身份字段的 $or，而不是被静默丢弃
func TestBareTermExpandsToIdentityFields(t *testing.T) {
	q, err := SearchToMongoDB("example.com", assetKeyword())
	if err != nil {
		t.Fatalf("bare term must not error, got %v", err)
	}
	if len(q) != 1 {
		t.Fatalf("want 1 clause, got %d", len(q))
	}
	out := asJSON(q[0])
	if !strings.Contains(out, "\"$or\"") {
		t.Fatalf("bare term should expand to $or, got %s", out)
	}
	for _, want := range []string{"\"host\"", "\"title\"", "\"ip\""} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in bare-term OR, got %s", want, out)
		}
	}
	// 正文类大字段不能参与，否则一个裸词会触发全集合大文本扫描
	for _, ban := range []string{"\"body\"", "\"rawheaders\"", "\"metadata\""} {
		if strings.Contains(out, ban) {
			t.Fatalf("bare term must not scan %s, got %s", ban, out)
		}
	}
}

func TestBareTermEscapesRegexMeta(t *testing.T) {
	q, err := SearchToMongoDB("a.*b", assetKeyword())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(asJSON(q[0]), "a.*b") {
		t.Fatalf("regex metacharacters must be escaped, got %s", asJSON(q[0]))
	}
}

func TestBareTermWithAnd(t *testing.T) {
	q, err := SearchToMongoDB("example.com && ip=1.1.1.1", assetKeyword())
	if err != nil {
		t.Fatalf("bare term combined with && must not error, got %v", err)
	}
	out := asJSON(q[0])
	if !strings.Contains(out, "\"$and\"") || !strings.Contains(out, "\"$or\"") {
		t.Fatalf("expected $and($or(bare), ip), got %s", out)
	}
}

func TestOrOperatorWorks(t *testing.T) {
	q, err := SearchToMongoDB("domain=a.com || ip=1.1.1.1", assetKeyword())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(asJSON(q[0]), "\"$or\"") {
		t.Fatalf("expected $or, got %s", asJSON(q[0]))
	}
}

func TestDanglingOperatorIsRejected(t *testing.T) {
	if _, err := SearchToMongoDB("domain=a.com &&", assetKeyword()); err == nil {
		t.Fatal("expected error for dangling &&")
	}
}

func TestKnownFieldStillWorks(t *testing.T) {
	q, err := SearchToMongoDB("domain=example.com", assetKeyword())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := asJSON(q[0])
	if !strings.Contains(out, "\"host\"") || !strings.Contains(out, "example.com") {
		t.Fatalf("expected host regex, got %s", out)
	}
}

func TestEmptyExpressionReturnsAll(t *testing.T) {
	q, err := SearchToMongoDB("", assetKeyword())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q) != 1 || len(q[0]) != 0 {
		t.Fatalf("empty expression should yield one empty clause, got %s", asJSON(q))
	}
}

func TestGetSearchQueryRejectsUnknownIndex(t *testing.T) {
	if _, err := GetSearchQuery(models.SearchRequest{Index: "nope"}); err == nil {
		t.Fatal("expected error for unknown index")
	}
}
