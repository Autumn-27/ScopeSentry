package project

import "testing"

func TestProjectScopeEmptyWhenBlank(t *testing.T) {
	if !projectScopeEmpty("  \n\t", "") {
		t.Fatal("expected blank target to be empty scope")
	}
}

func TestProjectScopeEmptyWhenIgnored(t *testing.T) {
	if !projectScopeEmpty("example.com", "example.com") {
		t.Fatal("expected ignored target to be empty scope")
	}
}

func TestProjectScopeNotEmptyWhenTargetRemains(t *testing.T) {
	if projectScopeEmpty("example.com", "") {
		t.Fatal("expected a valid target to remain")
	}
}
