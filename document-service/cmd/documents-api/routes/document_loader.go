package routes

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/infrastructure/apis/gcp"
	"document-service/internal/repository"
	"document-service/internal/services"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	service := services.NewDocumentLoadService(repo)
	handler := handlers.NewDocumentLoaderHandler(service)
	d.router.Post("/files/upload", handler.HandleDocumentUploadSignedURLRequest())
	d.router.Post("/files/download/{user_id}", handler.HandleDocumentDownloadSignedURLRequest())
	d.router.Get("/files/{user_id}", handler.HandleDocumentsListByUser())
	d.router.Get("/files/download/{user_id}/all", handler.HandleReturnAllDownloadURL())
	d.router.Delete("/files/", handler.HandleDeleteSelectedFiles())
	d.router.Delete("/files/{user_id}/all", handler.HandleDeleteAllFiles())
	d.ListRoutes()

}
func (d *DocumentLoaderRoutes) UseMiddlewares() {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	d.router.Use(middleware.Logger)
	d.router.Use(middleware.Recoverer)
	d.router.Use(corsMiddleware.Handler)
}
func (d *DocumentLoaderRoutes) ListRoutes() {
	fmt.Println("Available routes: ")
	for _, route := range d.router.Routes() {
		fmt.Println(route.Pattern)
	}
}
