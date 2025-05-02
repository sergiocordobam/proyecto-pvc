package repository

import (
	"context"
	"document-service/internal/models"
)

type ObjectStorageRepositoryInterface interface {
	GenerateUploadSignedURL(document models.Document) (string, error)
	GetDocumentData(ctx context.Context, userID uint64, documentName string) (models.Document, error)
	GetUserDocuments(ctx context.Context, userID uint64) (models.Document, error)
}
