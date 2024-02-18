package config

import (
	"fmt"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"github.com/spf13/viper"
)

// Config is the configuration for the application
type Config struct {
	DBDSN    string `mapstructure:"dbdsn"`
	RedisDSN string `mapstructure:"redisdsn"`
	// RedisStream     string `mapstructure:"redisstream"`
}

// LoadConfigFromFileOrEnvVar loads the configuration from a file or environment variable
func LoadConfigFromFileOrEnvVar(cfgFilePath string) (*Config, error) {
	var C Config
	viper.SetConfigFile(cfgFilePath)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Printf("warning: configuration file not found (will be loaded from environment variables): %s\n", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		return &C, fmt.Errorf("unable to decode into struct, %w", err)
	}
	return &C, nil
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() bool {
	if c.DBDSN == "" {
		return false
	}
	_, err := dsn.New(c.DBDSN)
	return err == nil
}
