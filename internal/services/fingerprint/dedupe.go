package fingerprint

import (
	"context"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// fingerprintMatch 描述一次指纹写入命中的已有文档
type fingerprintMatch struct {
	// MatchBy 取值：fingerprint_id（按外部 ID 命中）/ _id（按同名旧文档命中）/ 空（未命中，插入）
	MatchBy       string
	ExistingHexID string
}

// decideFingerprintMatch 决定该更新哪一条已有指纹：
//  1. 已有同 fingerprint_id 的文档 → 正常更新（幂等重跑一键更新）
//  2. 没有，但存在同名文档 → 更新这条旧文档并回填 fingerprint_id
//     （历史版本里 Add 生成随机 ID、一键更新用上游 ID，才会出现"同一指纹两份"）
//  3. 都没有 → 插入
func decideFingerprintMatch(byFingerprintID, byName string) fingerprintMatch {
	switch {
	case byFingerprintID != "":
		return fingerprintMatch{MatchBy: "fingerprint_id", ExistingHexID: byFingerprintID}
	case byName != "":
		return fingerprintMatch{MatchBy: "_id", ExistingHexID: byName}
	default:
		return fingerprintMatch{MatchBy: ""}
	}
}

// shouldReuseFingerprintID 判断「手动新增/导入」时是否应复用已有文档：
// 同名已存在就不该再插一条，直接更新它。
func shouldReuseFingerprintID(existingHexID string) bool {
	return existingHexID != ""
}

// lookupHexID 按条件取一条指纹文档的 _id（hex）；未命中返回空串且不报错
func (s *service) lookupHexID(ctx context.Context, filter bson.M) string {
	doc, err := s.commonRepo.FindOne(ctx, "FingerprintRules", filter, bson.M{"_id": 1})
	if err != nil || doc == nil {
		return ""
	}
	switch v := doc["_id"].(type) {
	case primitive.ObjectID:
		return v.Hex()
	case string:
		return v
	default:
		return ""
	}
}

// fingerprintSetDoc 把整份指纹字段组装成 $set 内容（含 fingerprint_id 回填）
func fingerprintSetDoc(data models.FingerprintRule) bson.M {
	return bson.M{
		"name":            data.Name,
		"rule":            data.Rule,
		"fingerprint_id":  data.FingerprintId,
		"category":        data.Category,
		"parent_category": data.ParentCategory,
		"state":           data.State,
		"company":         data.Company,
	}
}
