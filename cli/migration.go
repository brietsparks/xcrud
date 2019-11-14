package cli

import (
	"database/sql"
	"fmt"
	"github.com/brietsparks/xcrud/data"
	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli"
)

// NewMigrateCommand returns a migration command tree that can be used by a urfave/cli instance
func NewMigrateCommand(name string, chVars chan data.Vars) cli.Command {
	var vars data.Vars
	var mig *migrate.Migrate

	return cli.Command{
		Name:  name,
		Usage: "execute migration operations",
		Before: func(c *cli.Context) error {
			vars = <-chVars

			url := data.MakeUrl(vars)
			db, err := sql.Open("postgres", url)

			if err != nil {
				return err
			}

			m, err := data.NewSchemaMigration(db, vars.Name)

			if err != nil {
				return err
			}

			mig = m
			return nil
		},
		Subcommands: []cli.Command{
			{
				Name:  "up",
				Usage: "execute migrations",
				Action: func(c *cli.Context) error {
					fmt.Println("migrating up...")
					return mig.Up()
				},
			},
			{
				Name:  "down",
				Usage: "rollback migrations",
				Action: func(c *cli.Context) error {
					fmt.Println("migrating down...")
					return mig.Down()
				},
			},
		},
	}
}
