package routes

import (
	"document-service/cmd/documents-api/handlers"
	configDomain "document-service/internal/domain/configsDomain"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type DocumentLoaderRoutes struct {
	router    *chi.Mux
	apiConfig *configDomain.Application
}

func NewDocumentLoaderRoutes(router *chi.Mux, apiConfig *configDomain.Application) *DocumentLoaderRoutes {
	return &DocumentLoaderRoutes{
		router:    router,
		apiConfig: apiConfig,
	}
}
func (d *DocumentLoaderRoutes) MapRoutes() {
	service := d.apiConfig.Service
	handler := handlers.NewDocumentLoaderHandler(service)
	d.router.Post("/files/upload", handler.HandleDocumentUploadSignedURLRequest())
	d.router.Post("/files/download/{user_id}", handler.HandleDocumentDownloadSignedURLRequest())
	d.router.Get("/files/{user_id}", handler.HandleDocumentsListByUser())
	d.router.Get("/files/download/{user_id}/all", handler.HandleReturnAllDownloadURL())
	d.router.Delete("/files/{user_id}/{file_name}", handler.HandleDeleteSelectedFile())
	d.router.Delete("/files/{user_id}/all", handler.HandleDeleteAllFiles())
	d.router.Post("/auth/documents", handler.HandleAuthDocuments())
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
