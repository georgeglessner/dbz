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

	// Test image naming
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"Latest version", "", "mysql:latest"},
		{"Specific version", "8.0", "mysql:8.0"},
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
}
