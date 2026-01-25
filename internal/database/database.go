package database

import (
	"database/sql"
	"fmt"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Database wraps the SQL database connection
type Database struct {
	DB *sql.DB
}

// New creates a new database connection
func New(cfg Config) (*Database, error) {
	// TODO: Implement actual database connection
	// This is a skeleton for future implementation
	return &Database{}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Ping checks database connectivity
func (d *Database) Ping() error {
	if d.DB != nil {
		return d.DB.Ping()
	}
	return fmt.Errorf("database not initialized")
}
