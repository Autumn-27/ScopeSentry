package project

import (
	"context"
	"fmt"
	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"strings"
	"unicode"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Repository interface {
	GetProjectsByTag(ctx context.Context, projection bson.M) ([]models.Project, error)
	GetTarget(ctx context.Context, id string) (string, error)
	GetTargets(ctx context.Context, ids []string) (string, error)

	AggregateTagCounts(ctx context.Context, search string) (map[string]int, error)
	ListProjectsByTag(ctx context.Context, search, tag string, pageIndex, pageSize int) ([]models.ProjectBrief, error)
	FindByID(ctx context.Context, id string) (*models.Project, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	InsertProject(ctx context.Context, p *models.Project) (string, error)
	UpsertProjectTarget(ctx context.Context, id string, target string) error
	CreateOrUpdateProjectSchedule(ctx context.Context, id string, hour int, state bool, name string) error
	RemoveProjectSchedule(ctx context.Context, id string) error
	DeleteProjects(ctx context.Context, ids []string) error
	DeleteProjectTargets(ctx context.Context, ids []string) error
	DeleteProjectAsset(ctx context.Context, ids []string) error
	UpdateProject(ctx context.Context, p *models.UpdateProject) error
	UpdateAssetsProject(ctx context.Context, rootDomains []string, projectID string, change bool) error
}

type repository struct {
	collection       *mongo.Collection
	targetCollection *mongo.Collection
	scheduled        *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		collection:       mongodb.DB.Collection("project"),
		targetCollection: mongodb.DB.Collection("ProjectTargetData"),
		scheduled:        mongodb.DB.Collection("ScheduledTasks"),
	}
}

func (r *repository) GetProjectsByTag(ctx context.Context, projection bson.M) ([]models.Project, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, fmt.Errorf("failed to find projects: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

func (r *repository) GetTarget(ctx context.Context, id string) (string, error) {
	var result struct {
		Target string `bson:"target"`
	}
	err := r.targetCollection.FindOne(
		ctx,
		bson.M{"id": id},
		options.FindOne().SetProjection(bson.M{"target": 1, "_id": 0}),
	).Decode(&result)

	if err != nil {
		return "", fmt.Errorf("failed to get target: %w", err)
	}
	return result.Target, nil
}

func (r *repository) GetTargets(ctx context.Context, ids []string) (string, error) {
	cursor, err := r.targetCollection.Find(
		ctx,
		bson.M{"id": bson.M{"$in": ids}},
		options.Find().SetProjection(bson.M{"target": 1, "_id": 0}),
	)
	if err != nil {
		return "", fmt.Errorf("failed to find targets: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Target string `bson:"target"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return "", fmt.Errorf("failed to decode targets: %w", err)
	}

	// 将所有target拼接成一个字符串
	var builder strings.Builder
	for i, result := range results {
		builder.WriteString(strings.TrimRightFunc(result.Target, unicode.IsSpace))
		if i != len(results)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String(), nil
}

func (r *repository) AggregateTagCounts(ctx context.Context, search string) (map[string]int, error) {
	var pipeline mongo.Pipeline
	if search != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"$or": []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"root_domain": bson.M{"$regex": search, "$options": "i"}},
		}}}})
	}
	pipeline = append(pipeline,
		bson.D{{Key: "$group", Value: bson.M{"_id": "$tag", "count": bson.M{"$sum": 1}}}},
		bson.D{{Key: "$sort", Value: bson.M{"count": -1}}},
	)
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	res := make(map[string]int)
	for cursor.Next(ctx) {
		var row struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		res[row.ID] = row.Count
	}
	return res, nil
}

