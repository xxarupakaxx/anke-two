package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetOptions(ctx context.Context, questionIDs []int) ([]*model.Option, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateOption(ctx context.Context, option *model.Option) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	err = db.Create(&option).Error
	if err != nil {
		return fmt.Errorf("failed to create :%w", err)
	}

	return nil
}

func (repo *GormRepository) DeleteOption(Ctx context.Context, questionID int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateOption(ctx context.Context, option *model.Option) error {
	panic("implement me")
}
