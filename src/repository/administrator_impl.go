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

func (repo *GormRepository) GetAdmins(ctx context.Context, questionnaireIDs []int) ([]*model.Administrator, error) {
	panic("implement me")
}

func (repo *GormRepository) DeleteAdmin(ctx context.Context, questionnaireID int) {
	panic("implement me")
}
