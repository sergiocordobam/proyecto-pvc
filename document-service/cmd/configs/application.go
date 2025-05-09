package configs

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/domain/configsDomain"
	"document-service/internal/domain/interfaces"
	"document-service/internal/infrastructure/message_queue"
	"document-service/internal/repository"
	"document-service/internal/services"
)

const (
	RetryStrategyLinear      = "linear"
	RetryStrategyConstant    = "constant"
	RetryStrategyExponential = "exponential"
)

// NewApplication initializes all application components
func NewApplication(gcpClient interfaces.StorageClientInterface,
	govCarpetaClient interfaces.GovCarpetaClientInterface,
	config configsDomain.Config,
) *configsDomain.Application {
	objectsRepo := repository.NewObjectStorageRepository(gcpClient, govCarpetaClient)
	documentsService := services.NewDocumentLoadService(objectsRepo)
	documentHandler := handlers.NewDocumentLoaderHandler(documentsService)
	newMessageHandler := handlers.NewMessageHandler(documentsService)
	queueConsumer := message_queue.NewRabbitMQConsumer(config.RabbitMQURL, newMessageHandler)
	return &configsDomain.Application{
		Repository:    objectsRepo,
		Service:       documentsService,
		Handler:       documentHandler,
		QueueConsumer: queueConsumer,
		Config:        config,
	}
}
