package containers

import (
	"fmt"
)

// MySQL implements the Database interface for MySQL
type MySQL struct{}

// NewMySQL creates a new MySQL database instance
func NewMySQL() Database {
	return &MySQL{}
}

func (m *MySQL) GetImage(version string) string {
	if version == "" || version == "latest" {
		return "mysql:latest"
	}
	return fmt.Sprintf("mysql:%s", version)
}

func (m *MySQL) GetDefaultPort() int {
	return 3306
}

func (m *MySQL) GetEnvironment(config ContainerConfig) map[string]string {
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
		"MYSQL_ROOT_PASSWORD": config.Password,
		"MYSQL_DATABASE":      dbName,
		"MYSQL_USER":          dbUser,
		"MYSQL_PASSWORD":      config.Password,
		// Allow root connections from any host
		"MYSQL_ROOT_HOST": "%",
	}
}

func (m *MySQL) GetConnectionInfo(config ContainerConfig, containerName string) ConnectionInfo {
	// Use provided values or defaults
	dbName := config.Database
	if dbName == "" {
		dbName = "testdb"
	}

	// DSN for Go applications (go-sql-driver/mysql)
	goDSN := fmt.Sprintf("root:%s@tcp(localhost:%d)/%s?tls=skip-verify", config.Password, config.Port, dbName)

	return ConnectionInfo{
		Host:     "localhost",
		Port:     config.Port,
		User:     "root",
		Password: config.Password,
		Database: dbName,
		DSN:      goDSN,
		// Note: For GUI clients like DBeaver, use:
		// jdbc:mysql://localhost:%d/<database>?allowPublicKeyRetrieval=true&useSSL=false
	}
}

func (m *MySQL) ExecuteSQL(containerID string, sqlFile string) error {
	// MySQL-specific SQL execution logic
	// This would involve copying the SQL file to the container and running mysql
	return fmt.Errorf("MySQL SQL execution not yet implemented")
}
