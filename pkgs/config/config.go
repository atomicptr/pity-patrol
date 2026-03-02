package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	UserAgent string    `toml:"user-agent,omitempty"`
	DebugMode bool      `toml:"debug-mode,omitempty"`
	Accounts  []Account `toml:"accounts"`
}

type Account struct {
	Identifier string `toml:"identifier,omitempty"`
	Game
}

func isValidGameType(t string) bool {
	switch t {
	case "endfield":
		return true
	default:
		return false
	}
}

type Game struct {
	Type string `toml:"game"`

	// Endfield
	Credentials string `toml:"credentials"`
	SkGameRole  string `toml:"sk-game-role"`
}

func (a *Account) GameName() string {
	switch a.Game.Type {
	case "endfield":
		return "Arknights: Endfield"
	default:
		panic(fmt.Sprintf("Unknown game type: %s", a.Game.Type))
	}
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
		if !isValidGameType(account.Type) {
			return nil, fmt.Errorf("invalid game type for account #%d '%s'", index+1, account.Type)
		}
	}

	return &conf, nil
}
