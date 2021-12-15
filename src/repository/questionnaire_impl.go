package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
	"gorm.io/gorm"
)

func (repo *GormRepository) GetQuestionnaire(ctx context.Context, id int) (*model.Questionnaire, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}
	questionnaire := model.Questionnaire{}
	err = db.
		Where("id = ?", id).
		First(&questionnaire).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questionnaire :%w", err)
	}

	return &questionnaire, nil
}

func (repo *GormRepository) CreateQuestionnaire(ctx context.Context, questionnaire *model.Questionnaire) (int, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get db:%w", err)
	}
}

func (repo *GormRepository) UpdateQuestionnaire(ctx context.Context, questionnaire *model.Questionnaire) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
}

func (repo *GormRepository) DeleteQuestionnaire(ctx context.Context, id int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
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
