package containers

import (
	"fmt"
)

// MariaDB implements the Database interface for MariaDB
type MariaDB struct{}

// NewMariaDB creates a new MariaDB database instance
func NewMariaDB() Database {
	return &MariaDB{}
}

func (m *MariaDB) GetImage(version string) string {
	if version == "" || version == "latest" {
		return "mariadb:latest"
	}
	return fmt.Sprintf("mariadb:%s", version)
}

func (m *MariaDB) GetDefaultPort() int {
	return 3306
}

func (m *MariaDB) GetEnvironment(config ContainerConfig) map[string]string {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "testdb"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "testuser"
	}

	return map[string]string{
		"MARIADB_ROOT_PASSWORD": config.Password,
		"MARIADB_DATABASE":      dbName,
		"MARIADB_USER":          dbUser,
		"MARIADB_PASSWORD":      config.Password,
		// Allow root connections from any host
		"MARIADB_ROOT_HOST": "%",
	}
}

func (m *MariaDB) GetDataPath() string {
	return "/var/lib/mysql"
}

func (m *MariaDB) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "testdb"
	}

	// DSN with tls=skip-verify to handle self-signed certificates in Docker
	primaryDSN := fmt.Sprintf("root:%s@tcp(localhost:%d)/%s?tls=skip-verify", config.Password, config.Port, dbName)

	return ConnectionInfo{
		Host:     "localhost",
		Port:     config.Port,
		User:     "root",
		Password: config.Password,
		Database: dbName,
		DSN:      primaryDSN,
	}
}

func (m *MariaDB) ExecuteSQL(containerID string, sqlFile string) error {
	// MariaDB-specific SQL execution logic
	// This would involve copying the SQL file to the container and running mysql
	return fmt.Errorf("MariaDB SQL execution not yet implemented")
}
