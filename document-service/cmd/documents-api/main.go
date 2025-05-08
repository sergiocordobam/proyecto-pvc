package main

import (
	"context"
	"document-service/cmd/documents-api/routes"
	"document-service/internal/infrastructure/apis/gcp"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	ctx := context.Background()
	storageClient, err := gcp.NewStorageClient(ctx, "document-service-api-storage")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	r := chi.NewRouter()
	router := routes.NewDocumentLoaderRoutes(r, storageClient)
	router.UseMiddlewares()
	router.MapRoutes()

	defer storageClient.Client.Close()

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
