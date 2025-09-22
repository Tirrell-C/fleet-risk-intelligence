-- Initial database setup for Fleet Risk Intelligence
-- This script runs when MySQL container starts for the first time

CREATE DATABASE IF NOT EXISTS fleet_dev;
USE fleet_dev;

-- The tables will be created by GORM AutoMigrate when the services start
-- This script just ensures the database exists

-- Add some basic indexes for performance
-- These will be recreated by GORM if they don't exist
SELECT 'Database initialization complete' as status;