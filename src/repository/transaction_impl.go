package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type ctxKey string

const txKey ctxKey = "transaction"


// Do Transaction用のメソッド
func (repo *GormRepository) Do(ctx context.Context, options *sql.TxOptions, callBack func(context.Context) error) error {
	txFunc := func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, txKey, tx)

		err := callBack(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	if options == nil {
		err := repo.db.Transaction(txFunc)
		if err != nil {
			return fmt.Errorf("failed to get transaction:%w", err)
		}
	} else {
		err := repo.db.Transaction(txFunc, options)
		if err != nil {
			return fmt.Errorf("failed to get transaction:%w", err)
		}
	}

	return nil
}
