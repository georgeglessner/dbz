# DBZ Setup Guide

This guide will help you set up and run the `dbz` CLI tool.

## Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - For cloning the repository

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/georgeglessner/dbz.git
cd dbz
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build the Binary

```bash
make build
```

Or manually:
```bash
go build -o dbz main.go
```

### 4. Install System-wide (Optional)

```bash
make install
```

Or manually:
```bash
sudo cp dbz /usr/local/bin/
sudo chmod +x /usr/local/bin/dbz
```

### 5. Test Installation

```bash
dbz --version
dbz --help
```

## Development Setup

### Using Make (Recommended)

```bash
# Build and run
make dev

# Run tests
make test

# Build for all platforms
make release

# Clean up
make clean
```

### Manual Development

```bash
# Build
go build -o dbz main.go

# Run tests
go test -v ./...

# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

## Usage Examples

### Basic Usage

```bash
# Create a PostgreSQL database
dbz create postgres

# Create MySQL with specific version
dbz create mysql@8.4

# List all databases
dbz list

# Delete a database
dbz delete postgres
```

### Advanced Usage

```bash
# Create with custom settings
dbz create postgres --port 5433 --password mypass --volume /my/data

# Create and run SQL file
dbz create postgres init.sql

# Seed data
dbz seed postgres users --rows 500

# Run migrations
dbz migrate schema.sql
```

## Troubleshooting

### Docker Issues

If you get Docker connection errors:

1. **Check Docker is running:**
   ```bash
   docker ps
   ```

2. **Check Docker permissions:**
   ```bash
   sudo usermod -aG docker $USER
   # Then log out and back in
   ```

3. **Check Docker host:**
   ```bash
   export DOCKER_HOST=unix:///var/run/docker.sock
   ```

### Build Issues

1. **Go version too old:**
   ```bash
   go version  # Should be 1.21+
   ```

2. **Missing dependencies:**
   ```bash
   go mod tidy
   go mod download
   ```

3. **Permission denied on install:**
   ```bash
   sudo make install
   ```

### Runtime Issues

1. **Port already in use:**
   ```bash
   # Use a different port
   dbz create postgres --port 5433
   ```

2. **Container conflicts:**
   ```bash
   # List and remove conflicting containers
   docker ps -a
   docker rm conflicting-container
   ```

## Testing

### Unit Tests

```bash
make test
```

### Integration Tests

```bash
# Requires Docker running
make build
./dbz create postgres test
dbz list
./dbz delete postgres
```

### Test Coverage

```bash
make coverage
# Opens coverage.html in browser
```

## Development Workflow

1. **Make changes** to the source code
2. **Run tests:** `make test`
3. **Build:** `make build`
4. **Test manually:** `./dbz create postgres`
5. **Clean up:** `make cleanup`


## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests: `make test`
6. Build and test: `make build && ./dbz --help`
7. Submit a pull request

## Next Steps

After setup, you can:

1. **Create your first database:**
   ```bash
   dbz create postgres
   ```

2. **Explore the CLI:**
   ```bash
   dbz --help
   dbz create --help
   ```

3. **Read the full documentation:**
   ```bash
   cat README.md
   ```

4. **Check out examples:**
   ```bash
   # Create with SQL file
   echo "CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100));" > init.sql
   dbz create postgres init.sql
   ```

## Getting Help

- **Help command:** `dbz --help`
- **Command help:** `dbz create --help`
- **GitHub Issues:** [Report bugs or request features](https://github.com/dbz/dbz/issues)
- **Documentation:** See README.md for full documentation

---

**Happy database management! 🚀**