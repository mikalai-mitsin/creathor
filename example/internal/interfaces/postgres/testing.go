package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewMockPostgreSQL(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, error) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	return sqlxDB, mock, nil
}
