package utils

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	configContent := `
api:
  base_url: "http://example.com"
server:
  port: 8080
`
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()
	tempFile.WriteString(configContent)
	config, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if config.API.BaseURL != "http://example.com" {
		t.Errorf("expected base_url to be http://example.com, got %s", config.API.BaseURL)
	}
}

func TestLoadConfig_Error(t *testing.T) {
	_, err := LoadConfig("non-existent.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}