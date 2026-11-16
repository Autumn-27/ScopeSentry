package tasksync

import (
	"strings"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/models"
)

func TestValidateRejectsEmptyIDs(t *testing.T) {
	err := Request{Option: OptionExisting, Project: "p1"}.Validate()
	if err == nil || !strings.Contains(err.Error(), "任务") {
		t.Fatalf("Validate() = %v, want an error about missing task ids", err)
	}
}

func TestValidateRejectsUnknownOption(t *testing.T) {
	err := Request{IDS: []string{"t1"}, Option: "merge"}.Validate()
	if err == nil || !strings.Contains(err.Error(), "不支持的同步方式") {
		t.Fatalf("Validate() = %v, want an error about the unsupported option", err)
	}
}

func TestValidateExistingRequiresProject(t *testing.T) {
	err := Request{IDS: []string{"t1"}, Option: OptionExisting, Project: "  "}.Validate()
	if err == nil || !strings.Contains(err.Error(), "项目") {
		t.Fatalf("Validate() = %v, want an error about the missing project", err)
	}
}

func TestValidateNewRequiresName(t *testing.T) {
	err := Request{IDS: []string{"t1"}, Option: OptionNew, Name: ""}.Validate()
	if err == nil || !strings.Contains(err.Error(), "名称") {
		t.Fatalf("Validate() = %v, want an error about the missing project name", err)
	}
}

func TestValidateAcceptsBothOptions(t *testing.T) {
	if err := (Request{IDS: []string{"t1"}, Option: OptionExisting, Project: "p1"}).Validate(); err != nil {
		t.Fatalf("existing option rejected: %v", err)
	}
	if err := (Request{IDS: []string{"t1"}, Option: OptionNew, Name: "demo"}).Validate(); err != nil {
		t.Fatalf("new option rejected: %v", err)
	}
}

func TestExtractTargetsDedupesAcrossTasks(t *testing.T) {
	tasks := []models.Task{
		{Target: "a.com\nb.com\n\na.com"},
		{Target: " b.com \nc.com "},
	}
	got, err := ExtractTargets(tasks)
	if err != nil {
		t.Fatalf("ExtractTargets returned error: %v", err)
	}
	want := []string{"a.com", "b.com", "c.com"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("ExtractTargets = %v, want %v", got, want)
	}
}

func TestExtractTargetsEmptyTargetsReturnsError(t *testing.T) {
	if _, err := ExtractTargets([]models.Task{{Target: ""}, {Target: "   \n  "}}); err == nil {
		t.Fatal("ExtractTargets accepted empty targets, want error")
	}
	if _, err := ExtractTargets(nil); err == nil {
		t.Fatal("ExtractTargets accepted no tasks, want error")
	}
}

func TestMergeExistingKeepsOrderAndDedupes(t *testing.T) {
	got := MergeExisting("a.com\nb.com", []string{"b.com", "c.com", "a.com"})
	want := []string{"a.com", "b.com", "c.com"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("MergeExisting = %v, want %v", got, want)
	}
}

func TestMergeExistingHandlesMissingExistingTarget(t *testing.T) {
	got := MergeExisting("", []string{"x.com", "x.com", " y.com "})
	want := []string{"x.com", "y.com"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("MergeExisting = %v, want %v", got, want)
	}
}

func TestCountNew(t *testing.T) {
	if got := CountNew("a.com", []string{"a.com", "b.com"}); got != 1 {
		t.Fatalf("CountNew = %d, want 1", got)
	}
	if got := CountNew("a.com", []string{"a.com"}); got != 0 {
		t.Fatalf("CountNew = %d, want 0", got)
	}
}
