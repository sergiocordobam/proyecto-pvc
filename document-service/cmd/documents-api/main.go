package main

import (
	"context"
	"document-service/internal/infrastructure/gcp"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/api/iterator"
)

type User struct {
	ID uint32 `json:"id"`
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ctx := context.Background()
	storageClient, err := gcp.NewStorageClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer storageClient.Client.Close()
	bkt := storageClient.Client.Bucket("document-user-admin")

	// User Creation Route
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userIDstr := strconv.Itoa(int(u.ID))
		folderName := userIDstr + "/"
		it := bkt.Objects(ctx, &storage.Query{Prefix: folderName})
		_, errObject := it.Next()

		if errors.Is(errObject, iterator.Done) { // User doesn't exist
			obj := bkt.Object(folderName + "test.txt")
			writer := obj.NewWriter(ctx)

			if _, err := writer.Write([]byte("Hello, World!")); err != nil {
				log.Printf("Error writing to object: %v", err)
				RespondWithError(w, http.StatusInternalServerError, "Cannot create user")
				return
			}
			if err := writer.Close(); err != nil {
				log.Printf("Error closing writer: %v", err)
				RespondWithError(w, http.StatusInternalServerError, "Cannot create user")
				return
			}

			currentTime := time.Now().Add(24 * time.Second)
			attrsToUpdate := storage.ObjectAttrsToUpdate{
				TemporaryHold: true,
				Retention: &storage.ObjectRetention{
					Mode:        "Locked",
					RetainUntil: currentTime,
				},
			}
			if _, err := obj.Update(ctx, attrsToUpdate); err != nil {
				log.Printf("Error updating object attributes: %v", err)
				RespondWithError(w, http.StatusInternalServerError, "Error updating user metadata")
				return
			}

			response := map[string]string{
				"message": "user created successfully",
				"user_id": userIDstr,
			}
			RespondWithJSON(w, http.StatusCreated, response)
			return
		}

		// User already exists
		response := map[string]string{
			"message": "user already exists",
		}
		RespondWithJSON(w, http.StatusCreated, response) // You may want to change code to 200 - OK.
	})

	// Get Signed URL Route
	r.Get("/user/document", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("file_name")
		userID := r.URL.Query().Get("user_id")

		if userID == "" || fileName == "" {
			RespondWithError(w, http.StatusBadRequest, "user_id and file_name are required")
			return
		}

		objFileName := userID + "/" + fileName
		signedConfiguration := &storage.SignedURLOptions{
			Scheme:  storage.SigningSchemeV4,
			Method:  "PUT",
			Expires: time.Now().Add(15 * time.Minute),
		}

		urlStr, errGeneratedURL := bkt.SignedURL(objFileName, signedConfiguration)
		if errGeneratedURL != nil {
			log.Printf("Error generating Signed URL: %v", errGeneratedURL)
			RespondWithError(w, http.StatusInternalServerError, "Cannot generate signed URL")
			return
		}

		response := map[string]string{
			"url":     urlStr,
			"user_id": userID,
		}
		RespondWithJSON(w, http.StatusOK, response)
	})

	// Ping Route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	//Hello Route
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Hello, World!",
		})
	})

	// Start the server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Helper Functions for Responding

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
