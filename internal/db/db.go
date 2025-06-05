package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DataBase *sql.DB

func Init() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	log.Println("Trying to establish connection.")

	var err error

	DataBase, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DataBase:", err)
	}

	if err := DataBase.Ping(); err != nil {
		log.Fatal("Failed to ping DataBase:", err)
	}

	log.Println("Database connection established.")
}

func RunMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Migration driver error: %v", err)
	}

	instance, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Migration init error: %v", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")
}
