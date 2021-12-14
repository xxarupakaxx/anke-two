package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) CreateAdmins(ctx context.Context, questionnaireID int, administrator []string) error {
	if questionnaireID < 0 {
		return ErrNotFormat
	}
	if len(administrator) == 0 {
		return nil
	}

	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	admins := make([]*model.Administrator, len(administrator))
	for _, a := range administrator {
		admins = append(admins, &model.Administrator{
			QuestionnaireID: questionnaireID,
			TraqID:          a,
		})
	}
	err = db.Create(&admins).Error
	if err != nil {
		return fmt.Errorf("failed to create administrators :%w", err)
	}

	return nil
}

func (repo *GormRepository) GetMyAdmins(ctx context.Context, traqID string) ([]*model.Administrator, error) {
	if traqID == "" {
		return nil, nil
	}

	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	myAdmins := make([]*model.Administrator, 0)

	err = db.
		Where("traq_id", traqID).
		Find(&myAdmins).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get My Admins  :%w", err)
	}

	return myAdmins, nil
}

func (repo *GormRepository) DeleteAdmin(ctx context.Context, questionnaireID int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&model.Administrator{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete admins :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}
