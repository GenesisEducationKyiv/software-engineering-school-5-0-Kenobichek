package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type dbManagerImpl struct {
	db     *sql.DB
	logger loggerManager
}

func NewDBManager(db *sql.DB, logger loggerManager) *dbManagerImpl {
	return &dbManagerImpl{db: db, logger: logger}
}

func (d *dbManagerImpl) InitDB(dsn string) error {
	d.logger.Infof("Trying to establish database connection.")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		if err := db.Close(); err != nil {
			d.logger.Errorf("DB close error: %v", err)
		}
		return fmt.Errorf("failed to ping database: %w", err)
	}
	d.logger.Infof("Database connection established successfully.")
	d.db = db
	return nil
}

func (d *dbManagerImpl) RunMigrations(migrationsPath string) error {
	if d.db == nil {
		return fmt.Errorf("database connection (dbConn) is nil in RunMigrations")
	}
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
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
	d.logger.Infof("Migrations applied successfully")
	return nil
}

func (d *dbManagerImpl) GetDB() *sql.DB {
	return d.db
}