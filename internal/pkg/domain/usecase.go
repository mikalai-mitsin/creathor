package domain

import (
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func NewService(m *configs.DomainConfig) *Layer {
	layer := &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.ServiceTypeName(),
		Variable: m.ServiceVariableName(),
	}
	return layer
}
