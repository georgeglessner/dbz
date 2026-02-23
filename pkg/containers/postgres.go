package containers

import (
	"fmt"
)

// PostgreSQL implements the Database interface for PostgreSQL
type PostgreSQL struct{}

// NewPostgreSQL creates a new PostgreSQL database instance
func NewPostgreSQL() Database {
	return &PostgreSQL{}
}

func (p *PostgreSQL) GetImage(version string) string {
	if version == "" || version == "latest" {
		return "postgres:latest"
	}
	return fmt.Sprintf("postgres:%s", version)
}

func (p *PostgreSQL) GetDefaultPort() int {
	return 5432
}

func (p *PostgreSQL) GetEnvironment(config ContainerConfig) map[string]string {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "postgres"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "postgres"
	}

	return map[string]string{
		"POSTGRES_PASSWORD": config.Password,
		"POSTGRES_USER":     dbUser,
		"POSTGRES_DB":       dbName,
	}
}

func (p *PostgreSQL) GetDataPath() string {
	return "/var/lib/postgresql"
}

func (p *PostgreSQL) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "postgres"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "postgres"
	}

	return ConnectionInfo{
		Host:     "localhost",
		Port:     config.Port,
		User:     dbUser,
		Password: config.Password,
		Database: dbName,
		DSN:      fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", dbUser, config.Password, config.Port, dbName),
	}
}

func (p *PostgreSQL) ExecuteSQL(containerID string, sqlFile string) error {
	// PostgreSQL-specific SQL execution logic
	// This would involve copying the SQL file to the container and running psql
	return fmt.Errorf("PostgreSQL SQL execution not yet implemented")
}
