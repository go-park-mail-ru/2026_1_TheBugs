package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
)

func TestFileRepo_Upload_OK(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "PutObject", func(_ *minio.Client, ctx context.Context, bucket, key string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		require.Equal(t, "test-bucket", bucket)
		require.Equal(t, "test-key.jpg", key)
		require.Equal(t, int64(1024), size)
		require.Equal(t, "image/jpeg", opts.ContentType)

		return minio.UploadInfo{}, nil
	})

	err := repo.Upload(context.Background(), "test-key.jpg", bytes.NewReader([]byte("test")), 1024, "image/jpeg")
	require.NoError(t, err)
}

func TestFileRepo_Upload_ClietError(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	expectedErr := errors.New("minio error")
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "PutObject", func(_ *minio.Client, ctx context.Context, bucket, key string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		return minio.UploadInfo{}, expectedErr
	})

	err := repo.Upload(context.Background(), "test-key.jpg", bytes.NewReader([]byte("test")), 1024, "image/jpeg")
	require.ErrorIs(t, err, expectedErr)
}

func TestFileRepo_Get_OK(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "GetObject", func(_ *minio.Client, ctx context.Context, bucket, key string, opts minio.GetObjectOptions) (*minio.Object, error) {
		require.Equal(t, "test-bucket", bucket)
		require.Equal(t, "test-key.jpg", key)
		return &minio.Object{}, nil
	})

	_, err := repo.Get(context.Background(), "test-key.jpg")
	require.NoError(t, err)
}

func TestFileRepo_Get_ClientError(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	expectedErr := errors.New("minio error")
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "GetObject", func(_ *minio.Client, ctx context.Context, bucket, key string, opts minio.GetObjectOptions) (*minio.Object, error) {
		return &minio.Object{}, expectedErr
	})

	_, err := repo.Get(context.Background(), "test-key.jpg")
	require.ErrorIs(t, err, expectedErr)
}

func TestFileRepo_Delete_OK(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "RemoveObject", func(_ *minio.Client, ctx context.Context, bucket, key string, opts minio.RemoveObjectOptions) error {
		require.Equal(t, "test-bucket", bucket)
		require.Equal(t, "test-key.jpg", key)
		return nil
	})

	err := repo.Delete(context.Background(), "test-key.jpg")
	require.NoError(t, err)
}

func TestFileRepo_Delete_ClientError(t *testing.T) {
	repo := &FileRepo{bucket: "test-bucket"}
	expectedErr := errors.New("minio error")
	monkey.PatchInstanceMethod(reflect.TypeOf(repo.client), "RemoveObject", func(_ *minio.Client, ctx context.Context, bucket, key string, opts minio.RemoveObjectOptions) error {
		require.Equal(t, "test-bucket", bucket)
		require.Equal(t, "test-key.jpg", key)
		return expectedErr
	})

	err := repo.Delete(context.Background(), "test-key.jpg")
	require.ErrorIs(t, err, expectedErr)
}
