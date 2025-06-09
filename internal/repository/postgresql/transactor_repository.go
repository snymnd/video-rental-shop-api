package postgresql

import (
	"context"
	"database/sql"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/util/logger"
)

type transactor struct {
	conn *sql.DB
}

func NewTxRepository(conn *sql.DB) *transactor {
	return &transactor{conn}
}

type TxKey struct {
	key string
}

type DBinf interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

const TRANSACTION_KEY = "txkey"

func (wt *transactor) WithTx(ctx context.Context, tFunc func(ctx context.Context) error) error {
	log := logger.GetLogger()
	// begin transaction
	tx, err := wt.conn.Begin()
	if err != nil {
		return customerrors.NewError(
			"transaction failed",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			log.Info("close transaction: %v", errTx)
		}
	}()

	// run callback
	if err := tFunc(context.WithValue(ctx, TxKey{TRANSACTION_KEY}, tx)); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return customerrors.NewError(
				"transaction failed",
				errRollback,
				customerrors.DatabaseExecutionError,
			)
		}
		return err
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return customerrors.NewError(
			"transaction failed",
			errCommit,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
