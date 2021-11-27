package database

import (
	"context"
	"database/sql"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
)

type ctxKey string

const (
	txKey ctxKey = "transaction"
)

type Tx struct {
	db *gorm.DB
}

func NewTransaction(db *gorm.DB) *Tx {
	return &Tx{db: db}
}

func (t *Tx) Do(ctx context.Context, options *sql.TxOptions, f func(context.Context) error) error {
	fc := func(txx *gorm.DB) error {
		ctx = context.WithValue(ctx, txKey, txx)

		err := f(ctx)

		if err != nil {
			return err
		}
		return nil
	}
	if options == nil {
		err := t.db.Transaction(fc)
		if err != nil {
			return err
		}
	} else {
		err := t.db.Transaction(fc, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTx(ctx context.Context) (*gorm.DB, error) {
	iDB, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok {
		return nil, model.ErrInvalidTx
	}

	return iDB.Session(&gorm.Session{
		Context: ctx,
	}), nil
}
