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
	FileName    string    `json:"fileName"`
	SignedUrl   string    `json:"signedUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
	ContentType string    `json:"contentType"`
}
type UploadResponse struct {
	StatusCode   int             `json:"status"`
	DocumentsURL []SignedUrlInfo `json:"documents"`
}
