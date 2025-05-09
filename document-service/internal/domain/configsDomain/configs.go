package configsDomain

import (
	"document-service/cmd/documents-api/handlers"
	"document-service/internal/domain/interfaces"
	"time"
)

type Config struct {
	BucketName     string                            `yaml:"bucket_name" json:"bucket_name"`
	RabbitMQURL    string                            `yaml:"rabbitMQURL" json:"rabbitMQURL"`
	QueueNames     []string                          `yaml:"queue_names" json:"queue_names"`
	GovCarpetaConf APIConfig                         `yaml:"gov_carpeta_conf" json:"gov_carpeta_conf"`
	StorageClient  interfaces.StorageClientInterface `yaml:"-" json:"-"`
}
type APIConfig struct {
	BaseURL        string               `yaml:"base_url" json:"base_url"`
	TimeOut        time.Duration        `yaml:"timeout" json:"timeout" default:"30s"`
	Retry          RetryConfig          `yaml:"retry" json:"retry"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker"`
}
type CircuitBreakerConfig struct {
	FailureRatio   float64       `yaml:"failure_ratio" json:"failure_ratio" default:"0.3"`    // 30%
	MinObservation int           `yaml:"min_observation" json:"min_observation" default:"10"` // 10 requests
	Window         time.Duration `yaml:"window" json:"window" default:"5000ms"`
	Cooldown       time.Duration `yaml:"cooldown" json:"cooldown" default:"1000ms" `
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
	Repository    interfaces.ObjectStorageRepositoryInterface
	Service       interfaces.DocumentServiceInterface
	Handler       *handlers.DocumentLoaderLoaderModules
	QueueConsumer interfaces.QueueConsumerInterface
	Config        Config
}
