package config

import (
	"fmt"
	"os"
)

type EnvironmentType string

const (
	EnvironmentDevelopment EnvironmentType = "development"
	EnvironmentProduction  EnvironmentType = "production"
)

type AppConfig struct {
	Environment      EnvironmentType
	APIHost          string
	APIPort          string
	GoogleMapsAPIKey string
	LogLevel         string
}

func ParseEnvironmentType(s string) (EnvironmentType, error) {
	switch s {
	case string(EnvironmentDevelopment):
		return EnvironmentDevelopment, nil
	case string(EnvironmentProduction):
		return EnvironmentProduction, nil
	default:
		return "", fmt.Errorf("invalid environment type: %s", s)
	}
}

func LoadConfig() *AppConfig {
	config := &AppConfig{
		Environment:      EnvironmentDevelopment,
		APIHost:          "",
		APIPort:          "8080",
		GoogleMapsAPIKey: "",
		LogLevel:         "info",
	}

	// Load environment
	envStr := os.Getenv("ENVIRONMENT")
	if envStr != "" {
		if envType, err := ParseEnvironmentType(envStr); err == nil {
			config.Environment = envType
		}
	}

	// Load API configuration
	if apiHost := os.Getenv("API_HOST"); apiHost != "" {
		config.APIHost = apiHost
	}

	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		config.APIPort = apiPort
	}

	// Load Google Maps API key
	if apiKey := os.Getenv("GOOGLE_MAPS_API_KEY"); apiKey != "" {
		config.GoogleMapsAPIKey = apiKey
	}

	// Load log level
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config
}
