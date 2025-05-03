package models

import (
	"time"
)

type FileUploadInfo struct {
	FileName     string `json:"fileName"`
	ContentType  string `json:"contentType"`
	DocumentType string `json:"documentType"`
	Size         int    `json:"size"`
}
type UploadRequest struct {
	UserID int              `json:"userId"`
	Files  []FileUploadInfo `json:"files"`
}
type SignedUrlInfo struct {
	FileName    string    `json:"fileName"`  // Nombre del archivo solicitado (útil para el cliente)
	SignedUrl   string    `json:"signedUrl"` // La URL firmada generada
	ExpiresAt   time.Time `json:"expiresAt"` // Fecha de expiración de la URL firmada
	ContentType string    `json:"contentType"`
}
type UploadResponse struct {
	StatusCode   int             `json:"status"`
	DocumentsURL []SignedUrlInfo `json:"documents"`
}
