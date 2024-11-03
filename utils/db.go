package utils

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/mattn/go-sqlite3"
)

func GetLocalDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./customers.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
        CREATE TABLE IF NOT EXISTS customers (
            id UUID PRIMARY KEY,
            first_name TEXT NOT NULL,
            middle_name TEXT,
            last_name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            phone_number TEXT
        );
        `

	if _, err = db.Exec(createTableSQL); err != nil {
		return nil, err
	}
	return db, nil
}

func RunMigrations(dbURL string) error {
	// Use the MIGRATIONS_DIR environment variable, or fallback to "./migrations"
	dir := os.Getenv("MIGRATIONS_DIR")
	if dir == "" {
		dir = "./migrations"
	}

	// Get the absolute path to the migrations directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	// Open the database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Instantiate the PostgreSQL driver for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Create a new migrate instance with the database driver
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absDir,
		"postgres", driver)
	if err != nil {
		return err
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
