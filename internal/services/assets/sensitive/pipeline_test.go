package sensitive

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// 防止 bson.D 字面量改写成 Key/Value 形式时把顺序或取值弄错：
// 直接按类型检查每个 stage 的 key 与关键取值。
func TestBuildAggregatePipelineStages(t *testing.T) {
	query := bson.M{"status": 1}
	pipeline := buildAggregatePipeline(query, 2, 10)

	wantKeys := []string{"$match", "$project", "$sort", "$group", "$sort", "$skip", "$limit", "$project"}
	if len(pipeline) != len(wantKeys) {
		t.Fatalf("expected %d stages, got %d", len(wantKeys), len(pipeline))
	}
	for i, want := range wantKeys {
		if len(pipeline[i]) == 0 {
			t.Fatalf("stage %d is empty", i)
		}
		if got := pipeline[i][0].Key; got != want {
			t.Fatalf("stage %d: expected key %q, got %q", i, want, got)
		}
	}

	// $match 必须原样带上查询条件
	match, ok := pipeline[0][0].Value.(bson.M)
	if !ok || match["status"] != 1 {
		t.Fatalf("$match value lost: %#v", pipeline[0][0].Value)
	}

	// 分页：pageIndex=2, pageSize=10 → skip 10, limit 10
	if got, _ := pipeline[5][0].Value.(int); got != 10 {
		t.Fatalf("expected skip 10, got %v", pipeline[5][0].Value)
	}
	if got, _ := pipeline[6][0].Value.(int); got != 10 {
		t.Fatalf("expected limit 10, got %v", pipeline[6][0].Value)
	}

	// $group 的 _id 必须是 url
	group, ok := pipeline[3][0].Value.(bson.D)
	if !ok || len(group) == 0 || group[0].Key != "_id" || group[0].Value != "$url" {
		t.Fatalf("$group _id should be $url, got %#v", pipeline[3][0].Value)
	}
}
