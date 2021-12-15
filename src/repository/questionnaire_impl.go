package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
	"gorm.io/gorm"
)

func (repo *GormRepository) GetQuestionnaire(ctx context.Context, id int) (*model.Questionnaire, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateQuestionnaire(ctx context.Context, questionnaire *model.Questionnaire) (int, error) {
	panic("implement me")
}

func (repo *GormRepository) UpdateQuestionnaire(ctx context.Context, questionnaire *model.Questionnaire) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteQuestionnaire(ctx context.Context, id int) error {
	panic("implement me")
}

func setUpResSharedTo(db *gorm.DB) error {
	resSharedTypes := []model.ResSharedTo{
		{
			Name: "administrators",
		},
		{
			Name: "respondents",
		},
		{
			Name: "public",
		},
	}
	for _, resSharedType := range resSharedTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", resSharedType.Name).
			FirstOrCreate(&resSharedType).Error
		if err != nil {
			return fmt.Errorf("failed to create resSharedType:%w", err)
		}
	}

	return nil
}
