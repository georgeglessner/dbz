-- Migration 002: Add user profile fields
-- Adds profile information to the users table

-- Add profile columns to users table
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS full_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS bio TEXT,
    ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP;

-- Create index for searching by name
CREATE INDEX IF NOT EXISTS idx_users_full_name ON users(full_name);

COMMENT ON COLUMN users.full_name IS 'User display name';
COMMENT ON COLUMN users.bio IS 'User biography/description';
COMMENT ON COLUMN users.avatar_url IS 'URL to user avatar image';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last login';
