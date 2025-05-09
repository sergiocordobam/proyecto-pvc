package gov_carpeta

import (
	"bytes"
	"document-service/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GovCarpetaClientInterface interface {
	AuthenticateDocument(idCitizen int, documentURL string, documentTitle string) (*models.AuthenticateDocumentResponse, error)
}

type GovCarpetaClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewGovCarpetaClient(baseURL string, timeout time.Duration) *GovCarpetaClient {
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}

	return &GovCarpetaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *GovCarpetaClient) AuthenticateDocument(idCitizen int, documentURL string, documentTitle string) (*models.AuthenticateDocumentResponse, error) {
	endpoint := "/apis/authenticateDocument"
	url := fmt.Sprintf("https://%s%s", c.baseURL, endpoint)

	requestBody := models.AuthenticateGcpDocumentRequest{
		IdCitizen:     idCitizen,
		UrlDocument:   documentURL,
		DocumentTitle: documentTitle,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseString string
	if err := json.NewDecoder(resp.Body).Decode(&responseString); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	response := models.AuthenticateDocumentResponse{
		Code:    resp.StatusCode,
		Message: responseString,
	}
	return &response, nil
}
