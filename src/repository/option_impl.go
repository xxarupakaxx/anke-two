package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetOptions(ctx context.Context, questionIDs []int) ([]*model.Option, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateOption(ctx context.Context, option *model.Option) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteOption(Ctx context.Context, questionID int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateOption(ctx context.Context, option *model.Option) error {
	panic("implement me")
}

