package repository

import (
	"context"
	"database/sql"
)

type Transaction interface {
	Do(ctx context.Context, options *sql.TxOptions, callBack func(context.Context) error) error
}
