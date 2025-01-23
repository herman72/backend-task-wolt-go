package utils

import (
	"backend-wolt-go/internal/models"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)


func LoadConfig(configPath string) (models.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return models.Config{}, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config models.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return models.Config{}, fmt.Errorf("could not decode config file: %w", err)
	}

	return config, nil
}