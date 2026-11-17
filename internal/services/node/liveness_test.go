package node

import (
	"testing"
	"time"
)

func mustLoad(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	return loc
}

func TestLivenessFreshHeartbeatIsAlive(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	res := ClassifyLiveness("1", "2026-11-17 11:59:30", now, loc, 300*time.Second)
	if !res.Alive || res.ClockSkew || res.ParseFailed {
		t.Fatalf("fresh heartbeat should be alive, got %+v", res)
	}
}

func TestLivenessBeyondTimeoutIsOffline(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	res := ClassifyLiveness("1", "2026-11-17 11:50:00", now, loc, 300*time.Second)
	if res.Alive {
		t.Fatalf("10 minutes old heartbeat should be offline, got %+v", res)
	}
}

// 节点跑在 UTC、服务器在 +08:00：按服务器时区解析会得到「8 小时前」，
// 以前直接判离线；现在应回退到 UTC 解释并判为在线。
func TestLivenessHandlesTimezoneMismatch(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc) // = 04:00 UTC
	res := ClassifyLiveness("1", "2026-11-17 03:59:40", now, loc, 300*time.Second)
	if !res.Alive {
		t.Fatalf("utc node should be alive after fallback, got %+v", res)
	}
	if !res.UsedUTC {
		t.Fatalf("expected the UTC interpretation to be used, got %+v", res)
	}
}

// 节点时钟快几分钟：差值可能为负，不应因此判离线。
func TestLivenessToleratesClockAhead(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	res := ClassifyLiveness("1", "2026-11-17 12:02:00", now, loc, 300*time.Second)
	if !res.Alive {
		t.Fatalf("clock skew should not mark node offline, got %+v", res)
	}
	if !res.ClockSkew {
		t.Fatalf("expected ClockSkew flag, got %+v", res)
	}
}

func TestLivenessParseFailureKeepsNode(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	res := ClassifyLiveness("1", "not-a-time", now, loc, 300*time.Second)
	if !res.Alive || !res.ParseFailed {
		t.Fatalf("unparsable timestamp must keep the node, got %+v", res)
	}
}

func TestLivenessNonRunningState(t *testing.T) {
	loc := mustLoad(t, "Asia/Shanghai")
	now := time.Now()
	for _, state := range []string{"2", "3", ""} {
		if res := ClassifyLiveness(state, "2026-11-17 12:00:00", now, loc, 300*time.Second); res.Alive {
			t.Fatalf("state %q should not be alive, got %+v", state, res)
		}
	}
}

func TestNodeTimeoutEnvOverride(t *testing.T) {
	t.Setenv("NODE_TIMEOUT_SECONDS", "600")
	if got := NodeTimeout(); got != 600*time.Second {
		t.Fatalf("expected 600s, got %v", got)
	}
	t.Setenv("NODE_TIMEOUT_SECONDS", "bad")
	if got := NodeTimeout(); got != DefaultNodeTimeout {
		t.Fatalf("invalid value should fall back to default, got %v", got)
	}
}
