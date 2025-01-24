package utils

import (
	"backend-wolt-go/internal/models"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads a YAML configuration file from the specified path
// and decodes it into a Config struct.
//
// Parameters:
// - configPath: Path to the YAML configuration file.
//
// Returns:
// - models.Config: A struct containing the configuration values.
// - error: An error if the file cannot be opened or the YAML cannot be decoded.
func LoadConfig(configPath string) (models.Config, error) {
	// Open the configuration file.
	file, err := os.Open(configPath)
	if err != nil {
		return models.Config{}, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	// Initialize a Config struct to hold the decoded configuration.
	var config models.Config

	// Create a new YAML decoder and decode the file into the Config struct.
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return models.Config{}, fmt.Errorf("could not decode config file: %w", err)
	}

	// Return the decoded configuration.
	return config, nil
}
