package cmd

import (
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
		return crontab.Start()
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}
