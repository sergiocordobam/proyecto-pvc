package repository

import (
	"context"
	"document-service/internal/models"
	"time"
)

type ObjectStorageRepositoryInterface interface {
	GenerateUploadSignedURL(document models.Document) (string, time.Time, error)
	GenerateDownloadSignedURL(document models.Document) (string, time.Time, error)
	GetDocumentData(ctx context.Context, userID int, documentName string) (models.Document, error)
	GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error)
	CreateUserDirectory(ctx context.Context, userID int) error
}
