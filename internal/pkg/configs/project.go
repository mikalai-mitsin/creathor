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
		GoVersion:      "1.24",
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

	for i, app := range project.Apps {
		app.Module = project.Module
		app.ProjectName = project.Name
		app.ProtoPackage = project.ProtoPackage()
		app.GRPCEnabled = project.GRPCEnabled
		app.HTTPEnabled = project.HTTPEnabled
		app.GatewayEnabled = project.GatewayEnabled
		app.KafkaEnabled = project.KafkaEnabled
		for i2, entity := range app.Entities {
			entity.Module = project.Module
			entity.ProjectName = project.Name
			entity.ProtoPackage = project.ProtoPackage()
			entity.GRPCEnabled = project.GRPCEnabled
			entity.HTTPEnabled = project.HTTPEnabled
			entity.GatewayEnabled = project.GatewayEnabled
			entity.KafkaEnabled = project.KafkaEnabled
			app.Entities[i2] = entity
		}
		app.ProjectConfig = project
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

func (p *Project) ErrsImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/errs"`, p.Module)
}
func (p *Project) LogImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/log"`, p.Module)
}
func (p *Project) KafkaImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/kafka"`, p.Module)
}

func (p *Project) UUIDImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/uuid"`, p.Module)
}

func (p *Project) UptraceImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/uptrace"`, p.Module)
}

func (p *Project) PostgresImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/postgres"`, p.Module)
}

func (p *Project) PointerImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/pointer"`, p.Module)
}

func (p *Project) HTTPImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/http"`, p.Module)
}

func (p *Project) GRPCImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/grpc"`, p.Module)
}

func (p *Project) ContainersImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/containers"`, p.Module)
}

func (p *Project) ConfigsImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/configs"`, p.Module)
}

func (p *Project) ClockImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/clock"`, p.Module)
}

func (p *Project) GatewayImportPath() string {
	return fmt.Sprintf(`"%s/internal/pkg/gateway"`, p.Module)
}
