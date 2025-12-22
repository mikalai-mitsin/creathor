package configs

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

func (m *EntityConfig) ImportPathKafkaHandlers() string {
	return fmt.Sprintf(`"%s/internal/app/%s/handlers/kafka/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasKafkaHandlers() string {
	return fmt.Sprintf("%sKafkaHandlers", m.LowerCamelName())
}

func (m *EntityConfig) KafkaHandlerConstructorName() string {
	return fmt.Sprintf("New%s", m.KafkaHandlerTypeName())
}

func (m *EntityConfig) KafkaHandlerTypeName() string {
	return fmt.Sprintf("%sHandler", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) KafkaConsumerGroup() string {
	return fmt.Sprintf(
		"%s.%s.%s",
		strcase.ToSnake(m.AppConfig.ProjectConfig.Name),
		strcase.ToSnake(m.AppConfig.Name),
		strcase.ToSnake(m.Name),
	)
}

func (m *EntityConfig) GetKafkaHandlerPrivateVariableName() string {
	return fmt.Sprintf("kafka%sHandler", strcase.ToCamel(m.Name))
}
