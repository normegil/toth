package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func Init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Configure(colored bool) {
	if colored {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
}
