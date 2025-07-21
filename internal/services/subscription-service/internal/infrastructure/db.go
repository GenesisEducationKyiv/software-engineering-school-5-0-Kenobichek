package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func InitDB(dsn string) (*sql.DB, error) {
	log.Println("Trying to establish database connection.")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		if err := db.Close(); err != nil {
			log.Printf("DB close error: %v", err)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database connection established successfully.")
	return db, nil
}

func RunMigrations(dbConn *sql.DB, migrationsPath string) error {
	if dbConn == nil {
		return fmt.Errorf("database connection (dbConn) is nil in RunMigrations")
	}
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migration driver error: %w", err)
	}
	instance, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}
	err = instance.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}
