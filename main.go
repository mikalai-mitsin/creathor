package main

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/018bf/creathor/internal/generators/auth"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/018bf/creathor/internal/generators/app"
	"github.com/018bf/creathor/internal/generators/pkg"

	"github.com/018bf/creathor/internal/domain"

	"github.com/018bf/creathor/internal/configs"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

var version string

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
	strcase.ConfigureAcronym("UUID", "uuid")
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initProject(ctx *cli.Context) error {
	project, err := configs.NewProject(path.Join(destinationPath, configPath))
	if err != nil {
		return err
	}
	if err := CreateLayout(project); err != nil {
		return err
	}
	if err := CreateCI(project); err != nil {
		return err
	}
	if err := CreateBuild(project); err != nil {
		return err
	}
	if err := CreateDeployment(project); err != nil {
		return err
	}
	pkgGenerator := pkg.NewGenerator(project)
	if err := pkgGenerator.Sync(); err != nil {
		return err
	}
	authGenerator := auth.NewGenerator(project)
	if err := authGenerator.Sync(); err != nil {
		return err
	}
	for _, m := range project.Domains {
		d := &domain.Domain{
			Name:        m.Model,
			Module:      project.Module,
			ProtoModule: project.ProtoPackage(),
			Models: []*domain.Model{
				domain.NewMainModel(m),
				domain.NewFilterModel(m),
				domain.NewCreateModel(m),
				domain.NewUpdateModel(m),
			},
			UseCase:     domain.NewUseCase(m),
			Repository:  domain.NewRepository(m),
			Interceptor: domain.NewInterceptor(m),
			GRPCHandler: domain.NewGRPCHandler(m),
			Auth:        project.Auth,
		}
		appGenerator := app.NewGenerator(d)
		if err := appGenerator.Sync(); err != nil {
			return err
		}
	}
	for _, model := range project.Domains {
		if err := CreateCRUD(model); err != nil {
			return err
		}
	}
	if err := postInit(project); err != nil {
		return err
	}
	if err := RenderTests(project); err != nil {
		return err
	}
	return nil
}

func postInit(project *configs.Project) error {
	fmt.Println("post init...")
	var errb bytes.Buffer
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
	goLines := exec.Command("golines", ".", "-w", "--ignore-generated")
	goLines.Dir = destinationPath
	fmt.Println(strings.Join(goLines.Args, " "))
	if err := goLines.Run(); err != nil {
		fmt.Println(errb.String())
	}
	clean := exec.Command("golangci-lint", "run", "./...", "--fix")
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
