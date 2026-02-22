# DBZ - Database CLI Tool

`dbz` is a powerful CLI tool that quickly and easily creates databases from the command line. It supports multiple database types and provides a simple interface for database management and migrations.

## Features

- 🚀 **Quick Setup** - Create databases in seconds with Docker containers
- 🗄️ **Multiple Database Support** - PostgreSQL, MySQL, MariaDB, SQLite, DuckDB, ClickHouse
- 📦 **Easy Installation** - Install with a single curl command
- 🔧 **Version Control** - Specify exact database versions (e.g., `mysql@8.4`)
- 🔄 **Migrations** - Run SQL migration files to create schemas and seed data
- 📋 **Container Management** - List, create, stop, start, reset, and delete database containers
- 🔒 **Secure** - Auto-generated passwords and secure defaults
- 🎯 **Auto Port Detection** - Automatically finds available ports when default ports are in use

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/georgeglessner/dbz/main/install.sh | bash
```

### From Source

```bash
git clone https://github.com/georgeglessner/dbz.git
cd dbz
go build -o dbz main.go
sudo mv dbz /usr/local/bin/
```

### Requirements

- Docker (for containerized databases)
- Go 1.21+ (for building from source)

## Usage

### Create Databases

```bash
# Create PostgreSQL with latest version
dbz create postgres

# Create MySQL with specific version
dbz create mysql@8.4

# Create MariaDB with custom port
dbz create mariadb --port 3307

# Create SQLite database (file-based, no container)
dbz create sqlite mydb

# Create database and run SQL initialization file
dbz create postgres init.sql
```

### List Databases

```bash
# List all database containers
dbz list

# Alias for list
dbz ls
```

### Delete Databases

```bash
# Delete container by name
dbz delete postgres

# Delete SQLite database file
dbz delete mydb.db
```

### Stop and Start Databases

```bash
# Stop a running container
dbz stop mysql-dev

# Start a stopped container
dbz start mysql-dev
```

### Reset Databases

```bash
# Reset (restart) a container - stops then starts it
dbz reset mysql-dev
```

### Update dbz

```bash
# Update dbz to the latest version from source
dbz update

# Update from a specific directory
dbz update --from /path/to/dbz
```

### Run Migrations

Execute SQL files to set up your database schema and populate it with data.

```bash
# Run a schema migration
dbz migrate myapp schema.sql

# Run seed data migration
dbz migrate myapp seed_data.sql

# Run a numbered migration file
dbz migrate myapp migrations/001_create_users.sql
```

**How it works:**
1. dbz connects to your running database using the provided database name
2. Reads the SQL file and executes each statement
3. Reports the results for each statement

**Typical workflow:**
1. Create your database: `dbz create postgres --name myapp --database myapp`
2. Create your schema SQL file with CREATE TABLE statements
3. Run the migration: `dbz migrate myapp schema.sql`
4. Add seed data: `dbz migrate myapp seed_data.sql`

**Migration files:**
Migration files are regular SQL files containing any valid SQL statements (CREATE TABLE, INSERT, ALTER TABLE, etc.). Organize them in a `migrations/` directory with numbered prefixes for ordering:
- `migrations/001_create_users.sql`
- `migrations/002_add_email_index.sql`
- `migrations/003_seed_test_data.sql`

## Supported Databases

### OLTP Databases
- **PostgreSQL** - `postgres`, `postgresql`
- **MySQL** - `mysql` 
- **MariaDB** - `mariadb`
- **SQLite** - `sqlite` (file-based, no container)

### OLAP Databases
- **DuckDB** - `duckdb`
- **ClickHouse** - `clickhouse`

## Command Reference

### Global Flags

- `--help, -h` - Show help
- `--version` - Show version

### Create Command

```bash
dbz create [database-type] [sql-file] [flags]
```

**Flags:**
- `--port, -p` - Port to expose (auto-assign if not specified)
- `--password` - Database password (auto-generate if not specified)
- `--database, -b` - Database name (default: testdb)
- `--name, -c` - Docker container name (auto-generated if not specified)
- `--user, -u` - Database user (default depends on database type)
- `--volume, -v` - Volume to mount for data persistence
- `--network, -n` - Docker network to join

**Examples:**
```bash
dbz create postgres
dbz create mysql@8.4 --port 3307 --password mypass
dbz create postgres init.sql --volume /my/data
dbz create mysql --database myapp --name myapp-db
dbz create postgres --database prod_db --user admin --name production-postgres
```

**Port Assignment:**
If no `--port` flag is specified and the default port is already in use, dbz will automatically find and use the next available port.

### Stop Command

Stop a running database container.

```bash
dbz stop [container-name]
```

**Examples:**
```bash
dbz stop postgres
dbz stop my-mysql-container
```

**Note:** The container can be started again with `dbz start`. For file-based databases (SQLite, DuckDB), this command is not applicable.

### Start Command

Start a previously stopped database container.

```bash
dbz start [container-name]
```

**Examples:**
```bash
dbz start postgres
dbz start my-mysql-container
```

**Note:** For file-based databases (SQLite, DuckDB), this command is not applicable.

### Reset Command

Reset (restart) a database container. This is equivalent to running `dbz stop` followed by `dbz start`.

```bash
dbz reset [container-name]
```

**Examples:**
```bash
dbz reset postgres
dbz reset my-mysql-container
```

**Note:** This does not delete any data. The container is stopped and started again. For a complete reset with data deletion, use `dbz delete` followed by `dbz create`.

### Migration Command

Execute SQL migration files against a running database container.

```bash
dbz migrate [database] [migration-file] [flags]
```

**Arguments:**
- `database` - Database name (matches container name or database name used during creation)
- `migration-file` - Path to SQL file containing migration statements

**Flags:**
- `--direction, -d` - Migration direction: `up` or `down` (default: `up`)

**Examples:**
```bash
# Run a schema migration
dbz migrate myapp schema.sql

