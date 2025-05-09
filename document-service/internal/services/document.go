package services

import (
	"context"
	"document-service/internal/models"
	"document-service/internal/repository"
	"document-service/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

const EmptyStr = ""

type DocumentLoadService struct {
	repository repository.ObjectStorageRepositoryInterface
	validator  validator.ReqValidatorInterface
}

func NewDocumentLoadService(repository repository.ObjectStorageRepositoryInterface) *DocumentLoadService {
	return &DocumentLoadService{
		repository: repository,
		validator:  validator.NewReqValidator(),
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
func (d *DocumentLoadService) DownloadFiles(ctx context.Context, downloadRequest models.DownloadRequest) (models.DownloadResponse, error) {
	errValidating := d.validator.ValidateUserID(downloadRequest.UserID)
	if errValidating != nil {
		return models.DownloadResponse{
			StatusCode: http.StatusBadRequest,
		}, errValidating
	}
	signedURLs := d.callDocumentsURLDownload(downloadRequest.UserID, downloadRequest.FileName)
	if len(signedURLs) == 0 {
		return models.DownloadResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	response := models.DownloadResponse{
		StatusCode:   http.StatusOK,
		DocumentsURL: signedURLs,
	}
	return response, nil

}

func (d *DocumentLoadService) GetUserDocuments(ctx context.Context, userID int) ([]models.Document, error) {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return []models.Document{}, err
	}
	return d.repository.GetUserDocuments(ctx, userID)
}
func (d *DocumentLoadService) callDocumentsURL(userID int, fileUploadInfo []models.FileUploadInfo) []models.SignedUrlInfo {
	var wg sync.WaitGroup
	var m sync.Mutex
	signedURLInfoList := make([]models.SignedUrlInfo, len(fileUploadInfo))
	for i, file := range fileUploadInfo {
		wg.Add(1)
		go func(i int, file models.FileUploadInfo) {
			defer wg.Done()
			var signedURLInfo models.SignedUrlInfo
			var url string
			var errObj error
			var expirationTime time.Time

			document := models.NewDocument(file.FileName, file.DocumentType, file.ContentType, file.Size, userID)
			url, expirationTime, errObj = d.repository.GenerateUploadSignedURL(document)
			if errObj != nil {
				return
			}

			signedURLInfo.FileName = document.Metadata.Name
			signedURLInfo.SignedUrl = url
			signedURLInfo.ExpiresAt = expirationTime
			signedURLInfo.ContentType = document.Metadata.ContentType
			signedURLInfo.CustomMetadata = models.CustomMetadataResponse{
				Status:       document.Metadata.Status,
				DocumentType: document.Metadata.Type,
			}
			m.Lock()
			signedURLInfoList[i] = signedURLInfo
			m.Unlock()

		}(i, file)
	}
	wg.Wait()
	return signedURLInfoList
}
func (d *DocumentLoadService) callDocumentsURLDownload(userID int, fileNames []string) []models.SignedUrlInfo {
	var wg sync.WaitGroup
	var m sync.Mutex
	signedURLInfoList := make([]models.SignedUrlInfo, len(fileNames))
	for i, name := range fileNames {
		wg.Add(1)
		go func(i int, file string) {
			defer wg.Done()
			var signedURLInfo models.SignedUrlInfo
			var url string
			var errObj error
			var expirationTime time.Time
			document := models.NewDocument(name, EmptyStr, EmptyStr, 0, userID)
			url, expirationTime, errObj = d.repository.GenerateDownloadSignedURL(document)

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

		}(i, name)
	}

	wg.Wait()
	return signedURLInfoList
}

func (d *DocumentLoadService) DownloadAllFilesFromUser(ctx context.Context, userID int) (models.DownloadAllResponse, error) {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return models.DownloadAllResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}
	allDocs, err := d.GetUserDocuments(ctx, userID)
	if err != nil {
		return models.DownloadAllResponse{
			StatusCode: http.StatusFailedDependency,
		}, errors.New("DownloadAllFilesFromUser: error getting user documents")
	}
	if len(allDocs) == 0 {
		return models.DownloadAllResponse{
			StatusCode: http.StatusNotFound,
		}, errors.New("DownloadAllFilesFromUser: user has no documents")
	}
	allFileNames := make([]string, len(allDocs))
	for i, doc := range allDocs {
		allFileNames[i] = doc.Metadata.Name
	}
	signedURLs := d.callDocumentsURLDownload(userID, allFileNames)
	if len(signedURLs) == 0 {
		return models.DownloadAllResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	onlyURLs := make([]string, len(signedURLs))
	for i, signedURL := range signedURLs {
		onlyURLs[i] = signedURL.SignedUrl
	}
	response := models.DownloadAllResponse{
		StatusCode: http.StatusOK,
		Files:      onlyURLs,
		UserID:     userID,
	}
	return response, nil
}

func (d *DocumentLoadService) DeleteSelectedFileInUserDirectory(ctx context.Context, userID int, fileName string) error {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return err
	}
	if fileName == EmptyStr {
		return errors.New("DeleteSelectedFileInUserDirectory: no file to delete")
	}

	document := models.NewDocument(fileName, EmptyStr, EmptyStr, 0, userID)
	err := d.repository.DeleteFile(ctx, document.Metadata.AbsPath)
	if err != nil {
		return errors.New("DeleteSelectedFileInUserDirectory: error deleting file")
	}
	return nil
}
func (d *DocumentLoadService) DeleteAllFilesInUserDirectory(ctx context.Context, userID int) error {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return err
	}
	allDocs, err := d.GetUserDocuments(ctx, userID)
	if err != nil {
		return errors.New("DeleteAllFilesInUserDirectory: error getting user documents")
	}
	if len(allDocs) == 0 {
		return errors.New("DeleteAllFilesInUserDirectory: user has no documents")
	}
	for _, doc := range allDocs {
		err := d.DeleteSelectedFileInUserDirectory(ctx, userID, doc.Metadata.Name)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return errors.New("DeleteAllFilesInUserDirectory: error deleting files")
	}
	return nil
}
func (d *DocumentLoadService) AuthDocuments(ctx context.Context, req models.AuthDocRequest) error {
	if len(req.Files) == 0 {
		return errors.New("AuthDocuments: no documents to authorize")
	}
	var wg sync.WaitGroup
	var errors []error
	for i, document := range req.Files {
		wg.Add(1)
		go func(i int, document string) {
			defer wg.Done()
			newDoc := models.NewDocument(document, EmptyStr, EmptyStr, 0, req.Owner)
			_, errAuth := d.repository.AuthDocument(ctx, newDoc)
			if errAuth != nil {
				log.Warn("AuthDocuments: error authorizing document", errAuth)
				errors = append(errors, errAuth)
				return
			}
			newMetadata := map[string]string{
				"status": models.VerifiedStatus,
			}
			err := d.repository.SetMetadata(ctx, newDoc.Metadata.AbsPath, newMetadata)
			if err != nil {
				log.Warn("AuthDocuments: error setting metadata")
				errors = append(errors, err)
				return
			}

		}(i, document)
	}
	wg.Wait()
	if len(errors) >= len(req.Files) {
		return fmt.Errorf("AuthDocuments: all documents failed to authorize: %v", errors)
	}
	return nil
}
