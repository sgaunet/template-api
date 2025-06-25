// Package config manages the application configuration.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"gopkg.in/yaml.v3"
)

// Config is the configuration for the application.
type Config struct {
	DBDSN    string `yaml:"dbdsn" env:"DB_DSN"`
	RedisDSN string `yaml:"redisdsn" env:"REDIS_DSN"`
	// RedisStream     string `mapstructure:"redisstream"`
}

// Load loads the configuration from a file and overrides with environment variables.
func Load(filename string) (Config, error) {
	var cfg Config

	// Load config from YAML file if it exists
	if _, err := os.Stat(filename); err == nil {
		//nolint:gosec
		yamlFile, err := os.ReadFile(filename)
		if err != nil {
			return cfg, fmt.Errorf("could not read yaml file: %w", err)
		}
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			return cfg, fmt.Errorf("could not parse yaml file: %w", err)
		}
	}

	// Parse environment variables and override YAML values
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
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
