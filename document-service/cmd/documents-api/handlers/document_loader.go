package handlers

import (
	"document-service/cmd/pkg"
	"document-service/internal/domain/interfaces"
	models2 "document-service/internal/domain/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DocumentLoaderLoaderModules struct {
	Service                interfaces.DocumentServiceInterface
	notificationsPublisher interfaces.SendNotificationServiceInterface
}

func NewDocumentLoaderHandler(service interfaces.DocumentServiceInterface, publisher interfaces.SendNotificationServiceInterface) *DocumentLoaderLoaderModules {
	return &DocumentLoaderLoaderModules{
		Service:                service,
		notificationsPublisher: publisher,
	}
}
func (h *DocumentLoaderLoaderModules) HandleDocumentUploadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		var uploadRequest models2.UploadRequest
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
func (h *DocumentLoaderLoaderModules) HandleDocumentDownloadSignedURLRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		var fileData models2.Files
		if err := json.NewDecoder(r.Body).Decode(&fileData); err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid request body: ", err.Error())
		}
		userID := chi.URLParam(r, "user_id")

		userIDInt, err := strconv.Atoi(userID)
		downloadRequest := models2.DownloadRequest{
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
func (h *DocumentLoaderLoaderModules) HandleDocumentsListByUser() http.HandlerFunc {
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
func (h *DocumentLoaderLoaderModules) HandleReturnAllDownloadURL() http.HandlerFunc {
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
func (h *DocumentLoaderLoaderModules) HandleDeleteSelectedFile() http.HandlerFunc {
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
func (h *DocumentLoaderLoaderModules) HandleDeleteAllFiles() http.HandlerFunc {
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
func (h *DocumentLoaderLoaderModules) HandleAuthDocuments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		var documents models2.AuthDocRequest
		if err := json.NewDecoder(r.Body).Decode(&documents); err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid request body: ", err.Error())
			return
		}
		err := h.Service.AuthDocuments(r.Context(), documents)
		if err != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error HandleAuthDocuments: %s", err.Error())
			return
		}
		// Send the notification
		notificationRequest := models2.NotificationMessage{
			User:  documents.Owner,
			Name:  documents.Name,
			Event: "document_auth",
			ExtraData: map[string]interface{}{
				"title":   "Ya se han autorizado tus documentos",
				"message": "Tus documentos han sido autorizados por el administrador",
				"files":   documents.Files,
			},
		}
		errNotificationsPublisher := h.notificationsPublisher.SendNotification(ctx, notificationRequest)
		if errNotificationsPublisher != nil {
			pkg.Error(w, http.StatusFailedDependency, "Error publishing notification: %s", errNotificationsPublisher.Error())
			return
		}

		pkg.Success(w, http.StatusOK, "ok")
	}
}
