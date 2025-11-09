package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"vicnotes/backend/config"
	"vicnotes/backend/utils"
)

// InitDB initializes the database connection with retry logic
func InitDB() (*sql.DB, error) {
	dbURL := config.GetDatabaseURL()
	var db *sql.DB

	retryConfig := utils.RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}

	err := utils.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}

		// Test the connection
		if err := db.Ping(); err != nil {
			return fmt.Errorf("failed to ping database: %w", err)
		}

		return nil
	}, retryConfig)

	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// RunMigrations creates necessary tables if they don't exist
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
