package services

import (
	"context"
	"document-service/internal/domain/interfaces"
	models2 "document-service/internal/domain/models"
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
	repository interfaces.ObjectStorageRepositoryInterface
	validator  interfaces.ReqValidatorInterface
}

func NewDocumentLoadService(repository interfaces.ObjectStorageRepositoryInterface) *DocumentLoadService {
	return &DocumentLoadService{
		repository: repository,
		validator:  validator.NewReqValidator(),
	}
}

func (d *DocumentLoadService) UploadFiles(ctx context.Context, uploadRequest models2.UploadRequest) (models2.UploadResponse, error) {
	errValidating := d.validator.ValidateUserID(uploadRequest.UserID)
	if errValidating != nil {
		return models2.UploadResponse{
			StatusCode: http.StatusBadRequest,
		}, errValidating
	}
	err := d.repository.CreateUserDirectory(ctx, uploadRequest.UserID)
	if err != nil {
		return models2.UploadResponse{
			StatusCode:   http.StatusFailedDependency,
			DocumentsURL: nil,
		}, err
	}
	signedURLs := d.callDocumentsURL(uploadRequest.UserID, uploadRequest.Files)
	response := models2.UploadResponse{
		StatusCode:   http.StatusOK,
		DocumentsURL: signedURLs,
	}
	return response, nil
}
func (d *DocumentLoadService) DownloadFiles(ctx context.Context, downloadRequest models2.DownloadRequest) (models2.DownloadResponse, error) {
	errValidating := d.validator.ValidateUserID(downloadRequest.UserID)
	if errValidating != nil {
		return models2.DownloadResponse{
			StatusCode: http.StatusBadRequest,
		}, errValidating
	}
	signedURLs := d.callDocumentsURLDownload(downloadRequest.UserID, downloadRequest.FileName)
	if len(signedURLs) == 0 {
		return models2.DownloadResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	response := models2.DownloadResponse{
		StatusCode:   http.StatusOK,
		DocumentsURL: signedURLs,
	}
	return response, nil

}

func (d *DocumentLoadService) GetUserDocuments(ctx context.Context, userID int) ([]models2.Document, error) {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return []models2.Document{}, err
	}
	docs, err := d.repository.GetUserDocuments(ctx, userID)
	if err != nil {
		return []models2.Document{}, errors.New("GetUserDocuments: error getting user documents")
	}
	filterDocs := make([]models2.Document, len(docs))
	for i, doc := range docs {
		if doc.Metadata.Name != EmptyStr {
			filterDocs[i] = doc
		}
	}
	return filterDocs, nil
}
func (d *DocumentLoadService) callDocumentsURL(userID int, fileUploadInfo []models2.FileUploadInfo) []models2.SignedUrlInfo {
	var wg sync.WaitGroup
	var m sync.Mutex
	signedURLInfoList := make([]models2.SignedUrlInfo, len(fileUploadInfo))
	for i, file := range fileUploadInfo {
		wg.Add(1)
		go func(i int, file models2.FileUploadInfo) {
			defer wg.Done()
			var signedURLInfo models2.SignedUrlInfo
			var url string
			var errObj error
			var expirationTime time.Time

			document := models2.NewDocument(file.FileName, file.DocumentType, file.ContentType, file.Size, userID)
			url, expirationTime, errObj = d.repository.GenerateUploadSignedURL(document)
			if errObj != nil {
				return
			}

			signedURLInfo.FileName = document.Metadata.Name
			signedURLInfo.SignedUrl = url
			signedURLInfo.ExpiresAt = expirationTime
			signedURLInfo.ContentType = document.Metadata.ContentType
			signedURLInfo.CustomMetadata = models2.CustomMetadataResponse{
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
func (d *DocumentLoadService) callDocumentsURLDownload(userID int, fileNames []string) []models2.SignedUrlInfo {
	var wg sync.WaitGroup
	var m sync.Mutex
	signedURLInfoList := make([]models2.SignedUrlInfo, len(fileNames))
	for i, name := range fileNames {
		wg.Add(1)
		go func(i int, file string) {
			defer wg.Done()
			var signedURLInfo models2.SignedUrlInfo
			var url string
			var errObj error
			var expirationTime time.Time
			document := models2.NewDocument(name, EmptyStr, EmptyStr, 0, userID)
			url, expirationTime, errObj = d.repository.GenerateDownloadSignedURL(document)

			if errObj != nil {
				log.Warn("callDocumentsURLDownload: error generating download signed URL", errObj)

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

func (d *DocumentLoadService) DownloadAllFilesFromUser(ctx context.Context, userID int) (models2.DownloadAllResponse, error) {
	if err := d.validator.ValidateUserID(userID); err != nil {
		return models2.DownloadAllResponse{
			StatusCode: http.StatusBadRequest,
		}, err
	}
	allDocs, err := d.GetUserDocuments(ctx, userID)
	if err != nil {
		return models2.DownloadAllResponse{
			StatusCode: http.StatusFailedDependency,
		}, errors.New("DownloadAllFilesFromUser: error getting user documents")
	}
	if len(allDocs) == 0 {
		return models2.DownloadAllResponse{
			StatusCode: http.StatusNotFound,
		}, errors.New("DownloadAllFilesFromUser: user has no documents")
	}
	allFileNames := make([]string, len(allDocs))
	for i, doc := range allDocs {
		allFileNames[i] = doc.Metadata.Name
	}
	signedURLs := d.callDocumentsURLDownload(userID, allFileNames)
	if len(signedURLs) == 0 {
		return models2.DownloadAllResponse{
			StatusCode: http.StatusNotFound,
		}, nil
	}
	onlyURLs := make([]string, len(signedURLs))
	for i, signedURL := range signedURLs {
		onlyURLs[i] = signedURL.SignedUrl
	}
	response := models2.DownloadAllResponse{
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

	document := models2.NewDocument(fileName, EmptyStr, EmptyStr, 0, userID)
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
func (d *DocumentLoadService) AuthDocuments(ctx context.Context, req models2.AuthDocRequest) error {
	if len(req.Files) == 0 {
		return errors.New("AuthDocuments: no documents to authorize")
	}
	var wg sync.WaitGroup
	var errors []error
	for i, document := range req.Files {
		wg.Add(1)
		go func(i int, document string) {
			defer wg.Done()
			newDoc := models2.NewDocument(document, EmptyStr, EmptyStr, 0, req.Owner)
			_, errAuth := d.repository.AuthDocument(ctx, newDoc)
			if errAuth != nil {
				log.Warn("AuthDocuments: error authorizing document", errAuth)
				errors = append(errors, errAuth)
				return
			}
			newMetadata := map[string]string{
				"status": models2.VerifiedStatus,
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
func (d *DocumentLoadService) TransferDocsToCurrentBucket(ctx context.Context, registerDocsReq models2.RegisterDocumentsMessage) error {
	if len(registerDocsReq.Documents) == 0 {
		return errors.New("TransferDocsToCurrentBucket: no documents to transfer")
	}
	var wg sync.WaitGroup
	var errors []error
	for i, document := range registerDocsReq.Documents {
		wg.Add(1)
		go func(i int, document string) {
			defer wg.Done()
			//downloadDoc := d.DownloadFiles()

		}(i, document)
	}
	wg.Wait()
	if len(errors) >= len(registerDocsReq.Documents) {
		return fmt.Errorf("TransferDocsToCurrentBucket: all documents failed to authorize: %v", errors)
	}
	return nil
}
