package runner

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/atomicptr/pity-patrol/pkgs/claimer"
	"github.com/atomicptr/pity-patrol/pkgs/config"
	"github.com/atomicptr/pity-patrol/pkgs/report/reporter"
)

func Run(cfg *config.Config) {
	for index, account := range cfg.Accounts {
		RunAccount(cfg, index, &account)

		SleepMs(500, 2000)
	}
}

func RunAccount(cfg *config.Config, index int, account *config.Account) {
	identifier := AccountIdentifier(account, index)

	log.Printf("%s claiming...", identifier)

	rep, err := claimer.Claim(cfg, account)
	if err != nil {
		message := fmt.Sprintf("%s could not claim rewards because: %s", identifier, err)
		log.Println(message)

		err := reporter.SendError(cfg, account, message)
		if err != nil {
			log.Printf("error when sending error: %s", err)
		}
		return
	}

	err = reporter.Send(cfg, account, rep)
	if err != nil {
		log.Printf("error when sending report: %s", err)
	}

	if !rep.WasClaimed {
		log.Printf("%s has already claimed reward\n", identifier)
		return
	}

	log.Printf("%s claimed reward successfully\n", identifier)
}

func AccountIdentifier(account *config.Account, index int) string {
	accountName := ""

	if account.Identifier != "" {
		accountName = fmt.Sprintf(" %s", account.Identifier)
	}

	return fmt.Sprintf("Account #%d%s [%s]", index+1, accountName, account.GameName())
}

func SleepMs(from, to int) {
	if to <= from {
		log.Panicf("SleepMs: `to` (%d) must be bigger than `from` (%d)\n", to, from)
	}

	delay := from + rand.Intn(to-from)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
