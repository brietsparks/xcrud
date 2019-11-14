package main

import (
	appcli "github.com/brietsparks/xcrud/cli"
	"github.com/brietsparks/xcrud/data"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	// log file
	file, err := os.OpenFile("error.log", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// logger
	l := logrus.New()
	l.Out = file
	l.Formatter = &logrus.JSONFormatter{}
	l.SetLevel(logrus.ErrorLevel)

	// cli app
	app := cli.NewApp()
	app.Writer = l.Writer()

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
	resourcesCommand := appcli.NewResourcesCommand("resources", chDataVars, l)

	app.Commands = []cli.Command{
		migrationCommand,
		resourcesCommand,
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
