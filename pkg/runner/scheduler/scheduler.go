package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/atomicptr/pity-patrol/pkg/config"
	"github.com/atomicptr/pity-patrol/pkg/runner"
	"github.com/netresearch/go-cron"
)

func Run(cfg *config.Config) {
	log.Println("Running in scheduler mode...")

	if len(cfg.Accounts) == 0 {
		log.Println("No accounts have been registered, exiting...")
		return
	}

	c := cron.New()

	for index, account := range cfg.Accounts {
		r := config.ResetTimeByAccountType(account.Type)

		if account.CheckinOffset > 0 {
			r.Add(account.CheckinOffset)
		}

		tzString := ""
		if r.TimeZone != "" {
			tzString = fmt.Sprintf("CRON_TZ=%s ", r.TimeZone)
		}

		cronString := fmt.Sprintf("%s%d %d * * *", tzString, r.Minute, r.Hour)

		_, err := c.AddFunc(cronString, func() {
			runner.SleepMs(5000, 60_000) // sleep 5s - 60s
			runner.RunAccount(cfg, index, &account)
		})

		if err != nil {
			log.Panicf("Could not create cron job for %s (Cron: `%s`) because: %s", runner.AccountIdentifier(&account, index), cronString, err)
			return
		}

		if cfg.DebugMode {
			ident := runner.AccountIdentifier(&account, index)
			log.Printf("[DEBUG] %s registered job at %s", ident, cronString)
		}
	}

	go startHeartbeat(c)

	c.Start()

	select {}
}

func startHeartbeat(c *cron.Cron) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		entries := c.Entries()

		if len(entries) == 0 {
			log.Println("Scheduler: No jobs found.")
			continue
		}

		var next *cron.Entry
		for _, entry := range entries {
			if next == nil || entry.Next.Before(next.Next) {
				next = &entry
			}
		}

		if next != nil {
			log.Printf("Scheduler: Next scheduled job runs at %s", next.Next.Format("2006-01-02 15:04:05"))
		}
	}
}
