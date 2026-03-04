package config

import "fmt"

type Account struct {
	Identifier    string `toml:"identifier,omitempty"`
	CheckinOffset int    `toml:"checkin-offset,omitempty"`
	Game
}

func (acc *Account) validate() error {
	switch acc.Type {
	case "endfield", "genshin", "starrail", "honkai", "zzz", "themis":
	default:
		return fmt.Errorf("invalid game type '%s'", acc.Type)
	}

	return nil
}

type Game struct {
	Type string `toml:"game"`

	// Endfield
	Credentials string `toml:"credentials"`
	SkGameRole  string `toml:"sk-game-role"`

	// Hoyo Games
	Cookie string `toml:"cookie"`
}

func (a *Account) GameName() string {
	switch a.Type {
	case "endfield":
		return "Arknights: Endfield"
	case "genshin":
		return "Genshin Impact"
	case "starrail":
		return "Honkai: Star Rail"
	case "honkai":
		return "Honkai Impact 3rd"
	case "zzz":
		return "Zenless Zone Zero"
	case "themis":
		return "Tears of Thermis"
	default:
		panic(fmt.Sprintf("Unknown game type: %s", a.Type))
	}
}
