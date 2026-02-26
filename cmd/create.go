package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dbz/dbz/pkg/containers"
	"github.com/spf13/cobra"
)

var (
	port          int
	password      string
	dbName        string
	dbUser        string
	volume        string
	network       string
	containerName string
)

var createCmd = &cobra.Command{
	Use:   "create [database] [sql-file]",
	Short: "Create a new database container",
	Long: `Create a new database container with the specified type and optional SQL file.

Supported databases:
  - postgres, postgresql
  - mysql, mariadb
  - sqlite (file-based, no container)
  - duckdb
  - clickhouse

Examples:
  dbz create postgres
  dbz create mysql@8.4
  dbz create postgres init.sql
  dbz create mysql --port 3307 --password mypass
  dbz create mysql --database myapp --user appuser --password secret123
  dbz create postgres --database prod_db --user admin
`,
	Args: cobra.RangeArgs(1, 3),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().IntVarP(&port, "port", "p", 0, "Port to expose (auto-assign if not specified)")
	createCmd.Flags().StringVarP(&password, "password", "d", "", "Database password (auto-generate if not specified)")
	createCmd.Flags().StringVarP(&dbName, "database", "b", "", "Database name (default: testdb)")
	createCmd.Flags().StringVarP(&dbUser, "user", "u", "", "Database user (default depends on database type)")
	createCmd.Flags().StringVarP(&volume, "volume", "v", "", "Volume to mount for data persistence")
	createCmd.Flags().StringVarP(&network, "network", "n", "", "Docker network to join")
	createCmd.Flags().StringVarP(&containerName, "name", "c", "", "Docker container name (auto-generated if not specified)")
}

func runCreate(cmd *cobra.Command, args []string) error {
	dbSpec := args[0]
	var sqlFile string
	var nameFromArg string

	if len(args) > 1 {
		// For SQLite and DuckDB, the second argument is the database name
		// For other databases, it's an SQL file
		if strings.ToLower(dbSpec) == "sqlite" || strings.ToLower(dbSpec) == "duckdb" {
			nameFromArg = args[1]
			if len(args) > 2 {
				sqlFile = args[2]
			}
		} else {
			sqlFile = args[1]
		}
	}

	// Use flag value if provided, otherwise use arg value (for SQLite/DuckDB), otherwise empty (will use default)
	finalDbName := dbName
	if finalDbName == "" {
		finalDbName = nameFromArg
	}

	// Parse database type and version
	parts := strings.Split(dbSpec, "@")
	dbType := strings.ToLower(parts[0])
	version := "latest"
	if len(parts) > 1 {
		version = parts[1]
	}

	// Validate database type
	if !isSupportedDatabase(dbType) {
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Check if SQL file exists
	if sqlFile != "" {
		if _, err := os.Stat(sqlFile); os.IsNotExist(err) {
			return fmt.Errorf("SQL file not found: %s", sqlFile)
		}
	}

	// Create database container
	config := containers.ContainerConfig{
		Type:          dbType,
		Version:       version,
		Port:          port,
		Password:      password,
		Database:      finalDbName,
		User:          dbUser,
		Volume:        volume,
		Network:       network,
		SQLFile:       sqlFile,
		Name:          finalDbName,
		ContainerName: containerName,
	}

	container, err := containers.CreateContainer(config)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	fmt.Printf("✅ Created %s %s database\n", container.Name, container.Version)
	fmt.Printf("   Host: localhost:%d\n", container.Port)
	fmt.Printf("   User: %s\n", container.User)
	fmt.Printf("   Password: %s\n", container.Password)
	fmt.Printf("   Database: %s\n", container.Database)

	if container.Volume != "" {
		fmt.Printf("   Volume: %s\n", container.Volume)
	}

	// Show connection strings for MySQL/MariaDB
	if dbType == "mysql" || dbType == "mariadb" {
		fmt.Printf("\n   Connection (Go): %s\n", container.DSN)
		fmt.Printf("   Connection (GUI clients like DBeaver):\n")
		fmt.Printf("      Host: localhost\n")
		fmt.Printf("      Port: %d\n", container.Port)
		fmt.Printf("      Database: %s\n", container.Database)
		fmt.Printf("      User: %s\n", container.User)
		fmt.Printf("      Password: %s\n", container.Password)
		fmt.Printf("      URL: jdbc:mysql://localhost:%d/%s?allowPublicKeyRetrieval=true&useSSL=false\n",
			container.Port, container.Database)
	}

	return nil
}

func isSupportedDatabase(dbType string) bool {
	supported := []string{
		"postgres", "postgresql",
		"mysql", "mariadb",
		"sqlite",
		"duckdb",
		"clickhouse",
	}

	for _, supported := range supported {
		if dbType == supported {
			return true
		}
	}
	return false
}
