package configsDomain

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/domain/interfaces"
	"time"
)

type Config struct {
	BucketName      string                            `yaml:"bucket_name" json:"bucket_name"`
	RabbitMQURL     string                            `yaml:"rabbitMQURL" json:"rabbitMQURL"`
	QueueNames      []string                          `yaml:"queue_names" json:"queue_names"`
	GovCarpetaConf  APIConfig                         `yaml:"gov_carpeta_conf" json:"gov_carpeta_conf"`
	PublisherConfig PublisherConfig                   `yaml:"publisher_config" json:"publisher_config"`
	StorageClient   interfaces.StorageClientInterface `yaml:"-" json:"-"`
}
type APIConfig struct {
	BaseURL string        `yaml:"base_url" json:"base_url"`
	TimeOut time.Duration `yaml:"timeout" json:"timeout" default:"30s"`
	Retry   RetryConfig   `yaml:"retry" json:"retry"`
}
type RetryConfig struct {
	Quantity int           `yaml:"quantity" json:"quantity" default:"3"`
	Strategy RetryStrategy `yaml:"strategy" json:"strategy" default:"linear"`
	Interval time.Duration `yaml:"interval" json:"interval" default:"1000ms"` // used for constant backoff strategy
	Min      time.Duration `yaml:"min" json:"min" default:"1000ms"`           // used for exponential backoff strategy
	Max      time.Duration `yaml:"max" json:"max" default:"3000ms"`           // used for exponential backoff strategy
}
type RetryStrategy string

type Application struct {
	Repository           interfaces.ObjectStorageRepositoryInterface
	DocumentsService     interfaces.DocumentServiceInterface
	Handler              *handlers.DocumentLoaderLoaderModules
	QueueConsumer        interfaces.QueueConsumerInterface
	NotificationsService interfaces.SendNotificationServiceInterface
	Config               Config
}
type PublisherConfig struct {
	Topic   string `yaml:"topic" json:"topic"`
	Project string `yaml:"project" json:"project"`
}
