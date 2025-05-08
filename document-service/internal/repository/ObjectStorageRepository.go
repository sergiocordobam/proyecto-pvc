package repository

import (
	"context"
	"document-service/internal/infrastructure/apis/gcp"
	"document-service/internal/models"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/gommon/log"
)

type ObjectStorageRepository struct {
	gcpclient gcp.StorageClientInterface
}

func NewObjectStorageRepository(client gcp.StorageClientInterface) *ObjectStorageRepository {
	return &ObjectStorageRepository{gcpclient: client}
}

func (o *ObjectStorageRepository) GenerateUploadSignedURL(document models.Document) (string, time.Time, error) {
	expiryTime := document.Metadata.CreationDate.Add(1 * time.Hour)
	url, err := o.gcpclient.GenerateSignedURL(document.Metadata.AbsPath, "up", expiryTime)
	if err != nil {
		return "", time.Time{}, err
	}
	return url, expiryTime, nil
}
func (o *ObjectStorageRepository) GenerateDownloadSignedURL(document models.Document) (string, time.Time, error) {
	if !o.gcpclient.ObjectExists(context.Background(), document.Metadata.OwnerID, document.Metadata.Name) {
		return "", time.Time{}, errors.New("object not found")
	}
	expiryTime := document.Metadata.CreationDate.Add(1 * time.Hour)
	url, err := o.gcpclient.GenerateSignedURL(document.Metadata.AbsPath, "down", expiryTime)
	if err != nil {
		return "", expiryTime, err
	}
	return url, expiryTime, nil
}

func (o *ObjectStorageRepository) GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error) {
	userDocumentsList := make([]models.Document, 0)
	objListAttributes, err := o.gcpclient.ListObjectsWithPrefix(ctx, fmt.Sprintf("%d/", userID))
	if err != nil {
		return userDocumentsList, err
	}
	for _, obj := range objListAttributes {
		metadata := models.NewMetadata(obj.Name, "", obj.ContentType, int(obj.Size), userID)
		userDocumentsList = append(userDocumentsList, models.Document{
			Metadata: metadata,
		})
	}
	return userDocumentsList, nil
}

func (o *ObjectStorageRepository) CreateUserDirectory(ctx context.Context, userID int) error {
	userIDDirectory := fmt.Sprintf("%d/", userID)
	if o.gcpclient.ObjectExists(ctx, userID, "/") {
		log.Warn("user directory already exists")
		return nil
	}
	bkt := o.gcpclient.GetBucketPointer()
	obj := bkt.Object(userIDDirectory)
	writer := obj.NewWriter(ctx)

	if err := writer.Close(); err != nil {
		return err
	}

	currentTime := time.Now().Add(24 * time.Second)
	attrsToUpdate := storage.ObjectAttrsToUpdate{
		TemporaryHold: true,
		Retention: &storage.ObjectRetention{
			Mode:        "Locked",
			RetainUntil: currentTime,
		},
	}
	if err := o.gcpclient.SetObjectAttributes(ctx, obj, attrsToUpdate); err != nil {
		return err
	}
	return nil
}

func (o *ObjectStorageRepository) DeleteFile(ctx context.Context, fileName string) error {
	return o.gcpclient.DeleteObject(ctx, fileName)
}
