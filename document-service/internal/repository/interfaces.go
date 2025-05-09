package repository

import (
	"context"
	"document-service/internal/models"
	"time"
)

type ObjectStorageRepositoryInterface interface {
	GenerateUploadSignedURL(document models.Document) (string, time.Time, error)
	GenerateDownloadSignedURL(document models.Document) (string, time.Time, error)
	GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error)
	CreateUserDirectory(ctx context.Context, userID int) error
	DeleteFile(ctx context.Context, fileName string) error
	SetMetadata(ctx context.Context, fileName string, metadata map[string]string) error
	AuthDocument(ctx context.Context, document models.Document) (string, error)
}
