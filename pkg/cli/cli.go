package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/atomicptr/pity-patrol/pkg/config"
	"github.com/atomicptr/pity-patrol/pkg/meta"
	"github.com/atomicptr/pity-patrol/pkg/runner"
	"github.com/atomicptr/pity-patrol/pkg/runner/scheduler"
	"github.com/atomicptr/pity-patrol/pkg/util"
)

func Run() error {
	log.Printf("Pity Patrol %s\n", meta.VersionString())

	if os.Getenv("GITHUB_ACTIONS") != "" || os.Getenv("GITLAB_CI") != "" {
		return fmt.Errorf("unauthorized environment")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.DebugMode {
		log.Printf("[DEBUG] Config: %s\n", util.ToPrettyString(cfg))
	}

	log.Printf("%d account/s configured", len(cfg.Accounts))

	if cfg.EnableScheduler {
		scheduler.Run(cfg)
		return nil
	}

	runner.Run(cfg)
	return nil
}
