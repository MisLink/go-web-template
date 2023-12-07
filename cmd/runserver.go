package cmd

import (
	"context"

	"github.com/MisLink/go-web-template/pkg/utils"
	"github.com/spf13/cobra"
)

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use: "runserver",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, cleanup, err := CreateServer()
		if err != nil {
			return err
		}
		defer cleanup()
		return utils.Lifecycle(cmd.Context(), func(ctx context.Context) error {
			return server.Start(ctx)
		}, func() error {
			return server.Close()
		})
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
}
