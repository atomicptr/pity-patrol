package runner

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/atomicptr/pity-patrol/pkgs/claimer"
	"github.com/atomicptr/pity-patrol/pkgs/config"
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

	claimed, err := claimer.Claim(cfg, account)
	if err != nil {
		log.Printf("%s could not claim rewards because: %s\n", identifier, err)
		return
	}

	if !claimed {
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
	if from <= to {
		panic("SleepMs: `from` must be bigger than `to`")
	}

	delay := from + rand.Intn(to-from)
	time.Sleep(time.Duration(delay) * time.Millisecond)

}
