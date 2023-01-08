package main

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const version = "0.2.1"

var (
	serviceName     string
	moduleName      string
	goVersion       string
	destinationPath = "."
	models          cli.StringSlice
	authEnabled     bool
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
			&cli.BoolFlag{
				Name:        "auth",
				Aliases:     []string{"a"},
				Usage:       "enable auth",
				Destination: &authEnabled,
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
	data := &Project{Name: serviceName, Module: moduleName, GoVersion: goVersion, Auth: authEnabled}
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
	if err := CreateDeployment(data); err != nil {
		return err
	}
	for _, model := range models.Value() {
		if err := CreateCRUD(Model{Model: model, Module: moduleName, Auth: authEnabled}); err != nil {
			return err
		}
	}
	if err := postInit(); err != nil {
		return err
	}
	return nil
}

func initModels(ctx *cli.Context) error {
	for _, model := range models.Value() {
		if err := CreateCRUD(Model{Model: model, Module: moduleName, Auth: authEnabled}); err != nil {
			return err
		}
	}
	if err := postInit(); err != nil {
		return err
	}
	return nil
}

func getModuleName() string {
	gomod, err := os.ReadFile(path.Join(destinationPath, "go.mod"))
	if err != nil {
		return serviceName
	}
	return modfile.ModulePath(gomod)
}

func postInit() error {
	fmt.Println("post init...")
	var errb bytes.Buffer
	generate := exec.Command("go", "generate", "./...")
	generate.Dir = destinationPath
	generate.Stderr = &errb
	if err := generate.Run(); err != nil {
		fmt.Println(strings.Join(generate.Args, " "))
		fmt.Println(errb.String())
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = destinationPath
	tidy.Stderr = &errb
	if err := tidy.Run(); err != nil {
		fmt.Println(strings.Join(tidy.Args, " "))
		fmt.Println(errb.String())
	}
	clean := exec.Command("golangci-lint", "run", "./...", "--fix")
	clean.Dir = destinationPath
	_ = clean.Run()
	return nil
}
