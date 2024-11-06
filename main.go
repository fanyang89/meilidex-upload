package main

import (
	"os"

	"github.com/fanyang89/zerologging/v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fanyang89/meilidex-upload/cmd"
)

func main() {
	zerologging.WithConsoleLog(zerolog.InfoLevel)

	app := cmd.NewApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Unexpected error")
	}
}
