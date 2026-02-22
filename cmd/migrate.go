package cmd

import (
	"fmt"

	"github.com/dbz/dbz/pkg/migrations"
	"github.com/spf13/cobra"
)

var (
	direction string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [database] [migration-file]",
	Short: "Run database migrations",
	Long: `Run database migrations from SQL files.

Examples:
  dbz migrate postgres init.sql              # Run migration file
  dbz migrate myapp schema.sql               # Run schema migration
  dbz migrate myapp seed.sql --direction up  # Run seed data migration
`,
	Args: cobra.ExactArgs(2),
	RunE: runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringVarP(&direction, "direction", "d", "up", "Migration direction (up/down)")
}

func runMigrate(cmd *cobra.Command, args []string) error {
	database := args[0]
	migrationFile := args[1]

	config := migrations.MigrationConfig{
		File:      migrationFile,
		Direction: direction,
		Database:  database,
	}

	err := migrations.RunMigration(config)
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	fmt.Printf("✅ Applied migration: %s to %s\n", migrationFile, database)

	return nil
}