func (r *repository) ListProjectsByTag(ctx context.Context, search, tag string, pageIndex, pageSize int) ([]models.ProjectBrief, error) {
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"root_domain": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if tag != "All" {
		filter = bson.M{"$and": []bson.M{filter, {"tag": tag}}}
	}
	opts := options.Find().
		SetProjection(bson.M{"_id": 1, "name": 1, "logo": 1, "AssetCount": 1, "tag": 1}).
		SetSort(bson.D{{Key: "AssetCount", Value: -1}}).
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var docs []struct {
		ID         primitive.ObjectID `bson:"_id"`
		Name       string             `bson:"name"`
		Logo       string             `bson:"logo"`
		AssetCount int                `bson:"AssetCount"`
		Tag        string             `bson:"tag"`
	}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	items := make([]models.ProjectBrief, 0, len(docs))
	for _, d := range docs {
		items = append(items, models.ProjectBrief{ID: d.ID.Hex(), Name: d.Name, Logo: d.Logo, AssetCount: d.AssetCount, Tag: d.Tag})
	}
	return items, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.Project, error) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}
	var p models.Project
	if err := r.collection.FindOne(ctx, bson.M{"_id": obj}).Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	cnt, err := r.collection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *repository) InsertProject(ctx context.Context, p *models.Project) (string, error) {
	res, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *repository) UpsertProjectTarget(ctx context.Context, id string, target string) error {
	_, err := r.targetCollection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"target": target}}, options.Update().SetUpsert(true))
	return err
}

func (r *repository) CreateOrUpdateProjectSchedule(ctx context.Context, id string, hour int, state bool, name string) error {
	data := bson.M{"id": id, "hour": hour, "state": state, "name": name, "type": "project"}
	_, err := r.scheduled.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": data}, options.Update().SetUpsert(true))
	return err
}

func (r *repository) RemoveProjectSchedule(ctx context.Context, id string) error {
	_, err := r.scheduled.DeleteMany(ctx, bson.M{"id": id})
	return err
}

func (r *repository) DeleteProjects(ctx context.Context, ids []string) error {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if o, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, o)
		}
	}
	_, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	return err
}

func (r *repository) DeleteProjectTargets(ctx context.Context, ids []string) error {
	_, err := r.targetCollection.DeleteMany(ctx, bson.M{"id": bson.M{"$in": ids}})
	return err
}

