package gcp

import (
	"context"
	"document-service/internal/models"
	"time"

	"cloud.google.com/go/storage"
)

type StorageClientInterface interface {
	ReadObjectData(ctx context.Context, path string) ([]byte, error)
	ObjectExists(ctx context.Context, userID int, path string) bool
	DeleteObject(ctx context.Context, path string) error
	SetObjectAttributes(ctx context.Context, objectHandler *storage.ObjectHandle, attrs storage.ObjectAttrsToUpdate) error
	GenerateSignedURL(path string, method string, metadata models.Metadata, expiry time.Time) (string, error)
	ListObjectsWithPrefix(ctx context.Context, prefix string) ([]storage.ObjectAttrs, error)
	GetBucketPointer() *storage.BucketHandle
	Close() error
}
