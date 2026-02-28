package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dbz",
	Short: "dbz - quickly and easily create databases from the command line",
	Long: `dbz is a CLI tool that creates databases as Docker containers.
It supports multiple database types including PostgreSQL, MySQL, MariaDB, SQLite, DuckDB, and ClickHouse.

Examples:
  dbz create postgres              # Create PostgreSQL database with latest version
  dbz create mysql@8.4            # Create MySQL database with specific version
  dbz create postgres db_file.sql  # Create PostgreSQL and run SQL file
  dbz list                         # List all running database containers
  dbz delete postgres              # Delete PostgreSQL container
`,
}

func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the version for the root command
func SetVersion(version string) {
	rootCmd.Version = version
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}
