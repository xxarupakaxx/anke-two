package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetTargets(ctx context.Context, questionnaireID int) ([]*model.Target, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	targets := make([]*model.Target, 0)

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Find(&targets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get targets :%w", err)
	}

	return targets, nil
}

func (repo *GormRepository) CreateTargets(ctx context.Context, questionnaireID int, targets []string) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteTargets(ctx context.Context, questionnaireID int) error {
	panic("implement me")
}
