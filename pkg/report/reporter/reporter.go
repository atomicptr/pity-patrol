package reporter

import (
	"log"

	"github.com/atomicptr/pity-patrol/pkg/config"
	"github.com/atomicptr/pity-patrol/pkg/report"
	"github.com/atomicptr/pity-patrol/pkg/report/discord"
)

func Send(cfg *config.Config, account *config.Account, r *report.Report) error {
	if r == nil {
		return nil
	}

	for _, reporter := range cfg.Reporters {
		if !reporter.ReportOnSuccess() {
			continue
		}

		var err error

		switch reporter.Type {
		case "discord":
			err = discord.Send(&reporter, cfg, account, r)
		default:
			log.Printf("invalid reporter type: %s", reporter.Type)
			continue
		}

		if err != nil {
			log.Printf("error while trying to send to '%s': %s", reporter.Type, err)
		}
	}

	return nil
}

func SendError(cfg *config.Config, account *config.Account, message string) error {
	for _, reporter := range cfg.Reporters {
		if !reporter.ReportOnFailure() {
			continue
		}

		var err error

		switch reporter.Type {
		case "discord":
			err = discord.SendError(&reporter, cfg, account, message)
		default:
			log.Printf("invalid reporter type: %s", reporter.Type)
			continue
		}

		if err != nil {
			log.Printf("error while trying to send to '%s': %s", reporter.Type, err)
		}
	}

	return nil
}
