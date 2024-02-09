package configs

import (
	"fmt"
	"log"
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
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
	UptraceEnabled bool           `yaml:"uptrace"`
	KafkaEnabled   bool           `yaml:"kafka"`
}

func NewProject(configPath string) (*Project, error) {
	project := &Project{
		Name:           "",
		Module:         "",
		GoVersion:      "1.20",
		Auth:           true,
		CI:             "github",
		Models:         nil,
		GRPCEnabled:    true,
		GatewayEnabled: false,
		RESTEnabled:    true,
		MakeEnabled:    false,
		TaskEnabled:    true,
		UptraceEnabled: false,
		KafkaEnabled:   false,
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, project); err != nil {
		log.Fatalf("error: %v", err)
	}
	if project.Auth {
		project.Models = append(project.Models, &ModelConfig{
			Model:        "User",
			Module:       project.Module,
			ProjectName:  project.Name,
			ProtoPackage: project.ProtoPackage(),
			Auth:         project.Auth,
			Params: []*Param{
				{Name: "FirstName", Type: "string", Search: true},
				{Name: "LastName", Type: "string", Search: true},
				{Name: "Password", Type: "string", Search: false},
				{Name: "Email", Type: "string", Search: true},
				{Name: "GroupID", Type: "models.GroupID", Search: false},
			},
			GRPCEnabled:    project.GRPCEnabled,
			GatewayEnabled: project.GatewayEnabled,
			RESTEnabled:    project.RESTEnabled,
			KafkaEnabled:   project.KafkaEnabled,
		})
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
