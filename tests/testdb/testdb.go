package testdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage   = "postgres:15"
	dbName         = "testdb"
	dbUser         = "testuser"
	dbPassword     = "testpass"
	dbPort         = "5432/tcp"
	startupTimeout = 30 * time.Second
)

type PG struct {
	SQL       *sql.DB
	DSN       string
	container testcontainers.Container
	DBMutex   sync.Mutex
}

func New(t testing.TB) *PG {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: defaultImage,
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
		},
		ExposedPorts: []string{dbPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(dbPort)).WithStartupTimeout(startupTimeout),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start container: %v", err)
	}

	host, err := c.Host(ctx)
	if err != nil {
		failTest(t, c, "container host", err)
	}

	port, err := c.MappedPort(ctx, nat.Port(dbPort))
	if err != nil {
		failTest(t, c, "mapped port", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, port.Port(), dbName)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		failTest(t, c, "sql.Open", err)
	}

	if err := waitForDB(sqlDB, startupTimeout); err != nil {
		failTest(t, c, "postgres not ready", err)
	}

	if err := runMigrations(dsn); err != nil {
		failTest(t, c, "apply migrations", err)
	}

	pg := &PG{SQL: sqlDB, DSN: dsn, container: c}
	t.Cleanup(func() {
		_ = pg.SQL.Close()
		_ = pg.container.Terminate(context.Background())
	})
	return pg
}

func waitForDB(db *sql.DB, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("database did not respond within %v", timeout)
}

func runMigrations(dsn string) error {
	dir, err := filepath.Abs("../../../migrations")
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("driver init: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+dir, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func failTest(t testing.TB, c testcontainers.Container, msg string, err error) {
	_ = c.Terminate(context.Background())
	t.Fatalf("%s: %v", msg, err)
}
