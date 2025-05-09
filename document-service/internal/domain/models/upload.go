package models

import (
	"time"
)

type FileUploadInfo struct {
	FileName       string            `json:"fileName"`
	ContentType    string            `json:"contentType"`
	DocumentType   string            `json:"documentType"`
	Size           int               `json:"size"`
	CustomMetadata CustomMetadataReq `json:"customMetadata"`
}
type UploadRequest struct {
	UserID int              `json:"userId"`
	Files  []FileUploadInfo `json:"files"`
}
type SignedUrlInfo struct {
	FileName       string                 `json:"fileName"`
	SignedUrl      string                 `json:"signedUrl"`
	ExpiresAt      time.Time              `json:"expiresAt"`
	ContentType    string                 `json:"contentType"`
	CustomMetadata CustomMetadataResponse `json:"customMetadata"`
}
type UploadResponse struct {
	StatusCode   int             `json:"status"`
	DocumentsURL []SignedUrlInfo `json:"documents"`
}
type CustomMetadataReq struct {
	Status       string `json:"status"`
	DocumentType string `json:"documentType"`
}
type CustomMetadataResponse struct {
	Status       string `json:"x-goog-meta-status"`
	DocumentType string `json:"x-goog-meta-document-type"`
}
