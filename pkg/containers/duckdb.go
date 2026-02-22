package containers

import (
	"fmt"
)

// DuckDB implements the Database interface for DuckDB
type DuckDB struct{}

// NewDuckDB creates a new DuckDB database instance
func NewDuckDB() Database {
	return &DuckDB{}
}

func (d *DuckDB) GetImage(version string) string {
	if version == "" || version == "latest" {
		return "duckdb/duckdb:latest"
	}
	return fmt.Sprintf("duckdb/duckdb:%s", version)
}

func (d *DuckDB) GetDefaultPort() int {
	return 0 // DuckDB doesn't expose a port by default
}

func (d *DuckDB) GetEnvironment(config ContainerConfig) map[string]string {
	return map[string]string{} // Minimal environment for DuckDB
}

func (d *DuckDB) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided database name or container name
	dbName := config.Database
	if dbName == "" {
		dbName = containerName
	}

	return ConnectionInfo{
		Host:     "localhost",
		Port:     0,
		User:     "",
		Password: "",
		Database: dbName,
		DSN:      fmt.Sprintf("duckdb://%s.duckdb", dbName),
	}
}

func (d *DuckDB) ExecuteSQL(containerID string, sqlFile string) error {
	// DuckDB-specific SQL execution logic
	return fmt.Errorf("DuckDB SQL execution not yet implemented")
}
