# DBZ Example SQL Files

This directory contains sample SQL files that demonstrate database creation, migrations, and seeding with dbz.

## Understanding Migrate vs Seed

**`dbz migrate`** - Executes SQL files to modify database schema or run arbitrary SQL:
- Creating tables, indexes, constraints
- Running migration files
- Executing seed data SQL files (if they contain INSERT statements)
- Any DDL or DML operations
- **Requires `--db` flag** to specify which database to run against

**`dbz seed`** - Generates fake data programmatically and inserts it:
- Uses faker to generate realistic test data
- Inserts directly into specified tables
- Good for quickly populating tables with test data
- No SQL file needed
- **Requires database name** as first argument

## Important: Schema First!

**You must create the database schema BEFORE running seed data!**

The seed files (like `seed_blog_data.sql`) contain INSERT statements that assume the tables already exist. If you try to run seed data before creating the schema, you'll get errors like:
```
ERROR: relation "users" does not exist
```

**Correct workflow:**
1. Create the database container
2. Run schema/migration files (creates tables)
3. Then run seed data files (inserts data)

## Files

### Database Schemas

- **`blog_schema.sql`** - A simple blog database with users, posts, and comments
- **`ecommerce_schema.sql`** - E-commerce database with products, orders, and customers

### Migrations

- **`001_create_users_table.sql`** - Initial migration creating users table
- **`002_add_user_profile.sql`** - Migration adding profile fields to users
- **`003_create_posts_table.sql`** - Migration creating posts table with foreign key

### Seed Data (SQL Files)

- **`seed_blog_data.sql`** - Sample INSERT statements for blog database
- **`seed_ecommerce_data.sql`** - Sample INSERT statements for e-commerce database

### Utility

- **`test_setup.sql`** - Quick test to verify database connection

## Usage Examples

### Create a database with schema

```bash
# Create PostgreSQL database
dbz create postgres --database blogdb

# Run schema creation (uses migrate command)
dbz migrate examples/blog_schema.sql --db blogdb
```

### Run migrations

```bash
# Run migration files in order
dbz migrate examples/migrations/001_create_users_table.sql --db myapp
dbz migrate examples/migrations/002_add_user_profile.sql --db myapp
dbz migrate examples/migrations/003_create_posts_table.sql --db myapp
```

### Seed with SQL file data

```bash
# Load sample data from SQL file (uses migrate command for SQL execution)
dbz migrate examples/seed_blog_data.sql --db blogdb
```

### Seed with generated fake data

```bash
# Generate and insert 100 fake users (uses seed command with faker)
dbz seed blogdb users --rows 100

# Generate and insert 50 fake products
dbz seed blogdb products --rows 50
```

## Quick Start Workflow

```bash
# 1. Create database container
dbz create postgres --database myapp --name myapp-db

# 2. Create schema (REQUIRED before seeding!)
dbz migrate examples/blog_schema.sql --db myapp
# or run numbered migrations:
dbz migrate examples/migrations/001_create_users_table.sql --db myapp
dbz migrate examples/migrations/002_add_user_profile.sql --db myapp

# 3. Load sample data from SQL file
dbz migrate examples/seed_blog_data.sql --db myapp

# OR generate fake data programmatically:
dbz seed myapp users --rows 100
dbz seed myapp posts --rows 50

# 4. Check what's running
dbz list

# 5. Connect to database
psql postgresql://postgres:password@localhost:5432/myapp
```

## Tips

- Migration files are numbered to ensure they run in the correct order
- Each migration should be idempotent (can be run multiple times safely)
- Use `IF NOT EXISTS` for CREATE statements to make migrations safer
- Use `ON CONFLICT` for INSERT statements in seed files to allow re-running
- The `seed` command generates fake data programmatically, while `migrate` executes SQL files
