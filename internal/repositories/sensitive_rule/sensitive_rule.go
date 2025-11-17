// sensitive_rule-------------------------------------
// @file      : sensitive_rule.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/10/28 22:26
// -------------------------------------------

package sensitive_rule

import (
	"context"
	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Autumn-27/ScopeSentry-go/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	RuleCount(ctx context.Context, filter bson.M) (int64, error)
	RuleList(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.SensitiveRuleItem, error)
	RuleInsert(ctx context.Context, item *models.SensitiveRuleItem) (primitive.ObjectID, error)
	RuleUpdate(ctx context.Context, id primitive.ObjectID, update bson.M) error
	RuleUpdateStateMany(ctx context.Context, ids []primitive.ObjectID, state bool) (int64, error)
	RuleDeleteMany(ctx context.Context, ids []primitive.ObjectID) (int64, error)
}

type repository struct {
	RuleCollection *mongo.Collection
}

func NewRepository() Repository {
	return &repository{
		RuleCollection: mongodb.DB.Collection("SensitiveRule"),
	}
}

func (r *repository) RuleCount(ctx context.Context, filter bson.M) (int64, error) {
	return r.RuleCollection.CountDocuments(ctx, filter)
}

func (r *repository) RuleList(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]models.SensitiveRuleItem, error) {
	cursor, err := r.RuleCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []models.SensitiveRuleItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) RuleInsert(ctx context.Context, item *models.SensitiveRuleItem) (primitive.ObjectID, error) {
	res, err := r.RuleCollection.InsertOne(ctx, item)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *repository) RuleUpdate(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.RuleCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *repository) RuleUpdateStateMany(ctx context.Context, ids []primitive.ObjectID, state bool) (int64, error) {
	res, err := r.RuleCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{"$set": bson.M{"state": state}})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *repository) RuleDeleteMany(ctx context.Context, ids []primitive.ObjectID) (int64, error) {
	res, err := r.RuleCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
