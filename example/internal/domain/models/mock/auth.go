package mock_models // nolint:stylecheck

import (
	"testing"

	"github.com/018bf/example/internal/domain/models"

	"github.com/jaswdr/faker"
)

func NewToken(t *testing.T) models.Token {
	t.Helper()
	return models.Token(faker.New().Internet().Password())
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
		Email:    faker.New().Internet().Email(),
		Password: faker.New().Internet().Password(),
	}
}
