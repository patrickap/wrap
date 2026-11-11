package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config map[string]any

func Load(path string) (Config, error) {
	loader := viper.New()
	loader.SetConfigType("yaml")

	if path == "" {
		path = "config.yaml"
	}

	loader.SetConfigFile(path)

	if err := loader.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file %q does not exist", path)
		}

		return nil, fmt.Errorf("read config file: %w", err)
	}

	contents, err := os.ReadFile(loader.ConfigFileUsed())

	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	settings := make(map[string]any)

	if err := yaml.Unmarshal(contents, &settings); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return Config(settings), nil
}
