package cmd

import (
	"context"
	"os"

	"github.com/MisLink/go-web-template/pkg/database/ent/migrate"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use: "migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, cleanup, err := CreateDatabase()
		if err != nil {
			return err
		}
		defer cleanup()
		opts := []schema.MigrateOption{migrate.WithForeignKeys(!withOutForeignKeys)}
		if drop {
			opts = append(opts, migrate.WithDropColumn(true), migrate.WithDropIndex(true))
		}
		if dryRun {
			return client.Schema.WriteTo(context.Background(), os.Stdout, opts...)
		}
		return client.Schema.Create(context.Background(), opts...)
	},
}

var (
	dryRun             bool
	drop               bool
	withOutForeignKeys bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "display migrate sql")
	migrateCmd.Flags().BoolVar(&drop, "delete", false, "generate sql with drop")
	migrateCmd.Flags().BoolVar(&withOutForeignKeys, "without-foreign-keys", false, "ignore foreign keys")
}
