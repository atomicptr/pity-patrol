package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	UserAgent       string     `toml:"user-agent,omitempty"`
	EnableScheduler bool       `toml:"enable-scheduler,omitempty"`
	DebugMode       bool       `toml:"debug-mode,omitempty"`
	Accounts        []Account  `toml:"accounts"`
	Reporters       []Reporter `toml:"reporters,omitempty"`
}

func Load() (*Config, error) {
	if path := os.Getenv("PITY_PATROL_CONFIG"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return FromPath(path)
		}
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find config dir: %w", err)
	}

	path := filepath.Join(configDir, "pity-patrol", "config.toml")
	return FromPath(path)
}

func FromPath(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("could not find file: %s", path)
	}

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, fmt.Errorf("failed to decode toml: %w", err)
	}

	for index, account := range conf.Accounts {
		err := account.validate()
		if err != nil {
			return nil, fmt.Errorf("invalid account #%d - %s", index+1, err)
		}
	}

	for index, reporter := range conf.Reporters {
		err := reporter.validate()
		if err != nil {
			return nil, fmt.Errorf("invalid reporter #%d - %s", index+1, err)
		}
	}

	return &conf, nil
}
