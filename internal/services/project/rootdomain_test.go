package project

import (
	"reflect"
	"testing"
)

func TestBuildRootDomainsAddsCompanyPrefix(t *testing.T) {
	got := BuildRootDomains([]string{"腾讯科技（深圳）有限公司"})
	want := []string{"CMP:腾讯科技（深圳）有限公司"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("company name should get CMP: prefix, got %v", got)
	}
}

func TestBuildRootDomainsNormalizesICP(t *testing.T) {
	got := BuildRootDomains([]string{"ICP:京ICP备11018762号-1", "京ICP备11018762号"})
	want := []string{"ICP:京ICP备11018762号"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ICP variants should normalize to one prefixed entry, got %v", got)
	}
}

func TestBuildRootDomainsDropsBlankLines(t *testing.T) {
	got := BuildRootDomains([]string{"example.com", "", "   ", "\r\n"})
	want := []string{"example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("blank lines must not produce entries, got %v", got)
	}
}

func TestBuildRootDomainsExtractsRootDomain(t *testing.T) {
	got := BuildRootDomains([]string{"a.b.example.co.uk", "APP:WeChat", "APP-ID:com.tencent.mm"})
	want := []string{"example.co.uk", "APP:WeChat", "APP-ID:com.tencent.mm"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestBuildRootDomainsKeepsOrderAndDedupes(t *testing.T) {
	got := BuildRootDomains([]string{"b.com", "a.com", "b.com"})
	want := []string{"b.com", "a.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected order-preserving dedupe, got %v", got)
	}
}
func TestBuildRootDomainsDoesNotMistakeICPHost(t *testing.T) {
	got := BuildRootDomains([]string{"icp.example.com"})
	want := []string{"example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("host containing icp must stay a domain, got %v", got)
	}
}
