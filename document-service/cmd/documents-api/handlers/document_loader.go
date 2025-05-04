package handlers

import (
	"document-service/cmd/documents-api/validator"
	"document-service/cmd/pkg"
	"document-service/internal/models"
	"document-service/internal/repository"
	"encoding/json"
	"net/http"
	"time"
)

type DocumentLoaderHandler struct {
	ObjectStorageRepository repository.ObjectStorageRepositoryInterface
	reqValidator            validator.ReqValidatorInterface
}

func NewDocumentLoaderHandler(objectStorageRepository repository.ObjectStorageRepositoryInterface) *DocumentLoaderHandler {
	return &DocumentLoaderHandler{
		ObjectStorageRepository: objectStorageRepository,
		reqValidator:            validator.NewReqValidator(),
	}
}
func (h *DocumentLoaderHandler) HandleDocumentUploadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		var uploadRequest models.UploadRequest
		if err := json.NewDecoder(r.Body).Decode(&uploadRequest); err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		errValidating := h.reqValidator.ValidateUserID(uploadRequest.UserID)
		if errValidating != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		documents := generateDocumentStructs(uploadRequest.UserID, uploadRequest.Files)
		err := h.ObjectStorageRepository.CreateUserDirectory(r.Context(), uploadRequest.UserID)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Failed to create user directory")
			return
		}

		var signedURLs []models.SignedUrlInfo
		for i, document := range documents {
			url, errObj := h.ObjectStorageRepository.GenerateUploadSignedURL(document)
			if errObj != nil {
				pkg.Error(w, http.StatusInternalServerError, "Failed to generate signed URL")
				return
			}
			documents[i].URL = url
			signedURLs = append(signedURLs, models.SignedUrlInfo{
				FileName:    document.Metadata.Name,
				SignedUrl:   document.URL,
				ExpiresAt:   time.Now().Add(15 * time.Minute), // Set appropriate expiration
				ContentType: document.Metadata.ContentType,
			})
		}

		response := models.UploadResponse{
			StatusCode:   http.StatusOK,
			DocumentsURL: signedURLs,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func generateDocumentStructs(userID int, fileUploadInfo []models.FileUploadInfo) []models.Document {
	documentsList := make([]models.Document, len(fileUploadInfo))
	for i, file := range fileUploadInfo {
		documentsList[i] = models.NewDocument(file.FileName, file.DocumentType, file.Size, userID)
	}
	return documentsList
}
