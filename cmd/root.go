package cmd

import (
	"fmt"

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

func init() {}

func Execute() error {
	return rootCmd.Execute()
}
