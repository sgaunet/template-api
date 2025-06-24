// Package config manages the application configuration.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"gopkg.in/yaml.v2"
)

// Config is the configuration for the application.
type Config struct {
	DBDSN    string `yaml:"dbdsn"`
	RedisDSN string `yaml:"redisdsn"`
	// RedisStream     string `mapstructure:"redisstream"`
}

// LoadConfigFromFile loads the configuration from a file.
func LoadConfigFromFile(filename string) (Config, error) {
	var yamlConfig Config
	//nolint:gosec
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return yamlConfig, fmt.Errorf("could not read yaml file: %w", err)
	}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return yamlConfig, fmt.Errorf("could not parse yaml file: %w", err)
	}
	return yamlConfig, nil
}

// LoadConfigFromEnvVar loads the configuration from environment variables.
func LoadConfigFromEnvVar() Config {
	var cfg Config
	cfg.DBDSN = os.Getenv("DBDSN")
	cfg.RedisDSN = os.Getenv("REDISDSN")
	return cfg
}

// ErrInvalidConfig is returned when the configuration is invalid.
var ErrInvalidConfig = errors.New("invalid configuration")

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.DBDSN == "" {
		return fmt.Errorf("%w: DBDSN is empty", ErrInvalidConfig)
	}
	if _, err := dsn.New(c.DBDSN); err != nil {
		return fmt.Errorf("%w: invalid DBDSN: %w", ErrInvalidConfig, err)
	}
	return nil
}
