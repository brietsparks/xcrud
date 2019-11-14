package cli

import (
	"fmt"
	"github.com/brietsparks/xcrud/data"
	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli"
)

func NewMigrateCommand(name string, chVars chan data.Vars) cli.Command {
	var vars data.Vars
	var mig *migrate.Migrate

	return cli.Command{
		Name:  name,
		Usage: "execute migration operations",
		Before: func(c *cli.Context) error {
			vars = <-chVars

			m, err := data.NewSchemaMigration(vars)

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
