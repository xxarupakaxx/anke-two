package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type IOption interface {
	InsertOption(ctx context.Context, lastID int, num int, body string) error
	UpdateOptions(ctx context.Context, options []string, questionID int) error
	DeleteOptions(ctx context.Context, questionID int) error
	GetOptions(ctx context.Context, questionIDs []int) ([]model.Options, error)
}
