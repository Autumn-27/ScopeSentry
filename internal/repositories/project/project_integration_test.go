//go:build integration

package project

import (
	"context"
	"testing"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// #242 回归：更新 A 项目的资产关联，不能把 B 项目已经关联好的资产清空。
func TestUpdateAssetsProjectKeepsOtherProjects(t *testing.T) {
	if mongodb.DB == nil {
		t.Skip("mongo not initialized")
	}
	ctx := context.Background()
	coll := mongodb.DB.Collection("RootDomain")

	domA := "proj-a-issue242.example"
	domB := "proj-b-issue242.example"
	filter := bson.M{"domain": bson.M{"$in": []string{domA, domB}}}
	_, _ = coll.DeleteMany(ctx, filter)
	t.Cleanup(func() { _, _ = coll.DeleteMany(ctx, filter) })

	if _, err := coll.InsertMany(ctx, []interface{}{
		bson.M{"domain": domA, "project": "PROJ_A", "icp": ""},
		bson.M{"domain": domB, "project": "PROJ_B", "icp": ""},
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	repo := NewRepository()
	if err := repo.UpdateAssetsProject(ctx, []string{domA}, "PROJ_A", true); err != nil {
		t.Fatalf("update: %v", err)
	}

	var docB bson.M
	if err := coll.FindOne(ctx, bson.M{"domain": domB}).Decode(&docB); err != nil {
		t.Fatalf("find B: %v", err)
	}
	if got, _ := docB["project"].(string); got != "PROJ_B" {
		t.Fatalf("project B association must be preserved, got %q", got)
	}

	var docA bson.M
	if err := coll.FindOne(ctx, bson.M{"domain": domA}).Decode(&docA); err != nil {
		t.Fatalf("find A: %v", err)
	}
	if got, _ := docA["project"].(string); got != "PROJ_A" {
		t.Fatalf("project A association should be kept, got %q", got)
	}
}
