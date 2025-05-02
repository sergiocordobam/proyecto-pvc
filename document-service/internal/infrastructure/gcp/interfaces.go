package gcp

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
)

type StorageClientInterface interface {
	WriteObject(ctx context.Context, path string, content []byte) error
	ReadObject(ctx context.Context, path string) ([]byte, error)
	ObjectExists(ctx context.Context, prefix string) (bool, error)
	DeleteObject(ctx context.Context, path string) error
	SetObjectAttributes(ctx context.Context, path string, attrs storage.ObjectAttrsToUpdate) error
	GenerateSignedURL(path string, method string, expiry time.Duration) (string, error)
	ListObjectsWithPrefix(ctx context.Context, prefix string) ([]string, error)
	UserExits(ctx context.Context, userID uint64) bool
}
