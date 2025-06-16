package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

var DataBase *sql.DB

// Init establishes a PostgreSQL database connection using the provided DSN.
// It returns the database handle on success or an error if the connection or ping fails.
func Init(dsn string) (*sql.DB, error) {
	log.Println("Trying to establish database connection.")

	var err error
	DataBase, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := DataBase.Ping(); err != nil {
		if closeErr := DataBase.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to close database: %w and failed to ping database: %v", closeErr, err)
		}
		DataBase = nil
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully.")
	return DataBase, nil
}
// RunMigrations applies all up migrations to the connected PostgreSQL database.
// Returns an error if the database connection is nil, migration driver initialization fails, or migration execution encounters an error other than no changes detected.
func RunMigrations(dbConn *sql.DB) error {
	if dbConn == nil {
		return fmt.Errorf("database connection (dbConn) is nil in RunMigrations")
	}

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}

	instance, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
