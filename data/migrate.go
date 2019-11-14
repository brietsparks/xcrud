package data

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"path/filepath"
	"runtime"
)

func NewSchemaMigration(vars Vars) (*migrate.Migrate, error) {
	url := MakeUrl(vars)

	d, err := sql.Open("postgres", url)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	driver, err := postgres.WithInstance(d, &postgres.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to migration db driver: %w", err)
	}

	_, path, _, _ := runtime.Caller(0)

	m, err := migrate.NewWithDatabaseInstance(
		"file://" + filepath.Dir(path) + "/migrations",
		vars.Name, driver,
	)

	return m, err
}
