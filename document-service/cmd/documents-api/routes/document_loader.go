package routes

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/infrastructure/gcp"
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
	d.ListRoutes()

}
func (d *DocumentLoaderRoutes) UseMiddlewares() {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Encabezados permitidos
		ExposedHeaders:   []string{"*"},                                                       // Encabezados que el frontend puede leer
		AllowCredentials: true,                                                                // Permite cookies, encabezados de autorización, etc.
		MaxAge:           3600,                                                                // Tiempo que el navegador puede cachear la respuesta de preflight
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
