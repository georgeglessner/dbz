package containers

import (
	"fmt"
)

// ClickHouse implements the Database interface for ClickHouse
type ClickHouse struct{}

// NewClickHouse creates a new ClickHouse database instance
func NewClickHouse() Database {
	return &ClickHouse{}
}

func (c *ClickHouse) GetImage(version string) string {
	if version == "" || version == "latest" {
		return "clickhouse/clickhouse-server:latest"
	}
	return fmt.Sprintf("clickhouse/clickhouse-server:%s", version)
}

func (c *ClickHouse) GetDefaultPort() int {
	return 8123 // HTTP port
}

func (c *ClickHouse) GetEnvironment(config ContainerConfig) map[string]string {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "default"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "default"
	}

	return map[string]string{
		"CLICKHOUSE_USER":     dbUser,
		"CLICKHOUSE_PASSWORD": config.Password,
		"CLICKHOUSE_DB":       dbName,
	}
}

func (c *ClickHouse) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "default"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "default"
	}

	return ConnectionInfo{
		Host:     "localhost",
		Port:     config.Port,
		User:     dbUser,
		Password: config.Password,
		Database: dbName,
		DSN:      fmt.Sprintf("clickhouse://%s:%s@localhost:%d/%s", dbUser, config.Password, config.Port, dbName),
	}
}

func (c *ClickHouse) ExecuteSQL(containerID string, sqlFile string) error {
	// ClickHouse-specific SQL execution logic
	return fmt.Errorf("ClickHouse SQL execution not yet implemented")
}
