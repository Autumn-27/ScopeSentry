package taskstate

import (
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
)

// 待开始的定时扫描任务在节点产出进度之前必须保持待开始，不能直接显示为运行中（issue #288）。
func TestDerivePendingStaysPending(t *testing.T) {
	progress, status := Derive(constants.TaskStatusPending, 0, 10)
	if progress != 0 {
		t.Fatalf("progress = %d, want 0", progress)
	}
	if status != constants.TaskStatusPending {
		t.Fatalf("status = %d, want pending(%d)", status, constants.TaskStatusPending)
	}
}

func TestDerivePendingPromotedWhenProgressAppears(t *testing.T) {
	progress, status := Derive(constants.TaskStatusPending, 3, 10)
	if progress != 30 {
		t.Fatalf("progress = %d, want 30", progress)
	}
	if status != constants.TaskStatusRunning {
		t.Fatalf("status = %d, want running(%d)", status, constants.TaskStatusRunning)
	}
}

func TestDeriveFinishedAtHundredPercent(t *testing.T) {
	progress, status := Derive(constants.TaskStatusRunning, 10, 10)
	if progress != 100 {
		t.Fatalf("progress = %d, want 100", progress)
	}
	if status != constants.TaskStatusFinished {
		t.Fatalf("status = %d, want finished(%d)", status, constants.TaskStatusFinished)
	}
}

func TestDeriveRunningKeepsRunning(t *testing.T) {
	progress, status := Derive(constants.TaskStatusRunning, 5, 10)
	if progress != 50 || status != constants.TaskStatusRunning {
		t.Fatalf("Derive = (%d, %d), want (50, %d)", progress, status, constants.TaskStatusRunning)
	}
}

func TestDerivePausedIsNotResurrected(t *testing.T) {
	progress, status := Derive(constants.TaskStatusPaused, 5, 10)
	if progress != 50 || status != constants.TaskStatusPaused {
		t.Fatalf("Derive = (%d, %d), want (50, %d)", progress, status, constants.TaskStatusPaused)
	}
}

func TestDeriveClampsOvershoot(t *testing.T) {
	progress, status := Derive(constants.TaskStatusRunning, 13, 10)
	if progress != 100 || status != constants.TaskStatusFinished {
		t.Fatalf("Derive = (%d, %d), want (100, %d)", progress, status, constants.TaskStatusFinished)
	}
}

func TestDeriveZeroTargetsKeepsStatus(t *testing.T) {
	progress, status := Derive(constants.TaskStatusPending, 0, 0)
	if progress != 0 || status != constants.TaskStatusPending {
		t.Fatalf("Derive = (%d, %d), want (0, %d)", progress, status, constants.TaskStatusPending)
	}
}
