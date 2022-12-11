package cmd

import (
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
		return server.Start()
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)
}
