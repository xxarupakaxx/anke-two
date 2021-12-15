package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetTargets(ctx context.Context, questionnaireID int) ([]*model.Target, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateTargets(ctx context.Context, questionnaireID int, targets []string) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteTargets(ctx context.Context, questionnaireID int) error {
	panic("implement me")
}

