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
	Name           string      `yaml:"name"`
	Module         string      `yaml:"module"`
	GoVersion      string      `yaml:"goVersion"`
	Auth           bool        `yaml:"auth"`
	CI             string      `yaml:"ci"`
	Apps           []AppConfig `yaml:"apps"`
	GRPCEnabled    bool        `yaml:"gRPC"`
	GatewayEnabled bool        `yaml:"gateway"`
	MakeEnabled    bool        `yaml:"make"`
	TaskEnabled    bool        `yaml:"task"`
	UptraceEnabled bool        `yaml:"uptrace"`
	KafkaEnabled   bool        `yaml:"kafka"`
	HTTPEnabled    bool        `yaml:"http"`
}

func NewProject(configPath string) (*Project, error) {
	project := &Project{
		Name:           "",
		Module:         "",
		GoVersion:      "1.20",
		Auth:           true,
		CI:             "github",
		Apps:           nil,
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
		project.Apps = append(project.Apps, AppConfig{
			Name:           "users",
			ProjectName:    project.Name,
			ProtoPackage:   project.ProtoPackage(),
			Auth:           project.Auth,
			HTTPEnabled:    project.HTTPEnabled,
			GRPCEnabled:    project.GRPCEnabled,
			GatewayEnabled: project.GatewayEnabled,
			KafkaEnabled:   project.KafkaEnabled,
			Entities: []EntityConfig{
				{
					Name:         "user",
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
				},
			},
		})
	}
	for i, app := range project.Apps {
		app.Module = project.Module
		app.Auth = project.Auth
		app.ProjectName = project.Name
		app.ProtoPackage = project.ProtoPackage()
		app.GRPCEnabled = project.GRPCEnabled
		app.HTTPEnabled = project.HTTPEnabled
		app.GatewayEnabled = project.GatewayEnabled
		for i2, entity := range app.Entities {
			entity.Module = project.Module
			entity.Auth = project.Auth
			entity.ProjectName = project.Name
			entity.ProtoPackage = project.ProtoPackage()
			entity.GRPCEnabled = project.GRPCEnabled
			entity.HTTPEnabled = project.HTTPEnabled
			entity.GatewayEnabled = project.GatewayEnabled
			app.Entities[i2] = entity
		}
		project.Apps[i] = app
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
		validation.Field(&p.Apps),
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
