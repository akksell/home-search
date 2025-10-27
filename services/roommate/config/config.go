package config

import (
	"fmt"
	"os"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type Environment string

func getEnvironment(value string) (Environment, error) {
	switch value {
	case EnvironmentDevelopment:
		return EnvironmentDevelopment, nil
	case EnvironmentProduction:
		return EnvironmentProduction, nil
	}

	return "", fmt.Errorf("Invalid Environment used: %s", value)
}

type AppConfig struct {
	Environment               Environment
	Port                      string
	GoogleProjectId           string
	RoommateStoreInstanceDB   string
	AddressWrapperServiceHost string
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{}

	environment, err := getEnvironment(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return nil, err
	}
	config.Environment = environment

	config.Port = os.Getenv("PORT")

	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("Failed to load configuration: no google project id provided")
	}
	config.GoogleProjectId = projectId

	roommateFirestoreDB := os.Getenv("ROOMMATE_STORE_DATABASE")
	if roommateFirestoreDB == "" {
		return nil, fmt.Errorf("Failed to load configuration: no roommate store name provided")
	}
	config.RoommateStoreInstanceDB = roommateFirestoreDB

	addressWrapperSvcHost := os.Getenv("ADDRESS_WRAPPER_SERVICE_HOST")
	if addressWrapperSvcHost == "" {
		return nil, fmt.Errorf("Failed to load configuration: no address wrapper service hostname provided")
	}
	config.AddressWrapperServiceHost = addressWrapperSvcHost

	return config, nil
}
