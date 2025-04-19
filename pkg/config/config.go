package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	v    *viper.Viper
	once sync.Once
)

// Config holds the configuration details required to initialize Viper
type Config struct {
	// Name is the name of the configuration file without extension
	Name string
	// Type is the format of the configuration file (e.g., "yaml", "json", "toml")
	Type string
	// Path is a slice of directories to search for the configuration file
	Path []string
}

// LoadConfig initializes the Viper configuration with the provided config details.
// It implements a singleton pattern to ensure configuration is loaded only once.
// Returns the initialized Viper instance or an error if configuration loading fails.
func LoadConfig(config *Config) (*viper.Viper, error) {
	var e error
	once.Do(func() {
		v = viper.New()
		if config.Name == "" {
			e = fmt.Errorf("config name not found")
			return
		}
		if config.Type == "" {
			e = fmt.Errorf("config type not found")
			return
		}
		if len(config.Path) == 0 {
			e = fmt.Errorf("config path not found")
			return
		}
		v.SetConfigName(config.Name)
		v.SetConfigType(config.Type)
		for _, c := range config.Path {
			v.AddConfigPath(c)
		}
		if err := v.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				fmt.Println("[CONFIG] Config file not found. Using Default config ...")
			default:
				e = fmt.Errorf("error reading config file: %w", err)
			}
		}
	})
	if e != nil {
		return nil, e
	}
	return v, nil
}

// GetInstance returns the singleton viper instance.
// Returns nil if the configuration hasn't been loaded yet.
func GetInstance() *viper.Viper {
	return v
}

// Reset clears the singleton instance, allowing for reinitialization.
// This is primarily useful for testing.
func Reset() {
	v = nil
	once = sync.Once{}
}

// LoadConfigFromEnv loads configuration from environment variables.
// Prefix is used to filter environment variables (e.g., "APP_" will load only variables starting with APP_).
func LoadConfigFromEnv(prefix string) (*viper.Viper, error) {
	once.Do(func() {
		v = viper.New()
		v.SetEnvPrefix(prefix)
		v.AutomaticEnv()
	})

	return v, nil
}
