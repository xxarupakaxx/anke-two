package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) CreateAdmins(ctx context.Context, questionnaireID int, administrator []string) error {
	panic("implement me")
}

func (repo *GormRepository) GetAdmins(ctx context.Context, questionnaireIDs []int) ([]*model.Administrator, error) {
	panic("implement me")
}

func (repo *GormRepository) DeleteAdmin(ctx context.Context, questionnaireID int) {
	panic("implement me")
}

