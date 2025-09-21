package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathKafkaRepositories() string {
	return fmt.Sprintf(`"%s/internal/app/%s/repositories/kafka/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasKafkaRepositories() string {
	return fmt.Sprintf("%sKafkaRepositories", m.LowerCamelName())
}

func (m *EntityConfig) EventProducerConstructorName() string {
	return fmt.Sprintf("New%s", m.EventProducerTypeName())
}

func (m *EntityConfig) EventProducerTypeName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) EventProducerInterfaceName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToLowerCamel(m.Name))
}
func (m *EntityConfig) GetEventProducerPrivateVariableName() string {
	return fmt.Sprintf("%sEventProducer", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) TopicName() string {
	return fmt.Sprintf(
		"%s.%s.%s.v1",
		strcase.ToSnake(m.AppConfig.ProjectConfig.Name),
		strcase.ToSnake(m.AppConfig.Name),
		strcase.ToSnake(m.Name),
	)
}
