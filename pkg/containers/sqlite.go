package containers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite implements the Database interface for SQLite (file-based, no container)
type SQLite struct{}

// NewSQLite creates a new SQLite database instance
func NewSQLite() Database {
	return &SQLite{}
}

func (s *SQLite) GetImage(version string) string {
	return "" // No container image for SQLite
}

func (s *SQLite) GetDefaultPort() int {
	return 0 // No port for SQLite
}

func (s *SQLite) GetEnvironment(config ContainerConfig) map[string]string {
	return map[string]string{} // No environment variables for SQLite
}

func (s *SQLite) GetDataPath() string {
	return "" // SQLite is file-based, no container volume needed
}

func (s *SQLite) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided database name or container name
	dbName := config.Database
	if dbName == "" {
		dbName = containerName
	}

	dbPath := fmt.Sprintf("%s.db", dbName)
	return ConnectionInfo{
		Host:     "",
		Port:     0,
		User:     "",
		Password: "",
		Database: dbPath,
		DSN:      dbPath,
	}
}

func (s *SQLite) ExecuteSQL(containerID string, sqlFile string) error {
	// For SQLite, containerID is actually the database file path
	dbPath := containerID

	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Read SQL file
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	// Execute SQL
	if _, err := db.Exec(string(sqlContent)); err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
