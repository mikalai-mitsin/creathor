package main

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/018bf/creathor/models"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var version = "0.4.0"

var (
	destinationPath = "."
	configPath      = "./creathor.yaml"
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
			&cli.StringFlag{
				Name:        "c",
				Usage:       "config path",
				Destination: &configPath,
				Required:    false,
				Value:       configPath,
			},
		},
		Action: initProject,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initProject(ctx *cli.Context) error {
	project, err := models.NewProject(path.Join(destinationPath, configPath))
	if err != nil {
		return err
	}
	if err := CreateLayout(project); err != nil {
		return err
	}
	if err := CreateCI(project); err != nil {
		return err
	}
	if err := CreateDI(project); err != nil {
		return err
	}
	if err := CreateBuild(project); err != nil {
		return err
	}
	if err := CreateDeployment(project); err != nil {
		return err
	}
	for _, model := range project.Models {
		if err := CreateCRUD(model); err != nil {
			return err
		}
	}
	if err := postInit(project); err != nil {
		return err
	}
	return nil
}

func postInit(project *models.Project) error {
	fmt.Println("post init...")
	var errb bytes.Buffer
	generate := exec.Command("go", "generate", "./...")
	generate.Dir = destinationPath
	generate.Stderr = &errb
	fmt.Println(strings.Join(generate.Args, " "))
	if err := generate.Run(); err != nil {
		fmt.Println(errb.String())
	}
	if project.RESTEnabled {
		swag := exec.Command("swag", "init", "-d", "./internal/interfaces/rest", "-g", "server.go", "--parseDependency", "-o", "./api/rest", "-ot", "json")
		swag.Dir = destinationPath
		swag.Stderr = &errb
		fmt.Println(strings.Join(swag.Args, " "))
		if err := swag.Run(); err != nil {
			fmt.Println(errb.String())
		}
	}
	if project.GRPCEnabled {
		bufUpdate := exec.Command("buf", "mod", "update")
		bufUpdate.Dir = path.Join(destinationPath, "api", "proto")
		bufUpdate.Stderr = &errb
		fmt.Println(strings.Join(bufUpdate.Args, " "))
		if err := bufUpdate.Run(); err != nil {
			fmt.Println(errb.String())
		}
		bufGenerate := exec.Command("buf", "generate")
		bufGenerate.Dir = destinationPath
		bufGenerate.Stderr = &errb
		fmt.Println(strings.Join(bufGenerate.Args, " "))
		if err := bufGenerate.Run(); err != nil {
			fmt.Println(errb.String())
		}
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = destinationPath
	tidy.Stderr = &errb
	fmt.Println(strings.Join(tidy.Args, " "))
	if err := tidy.Run(); err != nil {
		fmt.Println(errb.String())
	}
	clean := exec.Command("golangci-lint", "run", "./...", "--fix")
	clean.Dir = destinationPath
	fmt.Println(strings.Join(clean.Args, " "))
	if err := clean.Run(); err != nil {
		fmt.Println(errb.String())
	}
	return nil
}
