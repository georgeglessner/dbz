package migrations

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dbz/dbz/pkg/containers"
)

// MigrationConfig holds configuration for running migrations
type MigrationConfig struct {
	File      string
	Direction string // "up" or "down"
	Database  string // target database name (container name or database name)
}

// RunMigration runs a database migration
func RunMigration(config MigrationConfig) error {
	if config.File == "" {
		return fmt.Errorf("migration file is required")
	}

	if config.Database == "" {
		return fmt.Errorf("database name is required (use --db flag)")
	}

	// Find the container
	container, err := findContainer(config.Database)
	if err != nil {
		return err
	}

	// Read migration file
	content, err := os.ReadFile(config.File)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute SQL using docker exec
	return executeSQLInContainer(container, string(content))
}

// findContainer finds a container by name or database name
// Prefers exact container name matches over database name matches
func findContainer(name string) (*containers.ContainerInfo, error) {
	containerList, err := containers.ListContainers()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// First, look for exact container name match
	for i := range containerList {
		if containerList[i].Name == name {
			return &containerList[i], nil
		}
	}

	// Then, look for database name match
	var matches []*containers.ContainerInfo
	for i := range containerList {
		if containerList[i].Database == name {
			matches = append(matches, &containerList[i])
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no running database found with name: %s", name)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple databases found with name '%s'. Use container name instead: %s", name, matches[0].Name)
	}

	return matches[0], nil
}

// executeSQLInContainer executes SQL in a running container
func executeSQLInContainer(container *containers.ContainerInfo, sql string) error {
	ctx := context.Background()

	// Get docker client
	dockerClient, err := containers.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer dockerClient.Close()

	// Determine database type and credentials
	dbType := container.Type
	user := container.User
	password := container.Password
	database := container.Database

	if user == "" {
		switch dbType {
		case "postgres", "postgresql":
			user = "postgres"
		case "mysql", "mariadb":
			user = "root"
		}
	}

	// Build the command based on database type
	var cmdArgs []string

	switch dbType {
	case "postgres", "postgresql":
		cmdArgs = []string{
			"exec",
			"-i",
			container.Name,
			"psql",
			"-U", user,
			"-d", database,
		}

	case "mysql", "mariadb":
		cmdArgs = []string{
			"exec",
			"-i",
			container.Name,
			"mysql",
			"-u", user,
			"-p" + password,
			database,
		}

	default:
		return fmt.Errorf("unsupported database type for migration: %s", dbType)
	}

	// Execute docker command with SQL as stdin
	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
	cmd.Stdin = strings.NewReader(sql)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variable for PostgreSQL password
	if dbType == "postgres" || dbType == "postgresql" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
