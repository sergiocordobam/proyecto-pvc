package handlers

import (
	"document-service/cmd/pkg"
	"document-service/internal/models"
	"document-service/internal/services"
	"encoding/json"
	"net/http"
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
