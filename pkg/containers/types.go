package containers

import (
	"time"
)

// ContainerConfig holds configuration for creating a database container
type ContainerConfig struct {
	Type          string
	Version       string
	Port          int
	Password      string
	Database      string // Database name (defaults to "testdb" if empty)
	User          string // Database user (defaults depend on database type)
	Volume        string
	Network       string
	SQLFile       string
	Name          string // For file-based databases (SQLite/DuckDB)
	ContainerName string // Docker container name (auto-generated if empty)
}

// ContainerInfo holds information about a running container
type ContainerInfo struct {
	ID       string
	Name     string
	Type     string
	Version  string
	Status   string
	Port     int
	User     string
	Password string
	Database string
	DSN      string // Connection string for the database
	Volume   string
	Created  time.Time
}

// Database represents a database type with its configuration
type Database interface {
	// GetImage returns the Docker image name and tag
	GetImage(version string) string

	// GetDefaultPort returns the default port for this database
	GetDefaultPort() int

	// GetEnvironment returns environment variables for the container
	GetEnvironment(config ContainerConfig) map[string]string

	// GetConnectionInfo returns connection information for this database
	GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo

	// ExecuteSQL executes SQL commands (for initialization)
	ExecuteSQL(containerID string, sqlFile string) error
}

// ConnectionInfo holds database connection details
type ConnectionInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	DSN      string
}
