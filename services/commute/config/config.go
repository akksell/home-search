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
	Environment           EnvironmentType
	Port                  string
	LogLevel              string
	GoogleProjectId       string
	CommuteStoreDB        string
	AddressWrapperSvcHost string
	RoommateSvcHost       string
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

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		Environment: EnvironmentDevelopment,
		Port:        "8080",
		LogLevel:    "info",
	}

	envStr := os.Getenv("ENVIRONMENT")
	if envStr != "" {
		if envType, err := ParseEnvironmentType(envStr); err == nil {
			config.Environment = envType
		}
	}

	if apiPort := os.Getenv("PORT"); apiPort != "" {
		config.Port = apiPort
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("Failed to load configuration: GOOGLE_PROJECT_ID is required")
	}
	config.GoogleProjectId = projectId

	commuteStoreDB := os.Getenv("COMMUTE_STORE_DATABASE")
	if commuteStoreDB == "" {
		return nil, fmt.Errorf("Failed to load configuration: COMMUTE_STORE_DATABASE is required")
	}
	config.CommuteStoreDB = commuteStoreDB

	addressWrapperSvcHost := os.Getenv("ADDRESS_WRAPPER_SERVICE_HOST")
	if addressWrapperSvcHost == "" {
		return nil, fmt.Errorf("Failed to load configuration: ADDRESS_WRAPPER_SERVICE_HOST is required")
	}
	config.AddressWrapperSvcHost = addressWrapperSvcHost

	roommateSvcHost := os.Getenv("ROOMMATE_SERVICE_HOST")
	if addressWrapperSvcHost == "" {
		return nil, fmt.Errorf("Failed to load configuration: ROOMMATE_SERVICE_HOST is required")
	}
	config.RoommateSvcHost = roommateSvcHost

	return config, nil
}
