package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetQuestions(ctx context.Context, questionnaireID int) ([]*model.Question, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	questions := make([]*model.Question, 0)

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions :%w", err)
	}

	return questions, nil
}

func (repo *GormRepository) CreateQuestion(ctx context.Context, question *model.Question) (int, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get db:%w", err)
	}

	err = db.Create(&question).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create question:%w", err)
	}

	return question.ID, err
}

func (repo *GormRepository) DeleteQuestion(ctx context.Context, id int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateQuestion(ctx context.Context, question *model.Question) error {
	panic("implement me")
}
