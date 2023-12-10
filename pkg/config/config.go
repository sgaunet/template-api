package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DbDSN string `mapstructure:"dbdsn"`
	// RedisDSN        string `mapstructure:"redisdsn"`
	// RedisStream     string `mapstructure:"redisstream"`
}

func LoadConfigFromFile(cfgFilePath string) (*Config, error) {
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
		fmt.Fprintf(os.Stderr, "error: cannot read config file: %s\n", err)
	}
	err = viper.Unmarshal(&C)
	if err != nil {
		return &C, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return &C, nil
}
