package containers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// DockerClient wraps the Docker client for easier testing
type DockerClient struct {
	client *client.Client
}

// NewDockerClient creates a new Docker client
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Test connection
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	return &DockerClient{client: cli}, nil
}

// Close closes the Docker client
func (d *DockerClient) Close() error {
	return d.client.Close()
}

// CreateContainer creates a new container with the given configuration
func (d *DockerClient) CreateContainer(ctx context.Context, config ContainerConfig, db Database) (*ContainerInfo, error) {
	image := db.GetImage(config.Version)

	// Pull image if not exists
	if err := d.pullImage(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to pull image: %w", err)
	}

	// Generate container name (use provided name if available)
	containerName := config.ContainerName
	if containerName == "" {
		containerName = generateContainerName(config.Type)
	}

	// Ensure we have a password (generate if not provided)
	password := config.Password
	if password == "" {
		password = generatePassword()
		config.Password = password // Update config so GetEnvironment receives the password
	}

	// Determine database user and name (with defaults)
	dbUser := config.User
	dbName := config.Database
	if dbUser == "" {
		switch config.Type {
		case "postgres", "postgresql":
			dbUser = "postgres"
			if dbName == "" {
				dbName = "postgres"
			}
		case "mysql", "mariadb":
			dbUser = "root"
			if dbName == "" {
				dbName = "testdb"
			}
		}
	}
	// Apply defaults if still empty
	if dbName == "" {
		dbName = "testdb"
	}

	// Set up port binding
	port := config.Port
	if port == 0 {
		// Find an available port starting from the default
		defaultPort := db.GetDefaultPort()
		port = findAvailablePort(defaultPort)
		if port == 0 {
			return nil, fmt.Errorf("could not find an available port")
		}
		if port != defaultPort {
			fmt.Printf("Port %d is taken, using port %d instead\n", defaultPort, port)
		}
	}

	portBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: fmt.Sprintf("%d", port),
	}

	containerPort := nat.Port(fmt.Sprintf("%d/tcp", db.GetDefaultPort()))
	portMap := nat.PortMap{
		containerPort: []nat.PortBinding{portBinding},
	}

	// Environment variables
	envVars := []string{}
	for key, value := range db.GetEnvironment(config) {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	// Container configuration
	containerConfig := &container.Config{
		Image:        image,
		Env:          envVars,
		ExposedPorts: nat.PortSet{containerPort: struct{}{}},
		Labels: map[string]string{
			"dbz":          "true",
			"dbz.type":     config.Type,
			"dbz.version":  config.Version,
			"dbz.database": dbName,
			"dbz.user":     dbUser,
			"dbz.password": password,
		},
		Cmd: getMySQLCommand(config.Type), // Add command to set authentication plugin
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		PortBindings: portMap,
		AutoRemove:   false,
	}

	// Add volume if specified
	if config.Volume != "" {
		hostConfig.Binds = []string{fmt.Sprintf("%s:/data", config.Volume)}
	}

	// Network configuration
	networkConfig := &network.NetworkingConfig{}
	if config.Network != "" {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			config.Network: {},
		}
	}

	// Create container
	resp, err := d.client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get container info
	info, err := d.client.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	// Get connection info
	connInfo := db.GetConnectionInfo(config, containerName)

	// Execute SQL file if provided
	if config.SQLFile != "" {
		if err := d.executeSQLFile(ctx, resp.ID, config.SQLFile, db); err != nil {
			// Don't fail if SQL execution fails, just log the error
			fmt.Printf("Warning: failed to execute SQL file: %v\n", err)
		}
	}

	// For MySQL/MariaDB, ensure we have a user with native password authentication
	if config.Type == "mysql" || config.Type == "mariadb" {
		// Wait a bit for the database to be ready, then set up native password auth
		go func() {
			time.Sleep(10 * time.Second) // Give MySQL time to start
			if err := d.ensureNativePasswordAuth(ctx, resp.ID, password, db); err != nil {
				fmt.Printf("Warning: failed to set native password authentication: %v\n", err)
			}
		}()
	}

	return &ContainerInfo{
		ID:       resp.ID,
		Name:     containerName,
		Type:     config.Type,
		Version:  config.Version,
		Status:   info.State.Status,
		Port:     port,
		User:     dbUser,
		Password: password,
		Database: dbName,
		DSN:      connInfo.DSN,
		Volume:   config.Volume,
		Created:  time.Now(),
	}, nil
}

// DeleteContainer removes a container by name or ID
// findAvailablePort checks if a port is available and finds the next available one
func findAvailablePort(startPort int) int {
	port := startPort
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err == nil {
			listener.Close()
			return port
		}
		port++
		// Prevent infinite loop, cap at 65535
		if port > 65535 {
			return 0
		}
	}
}

// StopContainer stops a container by name or ID
func (d *DockerClient) StopContainer(ctx context.Context, containerName string) error {
	// Try to find container by name first
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	var containerID string
	for _, c := range containers {
		if c.Names[0] == "/"+containerName || c.ID[:12] == containerName {
			containerID = c.ID
			break
		}
	}

	if containerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	// Stop container
	if err := d.client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	return nil
}

