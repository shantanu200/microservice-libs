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

type Config struct {
	Name string
	Type string
	Path []string
}

func LoadConfig(config *Config) error {
	var e error

	once.Do(func() {
		v = viper.New()

		if config.Name != "" {
			e = fmt.Errorf("config name not found")
			return
		}

		if config.Type != "" {
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

	return e
}
