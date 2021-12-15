package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
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

