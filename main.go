package main

import (
	"bytes"

	"github.com/mikalai-mitsin/creathor/internal/app/generator/layout"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg"

	"github.com/mikalai-mitsin/creathor/internal/app/generator/app"

	"github.com/iancoleman/strcase"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/urfave/cli/v2"
)

var version string

var (
	destinationPath = "."
	configPath      = "./creathor.yaml"
)

func main() {
	application := &cli.App{
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
	strcase.ConfigureAcronym("UUID", "uuid")
	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initProject(_ *cli.Context) error {
	project, err := configs.NewProject(path.Join(destinationPath, configPath))
	if err != nil {
		return err
	}
	layoutGenerator := layout.NewGenerator(project)
	if err := layoutGenerator.Sync(); err != nil {
		return err
	}
	pkgGenerator := pkg.NewGenerator(project)
	if err := pkgGenerator.Sync(); err != nil {
		return err
	}
	for _, appConfig := range project.Apps {
		for i, entity := range appConfig.Entities {
			appConfig.Entities[i].AppConfig = &appConfig
			appConfig.Entities[i].Entities = append(appConfig.Entities[i].Entities,
				configs.NewMainEntity(entity),
				configs.NewFilterEntity(entity),
				configs.NewCreateEntity(entity),
				configs.NewUpdateEntity(entity),
				configs.NewDeleteEntity(entity),
			)
		}
		appGenerator := app.NewGenerator(&appConfig)
		if err := appGenerator.Sync(); err != nil {
			return err
		}
	}
	if err := postInit(project); err != nil {
		return err
	}
	return nil
}

type command struct {
	name string
	args []string
	dir  string
}

func postInit(project *configs.Project) error {
	fmt.Println("post init...")

	commands := []command{
		{name: "task", args: []string{"clean"}, dir: destinationPath},
		{name: "go", args: []string{"generate", "./..."}, dir: destinationPath},
		{name: "go", args: []string{"mod", "tidy"}, dir: destinationPath},
		{name: "task", args: []string{"docs"}, dir: destinationPath},
	}

	if project.GRPCEnabled || project.KafkaEnabled {
		bufCommands := []command{
			{name: "buf", args: []string{"dep", "update"}, dir: path.Join(destinationPath, "api", "proto")},
			{name: "buf", args: []string{"generate"}, dir: destinationPath},
		}
		commands = append(bufCommands, commands...)
	}

	for _, cmd := range commands {
		if err := runCommand(cmd); err != nil {
			fmt.Printf("Warning: command failed: %v\n", err)
		}
	}

	return nil
}

func runCommand(cmd command) error {
	execCmd := exec.Command(cmd.name, cmd.args...)
	execCmd.Dir = cmd.dir

	var errb bytes.Buffer
	execCmd.Stderr = &errb

	fmt.Println(strings.Join(append([]string{cmd.name}, cmd.args...), " "))

	if err := execCmd.Run(); err != nil {
		fmt.Println(errb.String())
		return err
	}

	return nil
}
