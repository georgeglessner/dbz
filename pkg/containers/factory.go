package containers

import (
	"fmt"
	"strings"
)

// DatabaseFactory creates database instances based on type
func DatabaseFactory(dbType string) (Database, error) {
	switch strings.ToLower(dbType) {
	case "postgres", "postgresql":
		return NewPostgreSQL(), nil
	case "mysql":
		return NewMySQL(), nil
	case "mariadb":
		return NewMariaDB(), nil
	case "sqlite":
		return NewSQLite(), nil
	case "duckdb":
		return NewDuckDB(), nil
	case "clickhouse":
		return NewClickHouse(), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
