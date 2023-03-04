package configs

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Project struct {
	Name           string         `yaml:"name"`
	Module         string         `yaml:"module"`
	GoVersion      string         `yaml:"goVersion"`
	Auth           bool           `yaml:"auth"`
	CI             string         `yaml:"ci"`
	Models         []*ModelConfig `yaml:"models"`
	GRPCEnabled    bool           `yaml:"gRPC"`
	GatewayEnabled bool           `yaml:"gateway"`
	RESTEnabled    bool           `yaml:"REST"`
	MakeEnabled    bool           `yaml:"make"`
	TaskEnabled    bool           `yaml:"task"`
}

func NewProject(configPath string) (*Project, error) {
	project := &Project{
		Name:        "",
		Module:      "",
		GoVersion:   "1.19",
		Auth:        true,
		CI:          "github",
		Models:      nil,
		GRPCEnabled: true,
		RESTEnabled: true,
		MakeEnabled: false,
		TaskEnabled: true,
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
		model.GRPCEnabled = project.GRPCEnabled
		model.GatewayEnabled = project.GatewayEnabled
		model.RESTEnabled = project.RESTEnabled
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
		validation.Field(&p.RESTEnabled),
		validation.Field(&p.GRPCEnabled),
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) ProtoPackage() string {
	return fmt.Sprintf("%spb", strcase.ToSnake(p.Name))
}
