package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manager handles container operations
type Manager struct {
	docker *DockerClient
}

// NewManager creates a new container manager
func NewManager() (*Manager, error) {
	docker, err := NewDockerClient()
	if err != nil {
		return nil, err
	}

	return &Manager{docker: docker}, nil
}

// Close closes the manager and its resources
func (m *Manager) Close() error {
	if m.docker != nil {
		return m.docker.Close()
	}
	return nil
}

// CreateContainer creates a new database container
func CreateContainer(config ContainerConfig) (*ContainerInfo, error) {
	// Handle SQLite and DuckDB separately (no container needed - file-based)
	if config.Type == "sqlite" {
		return createSQLiteDatabase(config)
	}
	if config.Type == "duckdb" {
		return createDuckDBDatabase(config)
	}

	manager, err := NewManager()
	if err != nil {
		return nil, err
	}
	defer manager.Close()

	// Get database implementation
	db, err := DatabaseFactory(config.Type)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return manager.docker.CreateContainer(ctx, config, db)
}

// createSQLiteDatabase creates a SQLite database file
func createSQLiteDatabase(config ContainerConfig) (*ContainerInfo, error) {
	db, err := DatabaseFactory(config.Type)
	if err != nil {
		return nil, err
	}

	// Use provided name or generate one
	dbName := config.Name
	if dbName == "" {
		dbName = generateContainerName(config.Type)
	}
	dbPath := fmt.Sprintf("%s.db", dbName)

	// Create directory if needed
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create the database file
	file, err := os.Create(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()

	// Execute SQL file if provided
	if config.SQLFile != "" {
		if err := db.ExecuteSQL(dbPath, config.SQLFile); err != nil {
			// Clean up on error
			os.Remove(dbPath)
			return nil, fmt.Errorf("failed to execute SQL file: %w", err)
		}
	}

	// Update config with the actual dbName for GetConnectionInfo
	config.Name = dbName
	config.Database = dbName
	connInfo := db.GetConnectionInfo(config, dbName)

	return &ContainerInfo{
		ID:       dbPath,
		Name:     dbName,
		Type:     config.Type,
		Version:  config.Version,
		Status:   "file",
		Port:     0,
		User:     connInfo.User,
		Password: "",
		Database: dbPath,
		DSN:      connInfo.DSN,
		Volume:   "",
		Created:  time.Now(),
	}, nil
}

// createDuckDBDatabase creates a DuckDB database file
func createDuckDBDatabase(config ContainerConfig) (*ContainerInfo, error) {
	db, err := DatabaseFactory(config.Type)
	if err != nil {
		return nil, err
	}

	// Use provided name or generate one
	dbName := config.Name
	if dbName == "" {
		dbName = generateContainerName(config.Type)
	}
	dbPath := fmt.Sprintf("%s.duckdb", dbName)

	// Create directory if needed
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create the database file
	file, err := os.Create(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()

	// Execute SQL file if provided
	if config.SQLFile != "" {
		if err := db.ExecuteSQL(dbPath, config.SQLFile); err != nil {
			// Clean up on error
			os.Remove(dbPath)
			return nil, fmt.Errorf("failed to execute SQL file: %w", err)
		}
	}

	// Update config with the actual dbName for GetConnectionInfo
	config.Name = dbName
	config.Database = dbName
	connInfo := db.GetConnectionInfo(config, dbName)

	return &ContainerInfo{
		ID:       dbPath,
		Name:     dbName,
		Type:     config.Type,
		Version:  config.Version,
		Status:   "file",
		Port:     0,
		User:     connInfo.User,
		Password: "",
		Database: dbPath,
		DSN:      connInfo.DSN,
		Volume:   "",
		Created:  time.Now(),
	}, nil
}

// DeleteContainer removes a container by name or ID
func DeleteContainer(containerName string, removeVolumes bool) error {
	// Handle SQLite separately - check if it's a .db file or if we should look for one
	if filepath.Ext(containerName) == ".db" {
		return deleteSQLiteDatabase(containerName)
	}

	// Check if it's a SQLite database name (without .db extension)
	if _, err := os.Stat(containerName + ".db"); err == nil {
		return deleteSQLiteDatabase(containerName)
	}

	// Handle DuckDB separately - check if it's a .duckdb file
	if filepath.Ext(containerName) == ".duckdb" {
		return deleteDuckDBDatabase(containerName)
	}

	// Check if it's a DuckDB database name (without .duckdb extension)
	if _, err := os.Stat(containerName + ".duckdb"); err == nil {
		return deleteDuckDBDatabase(containerName)
	}

	manager, err := NewManager()
	if err != nil {
		return err
	}
	defer manager.Close()

	ctx := context.Background()
	return manager.docker.DeleteContainer(ctx, containerName, removeVolumes)
}

// deleteSQLiteDatabase removes a SQLite database file
func deleteSQLiteDatabase(name string) error {
	dbPath := name
	if filepath.Ext(name) != ".db" {
		dbPath = fmt.Sprintf("%s.db", name)
	}

	if err := os.Remove(dbPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("database file not found: %s", dbPath)
		}
		return fmt.Errorf("failed to delete database file: %w", err)
	}

	return nil
}

// deleteDuckDBDatabase removes a DuckDB database file
func deleteDuckDBDatabase(name string) error {
	dbPath := name
	if filepath.Ext(name) != ".duckdb" {
		dbPath = fmt.Sprintf("%s.duckdb", name)
	}

	if err := os.Remove(dbPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("database file not found: %s", dbPath)
		}
		return fmt.Errorf("failed to delete database file: %w", err)
	}

	return nil
}