// StartContainer starts a stopped container by name or ID
func (d *DockerClient) StartContainer(ctx context.Context, containerName string) error {
	// Try to find container by name first
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	var containerID string
	for _, c := range containers {
		if c.Names[0] == "/"+containerName || c.ID[:12] == containerName {
			containerID = c.ID
			break
		}
	}

	if containerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	// Start container
	if err := d.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

func (d *DockerClient) DeleteContainer(ctx context.Context, containerName string, removeVolumes bool) error {
	// Try to find container by name first
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	var containerID string
	for _, c := range containers {
		if c.Names[0] == "/"+containerName || c.ID[:12] == containerName {
			containerID = c.ID
			break
		}
	}

	if containerID == "" {
		return fmt.Errorf("container not found: %s", containerName)
	}

	// Stop container if running
	if err := d.client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		// Ignore if container is already stopped
		if !strings.Contains(err.Error(), "is not running") {
			return fmt.Errorf("failed to stop container: %w", err)
		}
	}

	// Remove container with optional volume removal
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
	}
	if err := d.client.ContainerRemove(ctx, containerID, removeOptions); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// ListContainers returns all database containers created by dbz
func (d *DockerClient) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var dbContainers []ContainerInfo
	for _, c := range containers {
		// Check if this is a dbz container
		if c.Labels["dbz"] != "true" {
			continue
		}

		// Extract port information
		port := 0
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				port = int(p.PublicPort)
				break
			}
		}

		info := ContainerInfo{
			ID:       c.ID[:12],
			Name:     c.Names[0][1:], // Remove leading slash
			Type:     c.Labels["dbz.type"],
			Version:  c.Labels["dbz.version"],
			Status:   c.State,
			Port:     port,
			Database: c.Labels["dbz.database"],
			User:     c.Labels["dbz.user"],
			Password: c.Labels["dbz.password"],
			Created:  time.Unix(c.Created, 0),
		}

		dbContainers = append(dbContainers, info)
	}

	return dbContainers, nil
}

// pullImage pulls a Docker image if it doesn't exist locally
func (d *DockerClient) pullImage(ctx context.Context, image string) error {
	// Check if image exists
	_, _, err := d.client.ImageInspectWithRaw(ctx, image)
	if err == nil {
		return nil // Image already exists
	}

	// Pull image
	fmt.Printf("Pulling image %s...\n", image)
	reader, err := d.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Read and discard output (for progress indication)
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("failed to read pull output: %w", err)
	}

	return nil
}

// executeSQLFile executes an SQL file in a container
func (d *DockerClient) executeSQLFile(ctx context.Context, containerID string, sqlFile string, db Database) error {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Copy the SQL file to the container
	// 2. Execute it using the appropriate client for the database type
	// 3. Handle errors and output appropriately

	// For now, we'll just log that SQL execution is not yet implemented
	fmt.Printf("Note: SQL file execution not yet implemented for %s\n", db.GetImage(""))
	return nil
}

// GetContainerByName finds a container by name and returns its info
func (d *DockerClient) GetContainerByName(ctx context.Context, name string) (*ContainerInfo, error) {
	containers, err := d.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if c.Name == name || c.ID == name {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("container not found: %s", name)
}

// ensureNativePasswordAuth ensures MySQL/MariaDB uses native password authentication for compatibility
func (d *DockerClient) ensureNativePasswordAuth(ctx context.Context, containerID string, rootPassword string, db Database) error {
	// Wait a moment for the database to be ready
	fmt.Println("Configuring MySQL/MariaDB for native password authentication...")

	// SQL to set root user to use native password authentication
	sql := fmt.Sprintf(`
		-- Set root user to use native password authentication
		ALTER USER 'root'@'%%' IDENTIFIED WITH mysql_native_password BY '%s';
		-- Create a test user with native password authentication
		CREATE USER IF NOT EXISTS 'dbzuser'@'%%' IDENTIFIED WITH mysql_native_password BY '%s';
		-- Grant privileges
		GRANT ALL PRIVILEGES ON *.* TO 'dbzuser'@'%%';
		FLUSH PRIVILEGES;
	`, rootPassword, rootPassword)

	// Execute the SQL through the container
	resp, err := d.client.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          []string{"mysql", "-u", "root", "-p" + rootPassword, "-e", sql},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create exec: %w", err)
	}

	// Start the exec
	attachResp, err := d.client.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		return fmt.Errorf("failed to attach exec: %w", err)
	}
	defer attachResp.Close()

	// Read output (optional - for debugging)
	output, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return fmt.Errorf("failed to read exec output: %w", err)
	}

	// Check exec status
	execInspect, err := d.client.ContainerExecInspect(ctx, resp.ID)
	if err != nil {
		return fmt.Errorf("failed to inspect exec: %w", err)
	}

	if execInspect.ExitCode != 0 {
		return fmt.Errorf("SQL execution failed with exit code %d. Output: %s", execInspect.ExitCode, string(output))
	}

	fmt.Println("✅ Successfully configured native password authentication")
	return nil
}

// generateContainerName generates a unique container name
func generateContainerName(dbType string) string {
	return fmt.Sprintf("dbz-%s-%d", dbType, time.Now().Unix())
}

// generatePassword generates a cryptographically secure random password
func generatePassword() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based if crypto/rand fails
		return "dbz_" + fmt.Sprintf("%d", time.Now().Unix())[:8]
	}
	return "dbz_" + base64.URLEncoding.EncodeToString(b)[:12]
}

// getMySQLCommand returns the command to run MySQL/MariaDB with proper authentication settings
func getMySQLCommand(dbType string) []string {
	// For now, use default commands. We'll handle authentication plugin via user creation
	// The authentication plugin needs to be set during user creation, not at startup
	switch dbType {
	case "mysql", "mariadb":
		return nil // Use default entrypoint script
	default:
		return nil // Use default command for other databases
	}
}
