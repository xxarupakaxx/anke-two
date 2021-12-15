package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetValidation(ctx context.Context, questionID int) (*model.Validation, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil,fmt.Errorf("failed to get db:%w", err)
	}
}

func (repo *GormRepository) CreateValidation(ctx context.Context, validation *model.Validation) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
}

func (repo *GormRepository) DeleteValidation(ctx context.Context, questionID int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
}

func (repo *GormRepository) UpdateValidation(ctx context.Context, validation *model.Validation) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
}