// ListContainers returns all database containers
func ListContainers() ([]ContainerInfo, error) {
	manager, err := NewManager()
	if err != nil {
		return nil, err
	}
	defer manager.Close()

	ctx := context.Background()
	dockerContainers, err := manager.docker.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	// Add SQLite databases from current directory
	sqliteContainers, err := listSQLiteDatabases()
	if err != nil {
		return nil, err
	}

	// Add DuckDB databases from current directory
	duckdbContainers, err := listDuckDBDatabases()
	if err != nil {
		return nil, err
	}

	allContainers := append(dockerContainers, sqliteContainers...)
	allContainers = append(allContainers, duckdbContainers...)
	return allContainers, nil
}

// listSQLiteDatabases finds all .db files in current directory
func listSQLiteDatabases() ([]ContainerInfo, error) {
	var containers []ContainerInfo

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".db" {
			info, err := file.Info()
			if err != nil {
				continue
			}

			containers = append(containers, ContainerInfo{
				ID:       file.Name(),
				Name:     file.Name()[:len(file.Name())-3], // Remove .db extension
				Type:     "sqlite",
				Version:  "3",
				Status:   "file",
				Port:     0,
				User:     "",
				Password: "",
				Database: file.Name(),
				Volume:   "",
				Created:  info.ModTime(),
			})
		}
	}

	return containers, nil
}

// listDuckDBDatabases finds all .duckdb files in current directory
func listDuckDBDatabases() ([]ContainerInfo, error) {
	var containers []ContainerInfo

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".duckdb" {
			info, err := file.Info()
			if err != nil {
				continue
			}

			containers = append(containers, ContainerInfo{
				ID:       file.Name(),
				Name:     file.Name()[:len(file.Name())-7], // Remove .duckdb extension
				Type:     "duckdb",
				Version:  "latest",
				Status:   "file",
				Port:     0,
				User:     "",
				Password: "",
				Database: file.Name(),
				Volume:   "",
				Created:  info.ModTime(),
			})
		}
	}

	return containers, nil
}

// StopContainer stops a database container
func StopContainer(containerName string) error {
	// Handle SQLite and DuckDB separately (no container)
	if filepath.Ext(containerName) == ".db" || filepath.Ext(containerName) == ".duckdb" {
		return fmt.Errorf("cannot stop file-based database: %s", containerName)
	}

	manager, err := NewManager()
	if err != nil {
		return err
	}
	defer manager.Close()

	ctx := context.Background()
	return manager.docker.StopContainer(ctx, containerName)
}

// StartContainer starts a stopped database container
func StartContainer(containerName string) error {
	// Handle SQLite and DuckDB separately (no container)
	if filepath.Ext(containerName) == ".db" || filepath.Ext(containerName) == ".duckdb" {
		return fmt.Errorf("cannot start file-based database: %s", containerName)
	}

	manager, err := NewManager()
	if err != nil {
		return err
	}
	defer manager.Close()

	ctx := context.Background()
	return manager.docker.StartContainer(ctx, containerName)
}

// GetContainerInfo returns information about a container by name
func GetContainerInfo(containerName string) (*ContainerInfo, error) {
	manager, err := NewManager()
	if err != nil {
		return nil, err
	}
	defer manager.Close()

	ctx := context.Background()
	return manager.docker.GetContainerByName(ctx, containerName)
}
