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

func NewFileRepo(client *minio.Client, bucket string) (*FileRepo, error) {
	r := &FileRepo{
		client: client,
		bucket: bucket,
	}
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: BucketExists check failed: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: MakeBucket failed: %w", err)
		}
	}
	policyStr, err := client.GetBucketPolicy(ctx, bucket)
	if err != nil {
		fmt.Printf("minio: GetBucketPolicy error: %v\n", err)
	} else {
		fmt.Printf("minio: Policy for %s: %s\n", bucket, policyStr)
	}

	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"AllowPublicRead","Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)
	if err := client.SetBucketPolicy(ctx, bucket, policy); err != nil {
		return nil, fmt.Errorf("minio: SetBucketPolicy failed: %w", err)
	}

	return r, nil
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
