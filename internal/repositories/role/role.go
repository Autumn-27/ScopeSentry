package role

import (
	"context"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository() *Repository {
	return &Repository{
		collection: mongodb.DB.Collection("roles"),
	}
}

// Create 创建角色
func (r *Repository) Create(ctx context.Context, role *models.Role) error {
	role.ID = primitive.NewObjectID()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, role)
	return err
}

// GetByID 根据ID获取角色
func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *Repository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *Repository) List(ctx context.Context, page, size int64) ([]*models.Role, int64, error) {
	// 计算跳过的文档数
	skip := (page - 1) * size

	// 设置分页选项
	opts := options.Find().
		SetSkip(skip).
		SetLimit(size).
		SetSort(bson.D{{"created_at", -1}})

	// 查询总记录数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// 查询角色列表
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var roles []*models.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// Update 更新角色
func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, update *models.RoleUpdate) error {
	updateData := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if update.Name != "" {
		updateData["$set"].(bson.M)["name"] = update.Name
	}
	if update.Description != "" {
		updateData["$set"].(bson.M)["description"] = update.Description
	}
	if len(update.Permissions) > 0 {
		updateData["$set"].(bson.M)["permissions"] = update.Permissions
	}
	if update.Status != "" {
		updateData["$set"].(bson.M)["status"] = update.Status
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, updateData)
	return err
}

// Delete 删除角色
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
