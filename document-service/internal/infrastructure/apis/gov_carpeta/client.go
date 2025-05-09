package gov_carpeta

import (
	"bytes"
	"document-service/internal/domain/configsDomain"
	models2 "document-service/internal/domain/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type GovCarpetaClient struct {
	baseURL    string
	httpClient *retryablehttp.Client
}

func NewGovCarpetaClient(config configsDomain.APIConfig) *GovCarpetaClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = config.Retry.Quantity
	retryClient.RetryWaitMin = config.Retry.Min
	retryClient.RetryWaitMax = config.Retry.Max

	if config.Retry.Strategy == "exponential" {
		retryClient.Backoff = retryablehttp.DefaultBackoff
	} else {
		retryClient.Backoff = retryablehttp.LinearJitterBackoff
	}

	retryClient.HTTPClient.Timeout = config.TimeOut
	return &GovCarpetaClient{
		baseURL:    config.BaseURL,
		httpClient: retryClient,
	}
}

func (c *GovCarpetaClient) AuthenticateDocument(idCitizen int, documentURL string, documentTitle string) (*models2.AuthenticateDocumentResponse, error) {
	endpoint := "/apis/authenticateDocument"
	url := fmt.Sprintf("https://%s%s", c.baseURL, endpoint)

	requestBody := models2.AuthenticateGcpDocumentRequest{
		IdCitizen:     idCitizen,
		UrlDocument:   documentURL,
		DocumentTitle: documentTitle,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	req, err := retryablehttp.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
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
	response := models2.AuthenticateDocumentResponse{
		Code:    resp.StatusCode,
		Message: responseString,
	}
	return &response, nil
}
