package main

import (
	"os"

	"isp_checker/cmd/blacklist/cmd"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	if err = cmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