func (r *repository) UpdateProject(ctx context.Context, p *models.UpdateProject) error {
	objID, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return fmt.Errorf("invalid project ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.UpdateOne(ctx, filter, bson.M{"$set": p})
	return err
}

func (r *repository) DeleteProjectAsset(ctx context.Context, ids []string) error {
	for _, collectionName := range constants.AssetDBNames {
		collection := mongodb.DB.Collection(collectionName)
		_, err := collection.DeleteMany(ctx, bson.M{"project": bson.M{"$in": ids}})
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return nil
}

// UpdateAssetsProject 更新资产项目关联
func (r *repository) UpdateAssetsProject(ctx context.Context, rootDomains []string, projectID string, change bool) error {
	// 资产集合列表

	for _, collectionName := range constants.AssetDBNames {
		if change {
			err := r.assetUpdateProject(ctx, rootDomains, collectionName, projectID)
			if err != nil {
				return fmt.Errorf("failed to update assets for collection %s: %w", collectionName, err)
			}
		} else {
			err := r.assetAddProject(ctx, rootDomains, collectionName, projectID)
			if err != nil {
				return fmt.Errorf("failed to add assets for collection %s: %w", collectionName, err)
			}
		}
	}

	return nil
}

// assetAddProject 添加资产到项目
func (r *repository) assetAddProject(ctx context.Context, rootDomains []string, collectionName, projectID string) error {
	collection := mongodb.DB.Collection(collectionName)

	var query bson.M

	switch collectionName {
	case "RootDomain":
		query = r.buildRootDomainQuery(rootDomains, false)
	case "app":
		query = r.buildAppQuery(rootDomains, false)
	case "mp":
		query = r.buildMpQuery(rootDomains, false)
	default:
		query = bson.M{"rootDomain": bson.M{"$in": rootDomains}}
	}

	updateQuery := bson.M{"$set": bson.M{"project": projectID}}

	result, err := collection.UpdateMany(ctx, query, updateQuery)
	if err != nil {
		return err
	}

	logger.Info("Updated assets for collection",
		zap.String("collection", collectionName),
		zap.Int64("modified_count", result.ModifiedCount))
	return nil
}

// assetUpdateProject 更新资产项目关联
func (r *repository) assetUpdateProject(ctx context.Context, rootDomains []string, collectionName, projectID string) error {
	collection := mongodb.DB.Collection(collectionName)

	var query bson.M

	switch collectionName {
	case "RootDomain":
		query = r.buildRootDomainQuery(rootDomains, true)
	case "app":
		query = r.buildAppQuery(rootDomains, true)
	case "mp":
		query = r.buildMpQuery(rootDomains, true)
	default:
		query = bson.M{
			"$and": []bson.M{
				{"project": projectID},
				{"rootDomain": bson.M{"$nin": rootDomains}},
			},
		}
	}

	// 先清空不符合条件的项目关联
	updateQuery := bson.M{"$set": bson.M{"project": ""}}
	result, err := collection.UpdateMany(ctx, query, updateQuery)
	if err != nil {
		return err
	}

	logger.Info("Cleared assets project association",
		zap.String("collection", collectionName),
		zap.Int64("modified_count", result.ModifiedCount))

	// 然后添加符合条件的项目关联
	return r.assetAddProject(ctx, rootDomains, collectionName, projectID)
}

// buildRootDomainQuery 构建RootDomain查询条件
func (r *repository) buildRootDomainQuery(rootDomains []string, isUpdate bool) bson.M {
	regexPatterns := make([]string, len(rootDomains))
	for i, domain := range rootDomains {
		regexPatterns[i] = "^" + strings.ReplaceAll(domain, ".", "\\.")
	}
	regexPattern := strings.Join(regexPatterns, "|")

	if isUpdate {
		return bson.M{
			"$and": []bson.M{
				{"project": bson.M{"$exists": true}},
				{"domain": bson.M{"$nin": rootDomains}},
				{"company": bson.M{"$nin": rootDomains}},
				{
					"icp": bson.M{
						"$nin": rootDomains,
						"$not": bson.M{"$regex": regexPattern, "$options": "i"},
					},
				},
			},
		}
	}

	return bson.M{
		"$or": []bson.M{
			{"domain": bson.M{"$in": rootDomains}},
			{"company": bson.M{"$in": rootDomains}},
			{"icp": bson.M{"$regex": regexPattern, "$options": "i"}},
		},
	}
}

// buildAppQuery 构建app查询条件
func (r *repository) buildAppQuery(rootDomains []string, isUpdate bool) bson.M {
	regexPatterns := make([]string, len(rootDomains))
	for i, domain := range rootDomains {
		regexPatterns[i] = "^" + strings.ReplaceAll(domain, ".", "\\.")
	}
	regexPattern := strings.Join(regexPatterns, "|")

	if isUpdate {
		return bson.M{
			"$and": []bson.M{
				{"project": bson.M{"$exists": true}},
				{"company": bson.M{"$nin": rootDomains}},
				{
					"icp": bson.M{
						"$nin": rootDomains,
						"$not": bson.M{"$regex": regexPattern, "$options": "i"},
					},
				},
			},
		}
	}

	return bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$in": rootDomains}},
			{"bundleID": bson.M{"$in": rootDomains}},
			{"company": bson.M{"$in": rootDomains}},
			{"icp": bson.M{"$regex": regexPattern, "$options": "i"}},
		},
	}
}

// buildMpQuery 构建mp查询条件
func (r *repository) buildMpQuery(rootDomains []string, isUpdate bool) bson.M {
	regexPatterns := make([]string, len(rootDomains))
	for i, domain := range rootDomains {
		regexPatterns[i] = "^" + strings.ReplaceAll(domain, ".", "\\.")
	}
	regexPattern := strings.Join(regexPatterns, "|")

	if isUpdate {
		return bson.M{
			"$and": []bson.M{
				{"project": bson.M{"$exists": true}},
				{"name": bson.M{"$nin": rootDomains}},
				{"company": bson.M{"$nin": rootDomains}},
				{"bundleID": bson.M{"$nin": rootDomains}},
				{
					"icp": bson.M{
						"$nin": rootDomains,
						"$not": bson.M{"$regex": regexPattern, "$options": "i"},
					},
				},
			},
		}
	}

	return bson.M{
		"$or": []bson.M{
			{"company": bson.M{"$in": rootDomains}},
			{"icp": bson.M{"$regex": regexPattern, "$options": "i"}},
		},
	}
}
