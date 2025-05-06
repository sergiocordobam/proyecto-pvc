package services

import (
	"context"
	"document-service/internal/models"
)

type DocumentServiceInterface interface {
	UploadFiles(ctx context.Context, uploadRequest models.UploadRequest) (models.UploadResponse, error)
	DownloadFiles(ctx context.Context, downloadRequest models.DownloadRequest) (models.DownloadResponse, error)
	GetDocumentData(userID int, documentName string) (models.Document, error)
	GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error)
}
