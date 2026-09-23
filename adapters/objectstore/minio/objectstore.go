package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type objectStore struct {
	client *minio.Client
	bucket string
}

var _ storage.ObjectStore = (*objectStore)(nil)

func NewMinIOObjectStore(cfg config.MinIOConfig) (*objectStore, error) {

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	return &objectStore{client: client, bucket: cfg.Bucket}, nil
}

func (os *objectStore) PresignPut(ctx context.Context, key string, contentType string, expire time.Duration) (storage.PresignedUpload, error) {
	url, err := os.client.PresignedPutObject(ctx, os.bucket, key, expire)
	if err != nil {
		return storage.PresignedUpload{}, fmt.Errorf("failed to get presigned url: %w", err)
	}
	return storage.PresignedUpload{
		URL:       url.String(),
		ExpiresAt: time.Now().UTC().Add(expire),
	}, nil
}

func (os *objectStore) Stat(ctx context.Context, key string) (storage.ObjectInfo, error) {
	info, err := os.client.StatObject(ctx, os.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return storage.ObjectInfo{}, fmt.Errorf("failed to get stat: %w", err)
	}
	return storage.ObjectInfo{
		Key:         key,
		SizeBytes:   info.Size,
		ContentType: info.ContentType,
		ETag:        info.ETag,
	}, nil
}

func (os *objectStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	minioObject, err := os.client.GetObject(ctx, os.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return minioObject, nil
}

func (os *objectStore) Delete(ctx context.Context, key string) error {
	err := os.client.RemoveObject(ctx, os.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete the object: %w", err)
	}
	return nil
}