# Run a numbered migration
dbz migrate myapp migrations/001_create_users.sql

# Run with specific direction (for rollback migrations)
dbz migrate myapp rollback.sql --direction down

# Run example migrations from the examples directory
dbz migrate blogdb examples/blog_schema.sql
dbz migrate blogdb examples/seed_blog_data.sql
```

**Creating migration files:**
Migration files are standard SQL files. You can create them manually or generate them:

```sql
-- schema.sql - Create tables
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- seed_data.sql - Insert data
INSERT INTO users (name, email) VALUES 
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com');
```

**Tip:** When dbz creates a database, it shows the container name and connection details. Use the container name (from `--name`) or database name (from `--database`) as the first argument to `migrate`.

### Update Command

```bash
dbz update [flags]
```

**Flags:**
- `--from, -f` - Path to dbz source directory (default: current directory or auto-detected)

**Examples:**
```bash
dbz update
dbz update --from ~/Projects/dbz
```

## Environment Variables

- `DOCKER_HOST` - Docker daemon host (default: unix:///var/run/docker.sock)

## Examples

### Development Workflow

Complete example of creating a database and running migrations:

```bash
# 1. Create a PostgreSQL database (container will be named myapp-db)
dbz create postgres --database myapp --name myapp-db

# 2. Create your schema file (schema.sql)
cat > schema.sql << 'EOF'
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
EOF

# 3. Run the schema migration
dbz migrate myapp schema.sql

# 4. Add seed data (via SQL migration file)
cat > seed_data.sql << 'EOF'
INSERT INTO users (name, email) VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com');
EOF
dbz migrate myapp seed_data.sql

# 5. Verify the database
dbz list

# 6. Clean up when done
dbz delete myapp-db
```

### CI/CD Integration

Using dbz in automated testing pipelines:

```bash
# In your CI pipeline
# 1. Create test database
dbz create postgres --database testdb --name testdb --password testpass

# 2. Apply schema migrations
dbz migrate testdb migrations/001_create_tables.sql
dbz migrate testdb migrations/002_add_indexes.sql

# 3. Insert test data via SQL migration
dbz migrate testdb test_data.sql

# 4. Run your application tests...
# npm test, pytest, go test, etc.

# 5. Clean up test database
dbz delete testdb
```

### Multiple Databases

```bash
# Create multiple databases for microservices
dbz create postgres --name user-service --port 5432
dbz create postgres --name order-service --port 5433
dbz create mysql --name inventory-service --port 3306

# List all running databases
dbz list
```

## Examples

The `examples/` directory contains sample SQL files to help you get started:

### Database Schemas
- **`blog_schema.sql`** - Complete blog database with users, posts, comments
- **`ecommerce_schema.sql`** - E-commerce database with products, orders, customers

### Migrations
- **`migrations/`** - Numbered migration files showing schema evolution
  - `001_create_users_table.sql`
  - `002_add_user_profile.sql`
  - `003_create_posts_table.sql`

### Quick Start with Examples

Try out the included example files:

```bash
# Create a blog database with sample posts and comments
dbz create postgres --database blogdb
dbz migrate blogdb examples/blog_schema.sql      # Creates tables
dbz migrate blogdb examples/seed_blog_data.sql     # Inserts sample data

# Create an e-commerce database with products and orders
dbz create postgres --database shopdb
dbz migrate shopdb examples/ecommerce_schema.sql     # Creates tables
dbz migrate shopdb examples/seed_ecommerce_data.sql  # Inserts sample data

# Explore the examples directory for more SQL files:
# - examples/migrations/ - Numbered migration files
# - examples/*.sql - Complete schema and seed files
```

## Troubleshooting

### Docker Issues

If you encounter Docker-related errors:

1. **Docker not running**: Start Docker Desktop or Docker daemon
2. **Permission denied**: Add your user to the docker group or use sudo
3. **Port conflicts**: Use `--port` flag to specify a different port

### Connection Issues

1. **Container not ready**: Wait a few seconds for the database to start
2. **Wrong credentials**: Check the output for auto-generated passwords
3. **Network issues**: Ensure Docker networking is configured correctly

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add some new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## Future Improvements

Planned features for future releases:

- **Data Seeding** - Built-in command to generate and insert fake test data (e.g., `dbz seed mydb users --rows 100`). Will support common table types with realistic fake data.

- **Custom Configuration** - Support for `.dbz.yaml` configuration files to set default ports, versions, and other settings per project. Will allow defining database defaults, migration directories, and Docker settings in a config file.

## Development

```bash
# Clone the repository
git clone https://github.com/georgeglessner/dbz.git
cd dbz

# Install dependencies
go mod download

# Build the binary
go build -o dbz main.go

# Run tests
go test ./...

# Install locally
go install
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Docker integration using the official Docker Go SDK

---

**Happy database management! 🚀**