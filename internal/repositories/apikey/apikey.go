package apikey

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

const collectionName = "api_key"

type Repository interface {
	Create(ctx context.Context, key *models.ApiKey) error
	FindAll(ctx context.Context) ([]models.ApiKey, error)
	FindByHash(ctx context.Context, keyHash string) (*models.ApiKey, error)
	Delete(ctx context.Context, id string) error
	UpdateLastUsed(ctx context.Context, id primitive.ObjectID) error
	EnsureIndexes(ctx context.Context) error
}

type repository struct {
	collection *mongodriver.Collection
}

func NewRepository() Repository {
	return &repository{
		collection: mongodb.DB.Collection(collectionName),
	}
}

func (r *repository) Create(ctx context.Context, key *models.ApiKey) error {
	key.ID = primitive.NewObjectID()
	key.CreatedAt = time.Now()
	key.Enabled = true
	_, err := r.collection.InsertOne(ctx, key)
	return err
}

func (r *repository) FindAll(ctx context.Context) ([]models.ApiKey, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var keys []models.ApiKey
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *repository) FindByHash(ctx context.Context, keyHash string) (*models.ApiKey, error) {
	var key models.ApiKey
	err := r.collection.FindOne(ctx, bson.M{"keyHash": keyHash, "enabled": true}).Decode(&key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *repository) UpdateLastUsed(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"lastUsedAt": now}})
	return err
}

func (r *repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongodriver.IndexModel{
		Keys:    bson.D{{Key: "keyHash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}
