package main

import (
	"github.com/018bf/example/internal/containers"
	"github.com/urfave/cli/v2"
	"os"
)

const version = "0.1.0"

var configPath = ""

func main() {
	app := &cli.App{
		Name:    "example",
		Usage:   "service",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				EnvVars:     []string{"EXAMPLE_CONFIG_PATH"},
				TakesFile:   true,
				Value:       configPath,
				Destination: &configPath,
				HasBeenSet:  false,
			},
		},
		Action: runApp,
		Commands: []*cli.Command{
			{
				Name:      "migrate",
				Usage:     "Run migrations",
				Action:    runMigrations,
				ArgsUsage: "",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// runApp - run app
func runApp(context *cli.Context) error {
	app := containers.NewGRPCExample(configPath)
	err := app.Start(context.Context)
	if err != nil {
		return err
	}
	return nil
}

// runMigrations - migrate database
func runMigrations(context *cli.Context) error {
	app := containers.NewMigrate(configPath)
	err := app.Start(context.Context)
	if err != nil {
		return err
	}
	return nil
}
