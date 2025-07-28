package main

import (
	"bytes"

	"github.com/mikalai-mitsin/creathor/internal/app/generator/auth"

	"github.com/mikalai-mitsin/creathor/internal/app/generator/layout"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	appgenerator "github.com/mikalai-mitsin/creathor/internal/app/generator/app"
	"github.com/mikalai-mitsin/creathor/internal/app/generator/pkg"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"

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
	if project.Auth {
		authGenerator := auth.NewGenerator(project)
		if err := authGenerator.Sync(); err != nil {
			return err
		}
	}
	for _, appConfig := range project.Apps {
		ap := &app.App{
			Config:      appConfig,
			Name:        appConfig.Name,
			Module:      project.Module,
			ProtoModule: project.ProtoPackage(),
			Entities:    make([]*app.BaseEntity, len(appConfig.Entities)),
			Auth:        appConfig.Auth,
		}
		for i, entity := range appConfig.Entities {
			baseEntity := &app.BaseEntity{
				Config:      entity,
				AppConfig:   &appConfig,
				Name:        entity.Name,
				Module:      project.Module,
				ProtoModule: project.ProtoPackage(),
				Entities: []*app.Entity{
					app.NewMainEntity(entity),
					app.NewFilterEntity(entity),
					app.NewCreateEntity(entity),
					app.NewUpdateEntity(entity),
				},
				Auth: project.Auth,
			}
			ap.Entities[i] = baseEntity
		}
		appGenerator := appgenerator.NewGenerator(ap)
		if err := appGenerator.Sync(); err != nil {
			return err
		}
	}
	if err := postInit(project); err != nil {
		return err
	}
	return nil
}

func postInit(project *configs.Project) error {
	fmt.Println("post init...")
	var errb bytes.Buffer
	if project.GRPCEnabled {
		bufUpdate := exec.Command("buf", "dep", "update")
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
	clean := exec.Command("task", "clean")
	clean.Dir = destinationPath
	fmt.Println(strings.Join(clean.Args, " "))
	if err := clean.Run(); err != nil {
		fmt.Println(errb.String())
	}
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
	return nil
}
