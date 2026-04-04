package claimer

import (
	"fmt"

	"atomicptr.dev/pity-patrol/pkg/claimer/endfield"
	"atomicptr.dev/pity-patrol/pkg/claimer/hoyo"
	"atomicptr.dev/pity-patrol/pkg/config"
	"atomicptr.dev/pity-patrol/pkg/report"
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
