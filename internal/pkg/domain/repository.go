package domain

import (
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func NewRepository(m *configs.DomainConfig) *Layer {
	layer := &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.RepositoryTypeName(),
		Variable: m.RepositoryVariableName(),
	}
	return layer
}
