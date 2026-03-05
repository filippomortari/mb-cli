package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host   string
	APIKey string
}

func LoadConfig() (*Config, error) {
	host := os.Getenv("MB_HOST")
	if host == "" {
		return nil, fmt.Errorf("MB_HOST environment variable is required")
	}

	apiKey := os.Getenv("MB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MB_API_KEY environment variable is required")
	}

	return &Config{
		Host:   host,
		APIKey: apiKey,
	}, nil
}
