package config

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Environment string

const (
	EnvironmentDevelop    Environment = "development"
	EnvironmentProduction Environment = "production"
)

func getEnvironment(value string) (Environment, error) {
	switch value {
	case string(EnvironmentDevelop):
		return EnvironmentDevelop, nil
	case string(EnvironmentProduction):
		return EnvironmentProduction, nil
	}
	return "", fmt.Errorf("Invalid environment provided: %v", value)
}

type AppConfig struct {
	GoogleGeocoderAPIHost string
	// TODO: see if its possible to use a service account instead
	GoogleGeocoderAPIKey string
	Port                 string
	Environment          Environment
}

func LoadConfig() (*AppConfig, error) {
	ctx := context.Background()
	config := &AppConfig{
		GoogleGeocoderAPIHost: "https://geocode.googleapis.com/v4beta",
		GoogleGeocoderAPIKey:  "",
		Port:                  "8080",
	}

	if os.Getenv("PORT") != "" {
		config.Port = os.Getenv("PORT")
	}

	environment, err := getEnvironment(os.Getenv("ENVIRONMENT"))
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %w", err)
	}
	config.Environment = environment

	if geocoderHost := os.Getenv("GOOGLE_GEOCODER_API"); geocoderHost != "" {
		config.GoogleGeocoderAPIHost = geocoderHost
	}

	if apiKey := os.Getenv("GEOCODER_API_KEY"); apiKey != "" {
		config.GoogleGeocoderAPIKey = apiKey
	} else {
		geocodeSecretName := os.Getenv("GEOCODE_SECRET_NAME")
		geocodeSecretVersion := os.Getenv("GEOCODE_SECRET_VERSION")
		if geocodeSecretName == "" {
			return nil, fmt.Errorf("Failed to load configuration: no GEOCODE_SECRET_NAME supplied")
		}
		if geocodeSecretVersion == "" {
			return nil, fmt.Errorf("Failed to load configuration: no GEOCODE_SECRET_VERSION supplied")
		}

		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed to load configuration: %w", err)
		}
		defer client.Close()
		secretCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		apiKeySecret, err := client.AccessSecretVersion(secretCtx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: geocodeSecretName + "/versions/" + geocodeSecretVersion,
		})
		if err != nil {
			return nil, fmt.Errorf("Failed to load configuration: %w", err)
		}
		if err = verifySecretIntegrity(apiKeySecret.GetPayload().GetData(), apiKeySecret.GetPayload().GetDataCrc32C()); err != nil {
			return nil, fmt.Errorf("Failed to load configuration: geocode api unverified - %w", err)
		}
		config.GoogleGeocoderAPIKey = string(apiKeySecret.GetPayload().GetData())
	}

	return config, nil
}

func verifySecretIntegrity(payload []byte, payloadCrc32 int64) error {
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(payload, crc32c))
	if checksum != payloadCrc32 {
		return fmt.Errorf("Secret data corruption detected")
	}
	return nil
}
