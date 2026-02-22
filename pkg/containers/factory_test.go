package containers

import (
	"testing"
)

func TestDatabaseFactory(t *testing.T) {
	tests := []struct {
		name    string
		dbType  string
		wantErr bool
	}{
		{"PostgreSQL", "postgres", false},
		{"PostgreSQL alias", "postgresql", false},
		{"MySQL", "mysql", false},
		{"MariaDB", "mariadb", false},
		{"SQLite", "sqlite", false},
		{"DuckDB", "duckdb", false},
		{"ClickHouse", "clickhouse", false},
		{"Invalid type", "invalid", true},
		{"Empty type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := DatabaseFactory(tt.dbType)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatabaseFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && db == nil {
				t.Errorf("DatabaseFactory() returned nil database")
			}
		})
	}
}

func TestPostgreSQL(t *testing.T) {
	db := NewPostgreSQL()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "postgres:latest"},
		{"Latest version explicit", "latest", "postgres:latest"},
		{"Specific version", "15", "postgres:15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetImage(tt.version); got != tt.want {
				t.Errorf("GetImage() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test default port
	if got := db.GetDefaultPort(); got != 5432 {
		t.Errorf("GetDefaultPort() = %v, want %v", got, 5432)
	}

	// Test environment variables
	config := ContainerConfig{Password: "testpass"}
	env := db.GetEnvironment(config)
	if env["POSTGRES_PASSWORD"] != "testpass" {
		t.Errorf("Environment password not set correctly")
	}

	// Test connection info
	config2 := ContainerConfig{Password: "testpass", Port: 5432}
	connInfo := db.GetConnectionInfo(config2, "test")
	if connInfo.User != "postgres" {
		t.Errorf("ConnectionInfo.User = %v, want %v", connInfo.User, "postgres")
	}
}

func TestMySQL(t *testing.T) {
	db := NewMySQL()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "mysql:latest"},
		{"Latest version explicit", "latest", "mysql:latest"},
		{"Specific version", "8.0", "mysql:8.0"},
		{"Another version", "8.4", "mysql:8.4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetImage(tt.version); got != tt.want {
				t.Errorf("GetImage() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test default port
	if got := db.GetDefaultPort(); got != 3306 {
		t.Errorf("GetDefaultPort() = %v, want %v", got, 3306)
	}

	// Test environment variables
	config := ContainerConfig{Password: "testpass", Database: "mydb", User: "myuser"}
	env := db.GetEnvironment(config)
	if env["MYSQL_ROOT_PASSWORD"] != "testpass" {
		t.Errorf("MYSQL_ROOT_PASSWORD not set correctly")
	}
	if env["MYSQL_DATABASE"] != "mydb" {
		t.Errorf("MYSQL_DATABASE not set correctly, got %v", env["MYSQL_DATABASE"])
	}
	if env["MYSQL_USER"] != "myuser" {
		t.Errorf("MYSQL_USER not set correctly, got %v", env["MYSQL_USER"])
	}

	// Test defaults
	configDefault := ContainerConfig{Password: "testpass"}
	envDefault := db.GetEnvironment(configDefault)
	if envDefault["MYSQL_DATABASE"] != "testdb" {
		t.Errorf("MYSQL_DATABASE default not set correctly, got %v", envDefault["MYSQL_DATABASE"])
	}
	if envDefault["MYSQL_USER"] != "testuser" {
		t.Errorf("MYSQL_USER default not set correctly, got %v", envDefault["MYSQL_USER"])
	}

	// Test connection info
	config2 := ContainerConfig{Password: "testpass", Port: 3306, Database: "testdb"}
	connInfo := db.GetConnectionInfo(config2, "test")
	if connInfo.User != "root" {
		t.Errorf("ConnectionInfo.User = %v, want %v", connInfo.User, "root")
	}
	if connInfo.Port != 3306 {
		t.Errorf("ConnectionInfo.Port = %v, want %v", connInfo.Port, 3306)
	}
	if connInfo.Database != "testdb" {
		t.Errorf("ConnectionInfo.Database = %v, want %v", connInfo.Database, "testdb")
	}
}

func TestMariaDB(t *testing.T) {
	db := NewMariaDB()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "mariadb:latest"},
		{"Latest version explicit", "latest", "mariadb:latest"},
		{"Specific version", "11", "mariadb:11"},
		{"Another version", "10.11", "mariadb:10.11"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetImage(tt.version); got != tt.want {
				t.Errorf("GetImage() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test default port
	if got := db.GetDefaultPort(); got != 3306 {
		t.Errorf("GetDefaultPort() = %v, want %v", got, 3306)
	}

	// Test environment variables
	config := ContainerConfig{Password: "testpass", Database: "mydb", User: "myuser"}
	env := db.GetEnvironment(config)
	if env["MARIADB_ROOT_PASSWORD"] != "testpass" {
		t.Errorf("MARIADB_ROOT_PASSWORD not set correctly")
	}
	if env["MARIADB_DATABASE"] != "mydb" {
		t.Errorf("MARIADB_DATABASE not set correctly, got %v", env["MARIADB_DATABASE"])
	}
	if env["MARIADB_USER"] != "myuser" {
		t.Errorf("MARIADB_USER not set correctly, got %v", env["MARIADB_USER"])
	}

	// Test defaults
	configDefault := ContainerConfig{Password: "testpass"}
	envDefault := db.GetEnvironment(configDefault)
	if envDefault["MARIADB_DATABASE"] != "testdb" {
		t.Errorf("MARIADB_DATABASE default not set correctly, got %v", envDefault["MARIADB_DATABASE"])
	}
	if envDefault["MARIADB_USER"] != "testuser" {
		t.Errorf("MARIADB_USER default not set correctly, got %v", envDefault["MARIADB_USER"])
	}

	// Test connection info
	config2 := ContainerConfig{Password: "testpass", Port: 3306, Database: "testdb"}
	connInfo := db.GetConnectionInfo(config2, "test")
	if connInfo.User != "root" {
		t.Errorf("ConnectionInfo.User = %v, want %v", connInfo.User, "root")
	}
	if connInfo.Port != 3306 {
		t.Errorf("ConnectionInfo.Port = %v, want %v", connInfo.Port, 3306)
	}
	if connInfo.Database != "testdb" {
		t.Errorf("ConnectionInfo.Database = %v, want %v", connInfo.Database, "testdb")
	}
}

func TestSQLite(t *testing.T) {
	db := NewSQLite()

	// Test image (SQLite doesn't use containers)
	if got := db.GetImage("latest"); got != "" {
		t.Errorf("GetImage() = %v, want empty string for SQLite", got)
	}

	// Test default port (SQLite doesn't use ports)
	if got := db.GetDefaultPort(); got != 0 {
		t.Errorf("GetDefaultPort() = %v, want 0 for SQLite", got)
	}

	// Test environment variables (should be empty)
	config := ContainerConfig{Password: "testpass"}
	env := db.GetEnvironment(config)
	if len(env) != 0 {
		t.Errorf("GetEnvironment() should return empty map for SQLite, got %v", env)
	}

	// Test connection info
	config2 := ContainerConfig{Database: "mydb"}
	connInfo := db.GetConnectionInfo(config2, "mydb")
	if connInfo.Database != "mydb.db" {
		t.Errorf("ConnectionInfo.Database = %v, want mydb.db", connInfo.Database)
	}
	if connInfo.DSN != "mydb.db" {
		t.Errorf("ConnectionInfo.DSN = %v, want mydb.db", connInfo.DSN)
	}

	// Test with container name fallback
	config3 := ContainerConfig{}
	connInfo3 := db.GetConnectionInfo(config3, "fallback")
	if connInfo3.Database != "fallback.db" {
		t.Errorf("ConnectionInfo.Database with fallback = %v, want fallback.db", connInfo3.Database)
	}
}

func TestDuckDB(t *testing.T) {
	db := NewDuckDB()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "duckdb/duckdb:latest"},
		{"Latest version explicit", "latest", "duckdb/duckdb:latest"},
		{"Specific version", "0.9", "duckdb/duckdb:0.9"},
		{"Another version", "0.10", "duckdb/duckdb:0.10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetImage(tt.version); got != tt.want {
				t.Errorf("GetImage() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test default port
	if got := db.GetDefaultPort(); got != 0 {
		t.Errorf("GetDefaultPort() = %v, want 0 for DuckDB", got)
	}

	// Test environment variables (should be empty)
	config := ContainerConfig{Password: "testpass"}
	env := db.GetEnvironment(config)
	if len(env) != 0 {
		t.Errorf("GetEnvironment() should return empty map for DuckDB, got %v", env)
	}

	// Test connection info
	config2 := ContainerConfig{Database: "mydb"}
	connInfo := db.GetConnectionInfo(config2, "mydb")
	if connInfo.Database != "mydb" {
		t.Errorf("ConnectionInfo.Database = %v, want mydb", connInfo.Database)
	}
	if connInfo.DSN != "duckdb://mydb.duckdb" {
		t.Errorf("ConnectionInfo.DSN = %v, want duckdb://mydb.duckdb", connInfo.DSN)
	}

	// Test with container name fallback
	config3 := ContainerConfig{}
	connInfo3 := db.GetConnectionInfo(config3, "fallback")
	if connInfo3.Database != "fallback" {
		t.Errorf("ConnectionInfo.Database with fallback = %v, want fallback", connInfo3.Database)
	}
}

func TestClickHouse(t *testing.T) {
	db := NewClickHouse()

	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "clickhouse/clickhouse-server:latest"},
		{"Latest version explicit", "latest", "clickhouse/clickhouse-server:latest"},
		{"Specific version", "24", "clickhouse/clickhouse-server:24"},
		{"Another version", "23.8", "clickhouse/clickhouse-server:23.8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := db.GetImage(tt.version); got != tt.want {
				t.Errorf("GetImage() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test default port
	if got := db.GetDefaultPort(); got != 8123 {
		t.Errorf("GetDefaultPort() = %v, want %v", got, 8123)
	}

	// Test environment variables
	config := ContainerConfig{Password: "testpass", Database: "mydb", User: "myuser"}
	env := db.GetEnvironment(config)
	if env["CLICKHOUSE_PASSWORD"] != "testpass" {
		t.Errorf("CLICKHOUSE_PASSWORD not set correctly")
	}
	if env["CLICKHOUSE_USER"] != "myuser" {
		t.Errorf("CLICKHOUSE_USER not set correctly, got %v", env["CLICKHOUSE_USER"])
	}
	if env["CLICKHOUSE_DB"] != "mydb" {
		t.Errorf("CLICKHOUSE_DB not set correctly, got %v", env["CLICKHOUSE_DB"])
	}

	// Test defaults
	configDefault := ContainerConfig{Password: "testpass"}
	envDefault := db.GetEnvironment(configDefault)
	if envDefault["CLICKHOUSE_USER"] != "default" {
		t.Errorf("CLICKHOUSE_USER default not set correctly, got %v", envDefault["CLICKHOUSE_USER"])
	}
	if envDefault["CLICKHOUSE_DB"] != "default" {
		t.Errorf("CLICKHOUSE_DB default not set correctly, got %v", envDefault["CLICKHOUSE_DB"])
	}

	// Test connection info
	config2 := ContainerConfig{Password: "testpass", Port: 8123, Database: "mydb", User: "myuser"}
	connInfo := db.GetConnectionInfo(config2, "test")
	if connInfo.User != "myuser" {
		t.Errorf("ConnectionInfo.User = %v, want %v", connInfo.User, "myuser")
	}
	if connInfo.Port != 8123 {
		t.Errorf("ConnectionInfo.Port = %v, want %v", connInfo.Port, 8123)
	}
	if connInfo.Database != "mydb" {
		t.Errorf("ConnectionInfo.Database = %v, want %v", connInfo.Database, "mydb")
	}
}
