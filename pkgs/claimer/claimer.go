package claimer

import (
	"fmt"

	"github.com/atomicptr/pity-patrol/pkgs/claimer/endfield"
	"github.com/atomicptr/pity-patrol/pkgs/claimer/hoyo"
	"github.com/atomicptr/pity-patrol/pkgs/config"
	"github.com/atomicptr/pity-patrol/pkgs/report"
)

func Claim(cfg *config.Config, account *config.Account) (*report.Report, error) {
	switch account.Type {
	case "endfield":
		return endfield.Claim(cfg, account)
	case "genshin", "starrail", "honkai", "zzz", "themis":
		return hoyo.Claim(cfg, account)
	default:
		return nil, fmt.Errorf("unknown game type: %s", account.Type)
	}
}
