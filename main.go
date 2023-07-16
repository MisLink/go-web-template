package main

import (
	"MODULE_NAME/cmd"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(s string, i ...any) {}))
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
