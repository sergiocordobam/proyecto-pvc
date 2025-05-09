package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

const (
	fileExtensions = "(pdf|jpg|jpeg|png|txt|doc|docx|docm|xls|xlsx|csv)"
)

type TempClient struct {
	client *http.Client
}

func NewTempFileClient() *TempClient {
	return &TempClient{
		client: http.DefaultClient,
	}
}
func (t *TempClient) DownloadFileFromPresignedURL(ctx context.Context, presignedURL string) ([]byte, string, error) {
	filename := t.getFileName(presignedURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, presignedURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("fallo al crear solicitud HTTP para URL %s: %w", presignedURL, err)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fallo al ejecutar solicitud HTTP GET para URL %s: %w", presignedURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("solicitud HTTP a %s falló con código de estado %d: %s", presignedURL, resp.StatusCode, string(errorBody))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("fallo al leer cuerpo de respuesta de URL %s: %w", presignedURL, err)
	}

	return data, filename, nil
}
func (t *TempClient) DetectContentType(bytesData []byte) (string, error) {
	contentType := http.DetectContentType(bytesData)
	if contentType == "" {
		return "", fmt.Errorf("no se pudo detectar el tipo de contenido")
	}
	return contentType, nil

}
func (t *TempClient) getFileName(url string) string {
	UID := uuid.New().String()
	defaultFilename := "file-" + UID

	if url == "" {
		return ""
	}

	regex := regexp.MustCompile(`([^/]+\.(?i)(pdf|jpg|jpeg|png|txt|doc|docx|docm|xls|xlsx|csv))($|\?)`)
	matches := regex.FindStringSubmatch(url)

	if len(matches) >= 2 && matches[1] != "" {
		return matches[1]
	}

	return defaultFilename
}
