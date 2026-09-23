package storage

import (
	"context"
	"io"
	"time"
)

type PresignedUpload struct {
	URL       string
	ExpiresAt time.Time
}

type ObjectInfo struct {
	Key         string
	SizeBytes   int64
	ContentType string
	ETag        string
}

type ObjectStore interface {
	PresignPut(ctx context.Context, key string, contentType string, expire time.Duration) (PresignedUpload, error)
	Stat(ctx context.Context, key string) (ObjectInfo, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}
