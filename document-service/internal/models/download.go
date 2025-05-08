package models

type DownloadRequest struct {
	UserID   int      `json:"userId"`
	FileName []string `json:"fileName"`
}
type DownloadResponse struct {
	StatusCode   int             `json:"status"`
	DocumentsURL []SignedUrlInfo `json:"documents"`
}
type Files struct {
	FileName []string `json:"file_names"`
}
type DownloadAllResponse struct {
	StatusCode int      `json:"status"`
	UserID     int      `json:"user_id"`
	Files      []string `json:"files"`
}
