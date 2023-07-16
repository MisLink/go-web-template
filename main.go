package main

import (
	"log"

	"MODULE_NAME/cmd"

	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(s string, i ...any) {}))
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
