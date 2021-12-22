package logging

import (
	"github.com/rs/zerolog/log"
)

func Setup() {
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
}
