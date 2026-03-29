package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type FileRepo struct {
	client *minio.Client
	bucket string
}

func NewFileRepo(client *minio.Client, bucket string) *FileRepo {
	return &FileRepo{
		client: client,
		bucket: bucket,
	}
}

// структуру сделать
func (r *FileRepo) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := r.client.PutObject(ctx, r.bucket, key, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("fileRepo.Upload: %w", err)
	}

	return nil
}

func (r *FileRepo) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := r.client.GetObject(ctx, r.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("fileRepo.Get: %w", err)
	}

	return obj, nil
}

func (r *FileRepo) Delete(ctx context.Context, key string) error {
	err := r.client.RemoveObject(ctx, r.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("fileRepo.Delete: %w", err)
	}

	return nil
}
