package tests

import (
	"os"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/config"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Setenv("MB_HOST", "https://metabase.example.com")
	os.Setenv("MB_API_KEY", "test-api-key")
	defer os.Unsetenv("MB_HOST")
	defer os.Unsetenv("MB_API_KEY")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Host != "https://metabase.example.com" {
		t.Errorf("expected host 'https://metabase.example.com', got '%s'", cfg.Host)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("expected api key 'test-api-key', got '%s'", cfg.APIKey)
	}
}

func TestLoadConfig_MissingHost(t *testing.T) {
	os.Unsetenv("MB_HOST")
	os.Setenv("MB_API_KEY", "test-api-key")
	defer os.Unsetenv("MB_API_KEY")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error when MB_HOST is missing")
	}

	expected := "MB_HOST environment variable is required"
	if err.Error() != expected {
		t.Errorf("expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestLoadConfig_MissingAPIKey(t *testing.T) {
	os.Setenv("MB_HOST", "https://metabase.example.com")
	os.Unsetenv("MB_API_KEY")
	defer os.Unsetenv("MB_HOST")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error when MB_API_KEY is missing")
	}

	expected := "MB_API_KEY environment variable is required"
	if err.Error() != expected {
		t.Errorf("expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestLoadConfig_BothMissing(t *testing.T) {
	os.Unsetenv("MB_HOST")
	os.Unsetenv("MB_API_KEY")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error when both env vars are missing")
	}
}
