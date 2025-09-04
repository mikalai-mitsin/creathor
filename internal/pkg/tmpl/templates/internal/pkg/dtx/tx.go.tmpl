package dtx

import (
	"database/sql"
	"errors"
)

//go:generate mockgen -source=tx.go -package=dtx -destination=tx_mock.go

type TX interface {
	GetSQLTx() *sql.Tx
	Commit() error
	Rollback() error
}

type DTX struct {
	sqlTx *sql.Tx
}

func NewTX() *DTX {
	return &DTX{sqlTx: nil}
}

func NewTXWithSQL(tx *sql.Tx) *DTX {
	return &DTX{sqlTx: tx}
}

func (tx DTX) GetSQLTx() *sql.Tx {
	return tx.sqlTx
}

func (tx DTX) Commit() error {
	if tx.sqlTx != nil {
		if err := tx.sqlTx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (tx DTX) Rollback() error {
	if tx.sqlTx != nil {
		if err := tx.sqlTx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return nil
			}
			return err
		}
	}
	return nil
}
