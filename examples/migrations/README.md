# Database Migrations

This folder contains example migration files demonstrating how to evolve your database schema over time.

## Naming Convention

Migration files should follow this naming pattern:
```
XXX_description.sql
```

Where:
- `XXX` is a sequential number (001, 002, 003, etc.)
- `description` is a brief description of what the migration does
- `.sql` is the file extension

## Running Migrations

### Option 1: Run individually

```bash
# Run migrations one at a time in order
dbz migrate 001_create_users_table.sql --db myapp
dbz migrate 002_add_user_profile.sql --db myapp
dbz migrate 003_create_posts_table.sql --db myapp
```

### Option 2: Use with created database

```bash
# Create database
dbz create postgres --database myapp

# Run all migrations
dbz migrate 001_create_users_table.sql --db myapp
dbz migrate 002_add_user_profile.sql --db myapp
dbz migrate 003_create_posts_table.sql --db myapp
```

## Best Practices

1. **Make migrations idempotent**: Use `IF NOT EXISTS` for CREATE statements and `IF EXISTS` for DROP statements
2. **One change per migration**: Each migration should do one thing (create a table, add a column, etc.)
3. **Never modify existing migrations**: Once a migration has been run, create a new one instead
4. **Test migrations**: Always test on a copy of your database before running on production
5. **Keep migrations in version control**: Track all migration files in git

## Example Migration Template

```sql
-- Migration XXX: Brief description
-- Date: YYYY-MM-DD
-- Author: Your Name

-- Up migration (forward)
CREATE TABLE IF NOT EXISTS table_name (
    id SERIAL PRIMARY KEY,
    column_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add index if needed
CREATE INDEX IF NOT EXISTS idx_table_name_column ON table_name(column_name);

-- Add comments
COMMENT ON TABLE table_name IS 'Description of what this table stores';
```

## Rollback Strategy

While these examples don't include rollback scripts, in production you should create corresponding down migrations:

```sql
-- Rollback for 001_create_users_table.sql
DROP TABLE IF EXISTS users CASCADE;
```

Name them with a suffix like `_rollback.sql`:
- `001_create_users_table_rollback.sql`

Or use a migration tool that tracks applied migrations and can run rollbacks automatically.
