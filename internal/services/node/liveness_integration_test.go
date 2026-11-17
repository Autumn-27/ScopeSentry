//go:build integration

package node

import (
	"context"
	"testing"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/database/redis"
	noderepo "github.com/Autumn-27/ScopeSentry/internal/repositories/node"
)

// #254 端到端回归：节点用 UTC 写心跳、服务器在 +08:00 时，
// 按服务器时区解析会得出「8 小时前」而把节点判离线。现在应判在线。
func TestGetNodeDataKeepsUTCNodeOnline(t *testing.T) {
	if redis.Client == nil {
		t.Skip("redis not initialized")
	}
	ctx := context.Background()
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("tz: %v", err)
	}
	// 服务器 12:00 (+08:00) == 04:00 UTC；节点按 UTC 写 03:59:40
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	stamp := now.UTC().Add(-20 * time.Second).Format(nodeTimeLayout)

	key := "node:it-node-utc"
	t.Cleanup(func() { redis.Client.Del(ctx, key) })
	if err := redis.Client.HSet(ctx, key, map[string]interface{}{
		"state":      "1",
		"updateTime": stamp,
		"running":    "0",
		"finished":   "0",
	}).Err(); err != nil {
		t.Fatalf("seed redis: %v", err)
	}

	lv := ClassifyLiveness("1", stamp, now, loc, NodeTimeout())
	if !lv.Alive {
		t.Fatalf("UTC heartbeat must be alive, got %+v", lv)
	}
	if !lv.UsedUTC {
		t.Fatalf("expected UTC interpretation, got %+v", lv)
	}

	repo := noderepo.NewRepository()
	nodes, err := repo.GetAllNodes(ctx)
	if err != nil {
		t.Fatalf("get nodes: %v", err)
	}
	found := false
	for _, n := range nodes {
		if n.Name == "it-node-utc" {
			found = true
			if n.UpdateTime != stamp {
				t.Fatalf("updateTime mismatch: %q", n.UpdateTime)
			}
		}
	}
	if !found {
		t.Fatal("seeded node not returned by repository")
	}
}

// 心跳过期的节点仍应判离线（确认放宽时区/时钟后没有把真正的死节点放过）。
func TestStaleNodeStillOffline(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
	stale := now.Add(-30 * time.Minute).Format(nodeTimeLayout)
	if lv := ClassifyLiveness("1", stale, now, loc, NodeTimeout()); lv.Alive {
		t.Fatalf("30 minutes old heartbeat must be offline, got %+v", lv)
	}
}
