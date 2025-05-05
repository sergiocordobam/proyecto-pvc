package services

import (
	"context"
	"document-service/internal/models"
	"document-service/internal/repository"
	validator "document-service/internal/validator"
	"net/http"
	"sync"
)

type DocumentLoadService struct {
	repository repository.ObjectStorageRepositoryInterface
	validator  validator.ReqValidatorInterface
}

func NewDocumentLoadService(repository repository.ObjectStorageRepositoryInterface) *DocumentLoadService {
	return &DocumentLoadService{
		repository: repository,
	}
}

func (d *DocumentLoadService) UploadFiles(ctx context.Context, uploadRequest models.UploadRequest) (models.UploadResponse, error) {
	errValidating := d.validator.ValidateUserID(uploadRequest.UserID)
	if errValidating != nil {
		return models.UploadResponse{
			StatusCode: http.StatusBadRequest,
		}, errValidating
	}
	err := d.repository.CreateUserDirectory(ctx, uploadRequest.UserID)
	if err != nil {
		return models.UploadResponse{
			StatusCode:   http.StatusFailedDependency,
			DocumentsURL: nil,
		}, err
	}
	signedURLs := d.callDocumentsURL(uploadRequest.UserID, uploadRequest.Files)
	response := models.UploadResponse{
		StatusCode:   http.StatusOK,
		DocumentsURL: signedURLs,
	}
	return response, nil
}

func (d DocumentLoadService) generateUploadSignedURLS(document models.Document) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentLoadService) generateDownloadSignedURLS(document models.Document) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentLoadService) GetDocumentData(userID int, documentName string) (models.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentLoadService) GetUserDocuments(userID int) (models.Document, error) {
	//TODO implement me
	panic("implement me")
}
func (d *DocumentLoadService) callDocumentsURL(userID int, fileUploadInfo []models.FileUploadInfo) []models.SignedUrlInfo {
	var wg sync.WaitGroup
	var m sync.Mutex
	signedURLInfoList := make([]models.SignedUrlInfo, len(fileUploadInfo))
	for i, file := range fileUploadInfo {
		go func(i int, file models.FileUploadInfo) {
			defer wg.Done()
			var signedURLInfo models.SignedUrlInfo
			document := models.NewDocument(file.FileName, file.DocumentType, file.ContentType, file.Size, userID)
			url, expirationTime, errObj := d.repository.GenerateUploadSignedURL(document)
			if errObj != nil {
				return
			}
			signedURLInfo.FileName = document.Metadata.Name
			signedURLInfo.SignedUrl = url
			signedURLInfo.ExpiresAt = expirationTime
			signedURLInfo.ContentType = document.Metadata.ContentType
			m.Lock()
			signedURLInfoList[i] = signedURLInfo
			m.Unlock()

		}(i, file)
	}
	wg.Wait()
	return signedURLInfoList
}
