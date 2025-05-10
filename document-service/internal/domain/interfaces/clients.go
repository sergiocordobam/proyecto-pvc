package interfaces

import (
	"context"
	"document-service/internal/domain/models"
	"time"

	"cloud.google.com/go/storage"
)

type StorageClientInterface interface {
	ObjectExists(ctx context.Context, userID int, path string) bool
	DeleteObject(ctx context.Context, path string) error
	SetObjectAttributes(ctx context.Context, objectHandler *storage.ObjectHandle, attrs storage.ObjectAttrsToUpdate) error
	GenerateSignedURL(path string, method string, metadata models.Metadata, expiry time.Time) (string, error)
	ListObjectsWithPrefix(ctx context.Context, prefix string) ([]storage.ObjectAttrs, error)
	GetBucketPointer() *storage.BucketHandle
	UploadFileBytes(ctx context.Context, fileName string, fileBytes []byte) error
	GetObjectAttributes(ctx context.Context, fileName string) (models.Document, error)
	Close() error
}
type GovCarpetaClientInterface interface {
	AuthenticateDocument(idCitizen int, documentURL string, documentTitle string) (*models.AuthenticateDocumentResponse, error)
}
type TempDownloadFilesClient interface {
	DownloadFileFromPresignedURL(ctx context.Context, presignedURL string) ([]byte, string, error)
	DetectContentType(bytesData []byte) (string, error)
}
type PubSubClientInterface interface {
	Publish(ctx context.Context, message []byte) (string, error)
	Close()
}
