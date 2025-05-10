package configs

import (
	"context"
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/domain/configsDomain"
	"document-service/internal/infrastructure/apis/gcp_bucket"
	message_queue2 "document-service/internal/infrastructure/apis/gcp_pub"
	"document-service/internal/infrastructure/apis/gov_carpeta"
	"document-service/internal/infrastructure/message_queue"
	"document-service/internal/repository"
	"document-service/internal/services"
	"log"

	"cloud.google.com/go/pubsub"
)

const (
	RetryStrategyLinear      = "linear"
	RetryStrategyConstant    = "constant"
	RetryStrategyExponential = "exponential"
)

// NewApplication initializes all application components
func NewApplication(ctx context.Context, config configsDomain.Config) *configsDomain.Application {
	gcpClient, err := gcp_bucket.NewStorageClient(ctx, config.BucketName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	config.StorageClient = gcpClient
	govCarpetaClient := gov_carpeta.NewGovCarpetaClient(config.GovCarpetaConf)
	pubsubClient, err := pubsub.NewClient(ctx, config.PublisherConfig.Project)
	if err != nil {
		log.Fatalf("Failed to create pubsub client: %v", err)
	}
	publisher := message_queue2.NewPubSubPublisher(pubsubClient, config.PublisherConfig.Project, config.PublisherConfig.Topic)
	objectsRepo := repository.NewObjectStorageRepository(gcpClient, govCarpetaClient)
	documentsService := services.NewDocumentLoadService(objectsRepo)
	notificationsService := services.NewSendNotificationsService(publisher)
	documentHandler := handlers.NewDocumentLoaderHandler(documentsService, notificationsService)
	newMessageHandler := handlers.NewMessageHandler(documentsService)
	queueConsumer := message_queue.NewRabbitMQConsumer(config.RabbitMQURL, newMessageHandler)
	return &configsDomain.Application{
		Repository:           objectsRepo,
		DocumentsService:     documentsService,
		Handler:              documentHandler,
		QueueConsumer:        queueConsumer,
		Config:               config,
		NotificationsService: notificationsService,
	}
}
