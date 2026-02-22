-- Test Setup SQL
-- Run this to verify your database connection and basic functionality

-- Create a test table
CREATE TABLE IF NOT EXISTS dbz_test (
    id SERIAL PRIMARY KEY,
    test_name VARCHAR(100) NOT NULL,
    test_value VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert test data
INSERT INTO dbz_test (test_name, test_value) VALUES
('connection_test', 'Database connection successful!'),
('timestamp_test', NOW()::TEXT),
('version_test', version());

-- Query to verify everything works
SELECT 
    'SUCCESS' as status,
    test_name,
    test_value,
    created_at
FROM dbz_test
ORDER BY created_at DESC;

-- Show database information
SELECT 
    current_database() as database_name,
    current_user as current_user,
    version() as postgres_version;

-- List all tables in current database
SELECT 
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;

-- Cleanup (optional - uncomment to remove test data)
-- DROP TABLE IF EXISTS dbz_test;
