package main

import (
	appcli "github.com/brietsparks/xcrud/cli"
	"github.com/brietsparks/xcrud/data"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	var envFilepath string
	chDataVars := make(chan data.Vars, 1)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "env, e",
			Usage:       "Load configuration from .env `FILE`",
			Destination: &envFilepath,
			Value: ".env",
		},
	}

	app.Before = func(context *cli.Context) error {
		vars, err := data.LoadEnvVars(envFilepath)

		if err != nil {
			return err
		}

		chDataVars <- vars

		return nil
	}

	migrationCommand := appcli.NewMigrateCommand("migrate", chDataVars)
	resourcesCommand := appcli.NewResourcesCommand("resources", chDataVars)

	app.Commands = []cli.Command{
		migrationCommand,
		resourcesCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
