package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetValidation(ctx context.Context, questionID int) (*model.Validation, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateValidation(ctx context.Context, validation *model.Validation) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteValidation(ctx context.Context, questionID int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateValidation(ctx context.Context, validation *model.Validation) error {
	panic("implement me")
}

