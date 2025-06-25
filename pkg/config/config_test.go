package config_test

import (
	"os"
	"testing"

	"github.com/sgaunet/template-api/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile_WithEnvVarOverride(t *testing.T) {
	// Create a temporary YAML file
	content := []byte("dbdsn: file-db-dsn\nredisdsn: file-redis-dsn\n")
	tmpfile, err := os.CreateTemp("", "test-config-*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write(content)
	assert.NoError(t, err)
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Set environment variables to override file values
	expectedDBDSN := "postgres://postgres:password@postgres-server:5432/postgres?sslmode=disable"
	t.Setenv("DB_DSN", expectedDBDSN)
	t.Setenv("REDIS_DSN", "env-redis-dsn")

	// Load configuration
	cfg, err := config.Load(tmpfile.Name())
	assert.NoError(t, err)

	// Assert that environment variables took precedence
	assert.Equal(t, expectedDBDSN, cfg.DBDSN)
	assert.Equal(t, "env-redis-dsn", cfg.RedisDSN)

	// Validate the final configuration
	err = cfg.Validate()
	assert.NoError(t, err)
}
