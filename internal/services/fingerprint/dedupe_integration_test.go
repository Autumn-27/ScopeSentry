//go:build integration

package fingerprint

import (
	"context"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// 复现 #340：历史数据里 Add 生成随机 ID，一键更新用上游 ID，
// 之前会 upsert 不中而插入第二条同名指纹。
func TestBatchAddUpdatesLegacySameNameRow(t *testing.T) {
	if mongodb.DB == nil {
		t.Skip("mongo not initialized")
	}
	ctx := context.Background()
	coll := mongodb.DB.Collection("FingerprintRules")
	name := "dedupe-test-" + t.Name()

	if _, err := coll.DeleteMany(ctx, bson.M{"name": name}); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	t.Cleanup(func() { _, _ = coll.DeleteMany(ctx, bson.M{"name": name}) })

	// 模拟历史数据：随机 ID
	legacy := models.FingerprintRule{
		Name:          name,
		Rule:          "legacy",
		FingerprintId: "SAA_LEGACY0001",
		Category:      "web",
		State:         true,
	}
	if _, err := coll.InsertOne(ctx, &legacy); err != nil {
		t.Fatalf("seed: %v", err)
	}

	svc := NewService()
	content := "fingerprint:\n  name: " + name + "\n  category: web\n  parent_category: server\n"
	_, ok, failed, err := svc.BatchAdd(ctx, []BatchAddItem{{
		ID:      "UPSTREAM_ID_1",
		Name:    name,
		Content: content,
	}})
	if err != nil || ok != 1 || failed != 0 {
		t.Fatalf("batch add failed: ok=%d failed=%d err=%v", ok, failed, err)
	}

	n, err := coll.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 row after update, got %d (duplicate regression)", n)
	}
	var got models.FingerprintRule
	if err := coll.FindOne(ctx, bson.M{"name": name}).Decode(&got); err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.FingerprintId != "UPSTREAM_ID_1" {
		t.Fatalf("fingerprint_id should be backfilled to upstream id, got %q", got.FingerprintId)
	}
}
