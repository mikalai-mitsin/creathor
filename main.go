package main

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

const version = "0.4.0"

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
	project, err := NewProject()
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
		model.Module = project.Module
		model.Auth = project.Auth
		model.ProjectName = project.Name
		model.ProtoPackage = project.ProtoPackage()
		if err := CreateCRUD(model); err != nil {
			return err
		}
	}
	if err := postInit(); err != nil {
		return err
	}
	return nil
}

func postInit() error {
	fmt.Println("post init...")
	var errb bytes.Buffer
	generate := exec.Command("go", "generate", "./...")
	generate.Dir = destinationPath
	generate.Stderr = &errb
	fmt.Println(strings.Join(generate.Args, " "))
	if err := generate.Run(); err != nil {
		fmt.Println(errb.String())
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = destinationPath
	tidy.Stderr = &errb
	fmt.Println(strings.Join(tidy.Args, " "))
	if err := tidy.Run(); err != nil {
		fmt.Println(errb.String())
	}
	swag := exec.Command("swag", "init", "-d", "./internal/interfaces/rest", "-g", "server.go", "--parseDependency", "-o", "./api", "-ot", "yaml")
	swag.Dir = destinationPath
	swag.Stderr = &errb
	fmt.Println(strings.Join(swag.Args, " "))
	if err := swag.Run(); err != nil {
		fmt.Println(errb.String())
	}
	buf := exec.Command("buf", "generate")
	buf.Dir = destinationPath
	buf.Stderr = &errb
	fmt.Println(strings.Join(buf.Args, " "))
	if err := buf.Run(); err != nil {
		fmt.Println(errb.String())
	}
	clean := exec.Command("golangci-lint", "run", "./...", "--fix")
	clean.Dir = destinationPath
	fmt.Println(strings.Join(clean.Args, " "))
	_ = clean.Run()
	return nil
}
