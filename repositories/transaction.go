package repositories

import (
	"context"
	"database/sql"
)

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) (interface{}, error)) (interface{}, error)
}

type transactor struct {
	Db *sql.DB
}

func NewTransactor(db *sql.DB) *transactor {
	return &transactor{
		Db: db,
	}
}

func (t *transactor) WithinTransaction(ctx context.Context, tFunc func(txCtx context.Context) (interface{}, error)) (interface{}, error) {
	tx, err := t.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	res, err := tFunc(injectTx(ctx, tx))
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, errCommit
	}
	return res, nil
}
