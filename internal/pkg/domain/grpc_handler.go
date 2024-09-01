package domain

import (
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func NewGRPCHandler(m *configs.DomainConfig) *Layer {
	return &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.GRPCHandlerTypeName(),
		Variable: m.GRPCHandlerVariableName(),
	}
}
