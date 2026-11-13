package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config map[string]any

func Load(path string) (Config, error) {
	var config Config

	loader := viper.New()

	if path != "" {
		loader.SetConfigFile(path)
	} else {
		loader.SetConfigName("wrap.config")
		loader.AddConfigPath("/etc/wrap/")
		loader.AddConfigPath("$HOME/.wrap")
		loader.AddConfigPath(".")
	}

	err := loader.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	err = loader.Unmarshal(&config)

	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return config, nil
}
