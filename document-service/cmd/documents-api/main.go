package main

import (
	"context"
	"document-service/cmd/documents-api/routes"
	"document-service/internal/infrastructure/apis/gcp"
	"document-service/internal/infrastructure/apis/gov_carpeta"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {

	ctx := context.Background()
	storageClient, err := gcp.NewStorageClient(ctx, "document-service-api-storage")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	r := chi.NewRouter()
	govCarpetaClient := gov_carpeta.NewGovCarpetaClient("govcarpeta-apis-4905ff3c005b.herokuapp.com", time.Duration(5)*time.Second)
	router := routes.NewDocumentLoaderRoutes(r, storageClient, govCarpetaClient)
	router.UseMiddlewares()
	router.MapRoutes()

	defer storageClient.Client.Close()

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
