package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigApp struct {
	// RedisDSN        string `yaml:"redisdsn"`
	// RedisStream     string `yaml:"redisstream"`
	DbDSN      string `yaml:"dbdsn"`
	DebugLevel string `yaml:"debuglevel"`
}

func ReadyamlConfigFile(filename string) (ConfigApp, error) {
	var yamlConfig ConfigApp
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

func GetConfigFromEnvVar() ConfigApp {
	var cfg ConfigApp
	cfg.DbDSN = os.Getenv("DBDSN")
	// cfg.RedisDSN = os.Getenv("REDISDSN")
	cfg.DebugLevel = os.Getenv("DEBUGLEVEL")
	return cfg
}
