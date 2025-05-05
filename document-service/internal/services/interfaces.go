package services

import (
	"context"
	"document-service/internal/models"
)

type DocumentServiceInterface interface {
	UploadFiles(ctx context.Context, uploadRequest models.UploadRequest) (models.UploadResponse, error)
	generateUploadSignedURLS(document models.Document) (string, error)
	generateDownloadSignedURLS(document models.Document) (string, error)
	GetDocumentData(userID int, documentName string) (models.Document, error)
	GetUserDocuments(userID int) (models.Document, error)
}
