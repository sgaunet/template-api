package config_test

import (
	"os"
	"testing"

	"github.com/sgaunet/template-api/pkg/config"
)

func TestLoadConfigFromEnvVar(t *testing.T) {
	err := os.Setenv("DBDSN", "postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable")
	if err != nil {
		t.Errorf("Error setting environment variable: %v", err)
	}
	cfg := config.LoadConfigFromEnvVar()
	if !cfg.IsValid() {
		t.Errorf("Invalid configuration")
	}
}
