package cmd

import (
	"context"

	"github.com/MisLink/go-web-template/pkg/utils"
	"github.com/spf13/cobra"
)

// crontabCmd represents the crontab command
var crontabCmd = &cobra.Command{
	Use: "crontab",
	RunE: func(cmd *cobra.Command, args []string) error {
		crontab, cleanup, err := CreateCrontab()
		if err != nil {
			return err
		}
		defer cleanup()
		ctx, cancel, err := crontab.Lock(cmd.Context())
		if err != nil {
			return err
		}
		defer cancel()
		return utils.Lifecycle(ctx, func(ctx context.Context) error {
			return crontab.Start()
		}, func() error {
			return crontab.Close()
		})
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}
