package config

import (
	"fmt"
	"os"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"gopkg.in/yaml.v2"
)

// Config is the configuration for the application
type Config struct {
	DBDSN    string `yaml:"DBDSN"`
	RedisDSN string `yaml:"REDISDSN"`
	// RedisStream     string `mapstructure:"redisstream"`
}

func LoadConfigFromFile(filename string) (Config, error) {
	var yamlConfig Config
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return yamlConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return yamlConfig, err
	}
	return yamlConfig, err
}

func LoadConfigFromEnvVar() Config {
	var cfg Config
	cfg.DBDSN = os.Getenv("DBDSN")
	cfg.RedisDSN = os.Getenv("REDISDSN")
	return cfg
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() bool {
	if c.DBDSN == "" {
		return false
	}
	_, err := dsn.New(c.DBDSN)
	return err == nil
}
