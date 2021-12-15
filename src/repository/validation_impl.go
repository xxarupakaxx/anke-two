package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetValidation(ctx context.Context, questionID int) (*model.Validation, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	var validation *model.Validation
	err = db.
		Where("question_id = ?", questionID).
		First(&validation).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get validation :%w", err)
	}

	return validation, err
}

func (repo *GormRepository) CreateValidation(ctx context.Context, validation *model.Validation) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	err = db.Create(&validation).Error
	if err != nil {
		return fmt.Errorf("failed to create validation :%w", err)
	}

	return nil
}

func (repo *GormRepository) DeleteValidation(ctx context.Context, questionID int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("question_id = ?", questionID).
		Delete(&model.Validation{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete validation :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

func (repo *GormRepository) UpdateValidation(ctx context.Context, validation *model.Validation) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("question_id = ?", validation.QuestionID).
		Updates(&validation)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update validation :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordUpdated
	}

	return nil
}
