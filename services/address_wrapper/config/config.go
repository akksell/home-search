package config

import "os"

type AppConfig struct {
	GoogleGeocoderAPIHost string
	// TODO: see if its possible to use a service account instead
	GoogleGeocoderAPIKey string
}

func LoadConfig() *AppConfig {
	config := &AppConfig{
		GoogleGeocoderAPIHost: "https://geocode.googleapis.com/v4beta",
		GoogleGeocoderAPIKey:  "",
	}

	if geocoderHost := os.Getenv("GOOGLE_GEOCODER_API"); geocoderHost != "" {
		config.GoogleGeocoderAPIHost = geocoderHost
	}

	if apiKey := os.Getenv("GEOCODER_API_KEY"); apiKey != "" {
		config.GoogleGeocoderAPIKey = apiKey
	}

	return config
}
