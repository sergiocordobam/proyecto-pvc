package services

import (
	"context"
	"document-service/internal/models"
)

type DocumentServiceInterface interface {
	UploadFiles(ctx context.Context, uploadRequest models.UploadRequest) (models.UploadResponse, error)
	DownloadFiles(ctx context.Context, downloadRequest models.DownloadRequest) (models.DownloadResponse, error)
	DownloadAllFilesFromUser(ctx context.Context, userID int) (models.DownloadAllResponse, error)
	GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error)
	DeleteSelectedFileInUserDirectory(ctx context.Context, userID int, files string) error
	DeleteAllFilesInUserDirectory(ctx context.Context, userID int) error
}
