package configs

import (
	"context"
	"document-service/internal/domain/configsDomain"
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
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	app := NewApplication(ctx, config)
	return app
}
func loadConfig() (configsDomain.Config, error) {
	configData, err := configFile.ReadFile(configFileName)
	if err != nil {
		return configsDomain.Config{}, err
	}

	var config configsDomain.Config
	errYaml := yaml.Unmarshal(configData, &config)
	if errYaml != nil {
		return configsDomain.Config{}, fmt.Errorf("error unmarshalling config: %w", errYaml)
	}
	return config, nil
}
