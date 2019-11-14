package data

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"path/filepath"
	"runtime"
)

// NewSchemaMigration creates a migration instance that can be used to
// apply and rollback schema changes to a database
func NewSchemaMigration(d *sql.DB, dbName string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(d, &postgres.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to create migration db driver: %w", err)
	}

	_, path, _, _ := runtime.Caller(0)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.Dir(path)+"/migrations",
		dbName, driver,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return m, nil
}
