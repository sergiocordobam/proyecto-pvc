package models

type AuthenticateGcpDocumentRequest struct {
	IdCitizen     int    `json:"idCitizen"`
	UrlDocument   string `json:"UrlDocument"`
	DocumentTitle string `json:"documentTitle"`
}

type AuthenticateDocumentResponse struct {
	Code    int    `json:"success"`
	Message string `json:"message,omitempty"`
}
type AuthDocRequest struct {
	Files []string `json:"files"`
	Owner int      `json:"owner"`
}
