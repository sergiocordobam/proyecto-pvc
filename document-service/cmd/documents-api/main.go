package main

import (
	"context"
	"document-service/cmd/documents-api/routes"
	"document-service/internal/infrastructure/apis/gcp"
	"encoding/json"
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
	// Create a new router
	r := chi.NewRouter()
	router := routes.NewDocumentLoaderRoutes(r, storageClient)
	router.UseMiddlewares()
	router.MapRoutes()

	defer storageClient.Client.Close()

	// Start the server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
