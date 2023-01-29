package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Project struct {
	Name      string   `yaml:"name"`
	Module    string   `yaml:"module"`
	GoVersion string   `yaml:"goVersion"`
	Auth      bool     `yaml:"auth"`
	CI        string   `yaml:"ci"`
	Models    []*Model `yaml:"models"`
}

func NewProject(configPath string) (*Project, error) {
	project := &Project{
		Name:      "",
		Module:    "",
		GoVersion: "",
		Auth:      false,
		CI:        "",
		Models:    nil,
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, project); err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, model := range project.Models {
		model.Module = project.Module
		model.Auth = project.Auth
		model.ProjectName = project.Name
		model.ProtoPackage = project.ProtoPackage()
	}
	return project, nil
}

func (p *Project) Validate() error {
	err := validation.ValidateStruct(
		p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Module, validation.Required),
		validation.Field(&p.GoVersion, validation.Required),
		validation.Field(&p.Auth, validation.Required),
		validation.Field(&p.CI),
		validation.Field(&p.Models),
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) ProtoPackage() string {
	return fmt.Sprintf("%spb", strcase.ToSnake(p.Name))
}
