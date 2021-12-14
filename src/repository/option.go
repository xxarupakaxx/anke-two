package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Option interface {
	GetOptions(ctx context.Context, questionIDs []int) ([]*model.Option, error)
	CreateOption(ctx context.Context, option *model.Option) error
	DeleteOption(ctx context.Context, questionID int) error
	UpdateOption(ctx context.Context, option *model.Option) error
}
