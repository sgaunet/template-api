package config

import (
	"fmt"

	"github.com/sgaunet/dsn/v2/pkg/dsn"
	"github.com/spf13/viper"
)

type Config struct {
	DbDSN string `mapstructure:"dbdsn"`
	// RedisDSN        string `mapstructure:"redisdsn"`
	// RedisStream     string `mapstructure:"redisstream"`
}

func LoadConfigFromFileOrEnvVar(cfgFilePath string) (*Config, error) {
	var C Config
	viper.SetConfigFile(cfgFilePath)
	// viper.SetConfigName(cfgFilePath) // name of config file (without extension)
	// viper.SetConfigType("yml")       // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	// viper.AddConfigPath(".") // optionally look for config in the working directory
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// return &C, fmt.Errorf("fatal error config file: %w", err)
		fmt.Printf("info: configuration file not found")
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		return &C, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return &C, nil
}

func (c *Config) IsValid() bool {
	if c.DbDSN == "" {
		return false
	}
	_, err := dsn.New(c.DbDSN)
	return err == nil
}
