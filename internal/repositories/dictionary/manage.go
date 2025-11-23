package dictionary

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type ManageRepository interface {
	List(ctx context.Context) ([]models.DictionaryMeta, error)
	ExistsByNameCategory(ctx context.Context, name, category string) (bool, error)
	InsertMeta(ctx context.Context, name, category, size string) (string, error)
	UpdateMetaSize(ctx context.Context, id string, size string) error
	DeleteMeta(ctx context.Context, ids []string) error

	UploadFile(ctx context.Context, filename string, content []byte) error
	DownloadFile(ctx context.Context, filename string) ([]byte, error)
	DeleteFile(ctx context.Context, filename string) error
	ReplaceFile(ctx context.Context, filename string, content []byte) error
}

type manageRepository struct {
	collection *mongo.Collection
	bucket     *gridfs.Bucket
}

func NewManageRepository() ManageRepository {
	return &manageRepository{
		collection: mongodb.DB.Collection("dictionary"),
	}
}

func (r *manageRepository) getBucket() (*gridfs.Bucket, error) {
	if r.bucket != nil {
		return r.bucket, nil
	}
	b, err := gridfs.NewBucket(mongodb.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %v", err)
	}
	r.bucket = b
	return b, nil
}

func (r *manageRepository) List(ctx context.Context) ([]models.DictionaryMeta, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []models.DictionaryMeta
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *manageRepository) ExistsByNameCategory(ctx context.Context, name, category string) (bool, error) {
	err := r.collection.FindOne(ctx, bson.M{"name": name, "category": category}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return err == nil, err
}

func (r *manageRepository) InsertMeta(ctx context.Context, name, category, size string) (string, error) {
	res, err := r.collection.InsertOne(ctx, bson.M{"name": name, "category": category, "size": size})
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (r *manageRepository) UpdateMetaSize(ctx context.Context, id string, size string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"size": size}})
	return err
}

func (r *manageRepository) DeleteMeta(ctx context.Context, ids []string) error {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objIDs = append(objIDs, oid)
	}
	_, err := r.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	return err
}

func (r *manageRepository) UploadFile(ctx context.Context, filename string, content []byte) error {
	bucket, err := r.getBucket()
	if err != nil {
		return err
	}
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return err
	}
	defer uploadStream.Close()
	_, err = uploadStream.Write(content)
	return err
}

func (r *manageRepository) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	bucket, err := r.getBucket()
	if err != nil {
		return nil, err
	}
	ds, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return nil, err
	}
	defer ds.Close()
	file := ds.GetFile()
	buf := make([]byte, file.Length)
	_, err = ds.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (r *manageRepository) DeleteFile(ctx context.Context, filename string) error {
	bucket, err := r.getBucket()
	if err != nil {
		return err
	}
	ds, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return err
	}
	// 获取文件ID并关闭read流
	file := ds.GetFile()
	ds.Close()
	if oid, ok := file.ID.(primitive.ObjectID); ok {
		return bucket.Delete(oid)
	}
	return fmt.Errorf("invalid gridfs file id type")
}

func (r *manageRepository) ReplaceFile(ctx context.Context, filename string, content []byte) error {
	// 删除旧文件再上传
	if err := r.DeleteFile(ctx, filename); err != nil {
		return err
	}
	return r.UploadFile(ctx, filename, content)
}
