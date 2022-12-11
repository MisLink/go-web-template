package cmd

import (
	"context"
	"os"

	"MODULE_NAME/pkg/database/ent/migrate"

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
		opts := []schema.MigrateOption{migrate.WithForeignKeys(withForeignKeys)}
		if delete {
			opts = append(opts, migrate.WithDropColumn(true), migrate.WithDropIndex(true))
		}
		if dryRun {
			return client.Schema.WriteTo(context.Background(), os.Stdout, opts...)
		} else {
			return client.Schema.Create(context.Background(), opts...)
		}
	},
}

var (
	dryRun          bool
	delete          bool
	withForeignKeys bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "display migrate sql")
	migrateCmd.Flags().BoolVar(&delete, "delete", false, "generate sql with drop")
	migrateCmd.Flags().BoolVar(&withForeignKeys, "with-foreign-keys", false, "create foreign keys")
}
