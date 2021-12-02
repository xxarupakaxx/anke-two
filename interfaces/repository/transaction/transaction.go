package transaction

import (
	"context"
	"database/sql"
)

type ITransaction interface {
	Do(context.Context, *sql.TxOptions, func(context.Context) error) error
}

