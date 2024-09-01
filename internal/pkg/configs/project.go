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
	Name           string          `yaml:"name"`
	Module         string          `yaml:"module"`
	GoVersion      string          `yaml:"goVersion"`
	Auth           bool            `yaml:"auth"`
	CI             string          `yaml:"ci"`
	Domains        []*DomainConfig `yaml:"domains"`
	GRPCEnabled    bool            `yaml:"gRPC"`
	GatewayEnabled bool            `yaml:"gateway"`
	MakeEnabled    bool            `yaml:"make"`
	TaskEnabled    bool            `yaml:"task"`
	UptraceEnabled bool            `yaml:"uptrace"`
	KafkaEnabled   bool            `yaml:"kafka"`
}

func NewProject(configPath string) (*Project, error) {
	project := &Project{
		Name:           "",
		Module:         "",
		GoVersion:      "1.20",
		Auth:           true,
		CI:             "github",
		Domains:        nil,
		GRPCEnabled:    true,
		GatewayEnabled: false,
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
		project.Domains = append(project.Domains, &DomainConfig{
			Model:        "user",
			Module:       project.Module,
			ProjectName:  project.Name,
			ProtoPackage: project.ProtoPackage(),
			Auth:         project.Auth,
			Params: []*Param{
				{Name: "FirstName", Type: "string", Search: true},
				{Name: "LastName", Type: "string", Search: true},
				{Name: "Password", Type: "string", Search: false},
				{Name: "Email", Type: "string", Search: true},
				{Name: "GroupID", Type: "entities.GroupID", Search: false},
			},
			GRPCEnabled:    project.GRPCEnabled,
			GatewayEnabled: project.GatewayEnabled,
			KafkaEnabled:   project.KafkaEnabled,
		})
	}
	for _, domain := range project.Domains {
		domain.Module = project.Module
		domain.Auth = project.Auth
		domain.ProjectName = project.Name
		domain.ProtoPackage = project.ProtoPackage()
		domain.GRPCEnabled = project.GRPCEnabled
		domain.GatewayEnabled = project.GatewayEnabled
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
		validation.Field(&p.Domains),
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
