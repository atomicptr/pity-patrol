package claimer

import (
	"fmt"

	"github.com/atomicptr/pity-patrol/pkgs/claimer/endfield"
	"github.com/atomicptr/pity-patrol/pkgs/config"
)

func Claim(cfg *config.Config, account *config.Account) (bool, error) {
	switch account.Type {
	case "endfield":
		return endfield.Claim(cfg, account)
	default:
		return false, fmt.Errorf("unknown game type: %s", account.Type)
	}
}
