package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetQuestions(ctx context.Context, questionnaireID int) (*model.Question, error) {
	panic("implement me")
}

func (repo *GormRepository) CreateQuestion(ctx context.Context, question *model.Question) error {
	panic("implement me")
}

func (repo *GormRepository) DeleteQuestion(ctx context.Context, id int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateQuestion(ctx context.Context, question *model.Question) error {
	panic("implement me")
}
