package configs

import (
	"context"
	"document-service/internal/domain/configsDomain"
	"document-service/internal/infrastructure/apis/gcp"
	"document-service/internal/infrastructure/apis/gov_carpeta"
	"embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

var (
	//go:embed config.yaml
	configFile embed.FS
)

func InitializeConfigsApp() *configsDomain.Application {
	ctx := context.Background()
	// Load configuration
	config, err := loadConfig()
	storageClient, err := gcp.NewStorageClient(ctx, config.BucketName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	config.StorageClient = storageClient
	govCarpetaClient := gov_carpeta.NewGovCarpetaClient(config.GovCarpetaConf)
	app := NewApplication(storageClient, govCarpetaClient, config)
	return app
}
func loadConfig() (configsDomain.Config, error) {
	configData, err := configFile.ReadFile(configFileName)
	if err != nil {
		return configsDomain.Config{}, err
	}

	var config configsDomain.Config
	fmt.Println("Config data:", string(configData))
	errYaml := yaml.Unmarshal(configData, &config)
	if errYaml != nil {
		return configsDomain.Config{}, fmt.Errorf("error unmarshalling config: %w", errYaml)
	}
	return config, nil
}
