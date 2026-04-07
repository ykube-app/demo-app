package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := pool.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := migrate(pool); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return pool, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}
