package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config map[string]any

func Load(path string) (Config, error) {
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

	contents, err := os.ReadFile(loader.ConfigFileUsed())

	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	config := make(map[string]any)
	err = yaml.Unmarshal(contents, &config)

	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return Config(config), nil
}
