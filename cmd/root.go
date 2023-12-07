package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	types "github.com/MisLink/go-web-template/types"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     types.ModuleName,
	Version: fmt.Sprintf("%s built at: %s", types.Version, types.BuiltAt),
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	ctx, reset := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	defer reset()
	return rootCmd.ExecuteContext(ctx)
}
