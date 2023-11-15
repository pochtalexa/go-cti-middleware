package flags

import (
	"flag"
	"github.com/rs/zerolog/log"
)

var (
	FlagRunAddr  string
	FlagLogin    string
	FlagPassword string
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ParseFlags() {

	defaultLogin := "agent"
	defaultPassword := "agent"

	flag.StringVar(&FlagRunAddr, "a", "localhost:9595", "middleware api addr")
	flag.StringVar(&FlagLogin, "l", defaultLogin, "login")
	flag.StringVar(&FlagPassword, "p", defaultPassword, "password")
	flag.Parse()

	log.Info().Msg("ParseFlags - success")
}
