package mock_models // nolint:stylecheck

import (
	"github.com/018bf/example/internal/domain/models"
	"testing"

	"syreclabs.com/go/faker"
)

func NewToken(t *testing.T) models.Token {
	t.Helper()
	return models.Token(faker.Internet().Password(30, 60))
}

func NewTokenPair(t *testing.T) *models.TokenPair {
	t.Helper()
	return &models.TokenPair{
		Access:  NewToken(t),
		Refresh: NewToken(t),
	}
}

func NewLogin(t *testing.T) *models.Login {
	t.Helper()
	return &models.Login{
		Email:    faker.Internet().Email(),
		Password: faker.Internet().Password(6, 12),
	}
}
