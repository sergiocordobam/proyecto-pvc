package repository

import (
	"context"
	"document-service/internal/infrastructure/gcp"
	"document-service/internal/models"
)

type ObjectStorageRepository struct {
	gcpclient gcp.StorageClientInterface
}

func NewObjectStorageRepository(client gcp.StorageClientInterface) *ObjectStorageRepository {
	return &ObjectStorageRepository{gcpclient: client}
}
func (o ObjectStorageRepository) GenerateUploadSignedURL(document models.Document) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o ObjectStorageRepository) GenerateDownloadSignedURL(document models.Document) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o ObjectStorageRepository) GetDocument(ctx context.Context, userID uint64, documentName string) (models.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (o ObjectStorageRepository) GetUserDocuments(ctx context.Context, userID uint32) (models.Document, error) {
	//TODO implement me
	panic("implement me")
}
