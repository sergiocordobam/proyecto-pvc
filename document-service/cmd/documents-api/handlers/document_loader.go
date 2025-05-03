package handlers

import (
	"document-service/cmd/documents-api/validator"
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
	}
}
func (h *DocumentLoaderHandler) HandleDocumentUploadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var uploadRequest models.UploadRequest
		if err := json.NewDecoder(r.Body).Decode(&uploadRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}

		errValidating := h.reqValidator.ValidateUserID(uploadRequest.UserID)
		if errValidating != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": errValidating.Error()})
			return
		}

		documents := generateDocumentStructs(uploadRequest.UserID, uploadRequest.Files)
		err := h.ObjectStorageRepository.CreateUserDirectory(r.Context(), uploadRequest.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create directory"})
			return
		}

		var signedURLs []models.SignedUrlInfo
		for i, document := range documents {
			url, errObj := h.ObjectStorageRepository.GenerateUploadSignedURL(document)
			if errObj != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate URL"})
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
