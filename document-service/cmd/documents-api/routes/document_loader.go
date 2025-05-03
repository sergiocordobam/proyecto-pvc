package routes

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/infrastructure/gcp"
	"document-service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DocumentLoaderRoutes struct {
	router    *chi.Mux
	gcpClient *gcp.StorageClient
}

func NewDocumentLoaderRoutes(router *chi.Mux, gcpClient *gcp.StorageClient) *DocumentLoaderRoutes {
	return &DocumentLoaderRoutes{router: router, gcpClient: gcpClient}
}
func (d *DocumentLoaderRoutes) MapRoutes() {
	repo := repository.NewObjectStorageRepository(d.gcpClient)
	handler := handlers.NewDocumentLoaderHandler(repo)
	d.router.Post("/files/upload", handler.HandleDocumentUploadSignedURLRequest())

}
func (d *DocumentLoaderRoutes) UseMiddlewares() {
	d.router.Use(middleware.Logger)
	d.router.Use(middleware.Recoverer)
}
