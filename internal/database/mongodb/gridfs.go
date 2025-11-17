package mongodb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// CreateGridFSFile 创建 GridFS 文件
func CreateGridFSFile(db *mongo.Database, filename string, content []byte) error {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return fmt.Errorf("failed to open upload stream: %v", err)
	}
	defer uploadStream.Close()

	if _, err := uploadStream.Write(content); err != nil {
		return fmt.Errorf("failed to write content: %v", err)
	}

	return nil
}

// ReadGridFSFile 读取 GridFS 文件
func ReadGridFSFile(db *mongo.Database, filename string) ([]byte, error) {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	downloadStream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open download stream: %v", err)
	}
	defer downloadStream.Close()

	buffer := make([]byte, downloadStream.GetFile().Length)
	_, err = downloadStream.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %v", err)
	}

	return buffer, nil
}

// DeleteGridFSFile 删除 GridFS 文件
func DeleteGridFSFile(db *mongo.Database, filename string) error {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	if err := bucket.Delete(filename); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
