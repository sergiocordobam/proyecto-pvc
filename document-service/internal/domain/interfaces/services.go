package interfaces

import (
	"context"
	models2 "document-service/internal/domain/models"
)

type DocumentServiceInterface interface {
	UploadFiles(ctx context.Context, uploadRequest models2.UploadRequest) (models2.UploadResponse, error)
	DownloadFiles(ctx context.Context, downloadRequest models2.DownloadRequest) (models2.DownloadResponse, error)
	DownloadAllFilesFromUser(ctx context.Context, userID int) (models2.DownloadAllResponse, error)
	GetUserDocuments(ctx context.Context, userID int) ([]models2.Document, error)
	DeleteSelectedFileInUserDirectory(ctx context.Context, userID int, files string) error
	DeleteAllFilesInUserDirectory(ctx context.Context, userID int) error
	AuthDocuments(ctx context.Context, request models2.AuthDocRequest) error
	TransferDocsToCurrentBucket(ctx context.Context, registerDocsReq models2.RegisterDocumentsMessage) error
}
