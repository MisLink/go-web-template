package main

import (
	"math/rand"
	"time"

	"MODULE_NAME/cmd"

	"go.uber.org/automaxprocs/maxprocs"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(s string, i ...any) {}))
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cmd.Execute()
}
