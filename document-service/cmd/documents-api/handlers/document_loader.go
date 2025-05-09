package handlers

import (
	"document-service/cmd/pkg"
	"document-service/internal/models"
	"document-service/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DocumentLoaderHandler struct {
	Service services.DocumentServiceInterface
}

func NewDocumentLoaderHandler(service services.DocumentServiceInterface) *DocumentLoaderHandler {
	return &DocumentLoaderHandler{
		Service: service,
	}
}
func (h *DocumentLoaderHandler) HandleDocumentUploadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		var uploadRequest models.UploadRequest
		if err := json.NewDecoder(r.Body).Decode(&uploadRequest); err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid request body: ", err.Error())
			return
		}
		uploadResponse, err := h.Service.UploadFiles(r.Context(), uploadRequest)
		if err != nil {
			pkg.Error(w, uploadResponse.StatusCode, "Error in HandleDocumentUploadSignedURLRequest: %s", err.Error())
			return
		}

		pkg.Success(w, uploadResponse.StatusCode, uploadResponse.DocumentsURL)
	}
}
func (h *DocumentLoaderHandler) HandleDocumentDownloadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		var fileData models.Files
		if err := json.NewDecoder(r.Body).Decode(&fileData); err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid request body: ", err.Error())
		}
		userID := chi.URLParam(r, "user_id")

		userIDInt, err := strconv.Atoi(userID)
		downloadRequest := models.DownloadRequest{
			UserID:   userIDInt,
			FileName: fileData.FileName,
		}
		uploadResponse, err := h.Service.DownloadFiles(r.Context(), downloadRequest)
		if err != nil {
			pkg.Error(w, uploadResponse.StatusCode, "Error HandleDocumentDownloadSignedURLRequest: %s", err.Error())
			return
		}

		pkg.Success(w, uploadResponse.StatusCode, uploadResponse.DocumentsURL)
	}
}
func (h *DocumentLoaderHandler) HandleDocumentsListByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userID := chi.URLParam(r, "user_id")
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid user ID: ", err.Error())
			return
		}
		documents, err := h.Service.GetUserDocuments(r.Context(), userIDInt)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error HandleGetAllUserDocumentsData: %s", err.Error())
			return
		}

		pkg.Success(w, http.StatusOK, documents)
	}
}
func (h *DocumentLoaderHandler) HandleReturnAllDownloadURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userID := chi.URLParam(r, "user_id")
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid user ID: ", err.Error())
			return
		}
		document, err := h.Service.DownloadAllFilesFromUser(r.Context(), userIDInt)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error HandleGetDocumentData: %s", err.Error())
			return
		}

		pkg.Success(w, http.StatusOK, document)
	}
}
func (h *DocumentLoaderHandler) HandleDeleteSelectedFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userID := chi.URLParam(r, "user_id")
		userIDInt, errParam := strconv.Atoi(userID)
		if errParam != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid user ID: ", errParam.Error())
			return
		}
		fileName := chi.URLParam(r, "file_name")
		err := h.Service.DeleteSelectedFileInUserDirectory(r.Context(), userIDInt, fileName)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error HandleDeleteSelectedFile: %s", err.Error())
			return
		}

		pkg.Success(w, http.StatusOK, "ok")
	}
}
func (h *DocumentLoaderHandler) HandleDeleteAllFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userID := chi.URLParam(r, "user_id")
		userIDInt, errParam := strconv.Atoi(userID)
		if errParam != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid user ID: ", errParam.Error())
			return
		}
		err := h.Service.DeleteAllFilesInUserDirectory(r.Context(), userIDInt)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error HandleDeleteSelectedFile: %s", err.Error())
			return
		}
		pkg.Success(w, http.StatusOK, "file deleted successfully")
	}
}
