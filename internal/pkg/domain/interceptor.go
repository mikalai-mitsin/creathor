package domain

import (
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func NewUseCase(m *configs.DomainConfig) *Layer {
	usecase := &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.UseCaseTypeName(),
		Variable: m.UseCaseVariableName(),
	}
	return usecase
}
