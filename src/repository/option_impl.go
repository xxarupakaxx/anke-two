package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetOptions(ctx context.Context, questionIDs []int) ([]*model.Option, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	options := make([]*model.Option, len(questionIDs))

	err = db.
		Where("question_id IN  ? ", questionIDs).
		Find(&options).Error
	if err != nil {
		return nil, fmt.Errorf("failed get options :%w", err)
	}

	return options, nil
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

func (repo *GormRepository) DeleteOption(ctx context.Context, questionID int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("question_id = ?", questionID).
		Delete(&model.Option{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete option :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

func (repo *GormRepository) UpdateOption(ctx context.Context, option *model.Option) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("id = ?", option.ID).
		Updates(&option)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to updates option:%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordUpdated
	}

	return nil
}
