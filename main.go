package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"os/exec"
	"path"
)

const version = "0.1.3"

var (
	serviceName     string
	moduleName      string
	goVersion       string
	destinationPath = "."
	models          cli.StringSlice
)

//go:embed templates/*
var content embed.FS

func main() {
	app := &cli.App{
		Name:    "Creathor",
		Usage:   "generate stub for service",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "destination",
				Aliases:     []string{"d"},
				Usage:       "module name",
				Destination: &destinationPath,
				Required:    false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "create base files and directories",
				Action: initProject,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "name",
						Aliases:     []string{"n"},
						Usage:       "service name",
						Destination: &serviceName,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "module",
						Usage:       "module name",
						Destination: &moduleName,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "go",
						Usage:       "go",
						Destination: &goVersion,
						Required:    false,
						Value:       "1.19",
					},
					&cli.StringSliceFlag{
						Name:        "model",
						Aliases:     []string{"m"},
						Usage:       "generate models",
						Required:    false,
						Destination: &models,
					},
				},
			},
			{
				Name:   "model",
				Usage:  "create base model and CRUD interfaces",
				Action: initModels,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:        "model",
						Aliases:     []string{"m"},
						Usage:       "model name",
						Destination: &models,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "module",
						Usage:       "module name",
						Destination: &moduleName,
						Required:    false,
						Hidden:      true,
						Value:       getModuleName(),
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initProject(ctx *cli.Context) error {
	data := &Project{Name: serviceName, Module: moduleName, GoVersion: goVersion}
	if err := CreateLayout(data); err != nil {
		return err
	}
	if err := CreateCI(); err != nil {
		return err
	}
	if err := CreateDI(data); err != nil {
		return err
	}
	if err := CreateBuild(data); err != nil {
		return err
	}
	for _, model := range models.Value() {
		if err := CreateCRUD(model, moduleName); err != nil {
			return err
		}
	}
	postInit()
	return nil
}

func initModels(ctx *cli.Context) error {
	for _, model := range models.Value() {
		if err := CreateCRUD(model, moduleName); err != nil {
			return err
		}
	}
	postInit()
	return nil
}

func getModuleName() string {
	gomod, err := os.ReadFile(path.Join(destinationPath, "go.mod"))
	if err != nil {
		return serviceName
	}
	return modfile.ModulePath(gomod)
}

func postInit() {
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = destinationPath
	if err := tidy.Run(); err != nil {
		fmt.Println(err.Error())
	}

	generate := exec.Command("go", "generate", "./...")
	generate.Dir = destinationPath
	if err := generate.Run(); err != nil {
		fmt.Println(err.Error())
	}

	clean := exec.Command("golangci-lint", "run", "./...", "--fix")
	clean.Dir = destinationPath
	_ = clean.Run()
}
