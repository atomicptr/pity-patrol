package cli

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/atomicptr/pity-patrol/pkgs/claimer"
	"github.com/atomicptr/pity-patrol/pkgs/config"
	"github.com/atomicptr/pity-patrol/pkgs/meta"
	"github.com/atomicptr/pity-patrol/pkgs/util"
)

func Run() error {
	log.Printf("Pity Patrol %s\n", meta.VersionString())

	if os.Getenv("GITHUB_ACTIONS") != "" || os.Getenv("GITLAB_CI") != "" {
		return fmt.Errorf("Unauthorized environment")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.DebugMode {
		log.Printf("[DEBUG] Config: %s\n", util.ToPrettyString(cfg))
	}

	log.Printf("%d account/s configured", len(cfg.Accounts))

	for index, account := range cfg.Accounts {
		accountName := ""

		if account.Identifier != "" {
			accountName = fmt.Sprintf(" %s", account.Identifier)
		}

		identifier := fmt.Sprintf("Account #%d%s [%s]", index+1, accountName, account.GameName())

		log.Printf("%s claiming...", identifier)

		claimed, err := claimer.Claim(cfg, &account)
		if err != nil {
			log.Printf("%s could not claim rewards because: %s\n", identifier, err)
			continue
		}

		if !claimed {
			log.Printf("%s has already claimed reward\n", identifier)
			continue
		}

		log.Printf("%s claimed reward successfully\n", identifier)

		delay := 500 + rand.Intn(1500)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	return nil
}
