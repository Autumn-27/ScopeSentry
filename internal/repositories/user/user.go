package user

import (
	"context"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	FindAll(ctx context.Context) ([]*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, user *models.User) error
	UpdateFields(ctx context.Context, id string, fields bson.M) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	collection *mongodriver.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection("user"),
	}
}

// Create 创建用户
func (r *repository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// FindByID 根据ID获取用户
func (r *repository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名获取用户
func (r *repository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 获取用户列表
func (r *repository) FindAll(ctx context.Context) ([]*models.User, error) {
	// 设置分页选项
	opts := options.Find().
		SetSort(bson.D{{"createdAt", -1}})

	// 查询用户列表
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Update 更新用户
func (r *repository) Update(ctx context.Context, id string, user *models.User) error {
	updateDoc := bson.M{
		"$set": user,
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectId}, updateDoc)
	return err
}

// UpdateFields 按字段更新（仅更新传入字段，避免零值覆盖）
func (r *repository) UpdateFields(ctx context.Context, id string, fields bson.M) error {
	updateDoc := bson.M{
		"$set": fields,
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectId}, updateDoc)
	return err
}

// Delete 删除用户
func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
