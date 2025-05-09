package main

import (
	"document-service/cmd/configs"
	"document-service/cmd/documents-api/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	app := configs.InitializeConfigsApp()
	rabbitMQConsumer := app.QueueConsumer
	errConsumer := rabbitMQConsumer.Connect()
	if errConsumer != nil {
		log.Fatalf("Error Connect RabbitMQ: %v", errConsumer)
	}
	defer rabbitMQConsumer.Connect()
	defer app.Config.StorageClient.Close()

	r := chi.NewRouter()
	router := routes.NewDocumentLoaderRoutes(r, app)
	router.UseMiddlewares()
	router.MapRoutes()

	for _, queueName := range app.Config.QueueNames {
		go rabbitMQConsumer.Consume(queueName)
	}

	log.Println("Server listening on :8080")

	go func() {
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatalf("Error al iniciar el servidor HTTP: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan //espera señal de interrupción
	log.Println("Recibida señal de terminación, cerrando...")

	log.Println("Aplicación terminada")
}
