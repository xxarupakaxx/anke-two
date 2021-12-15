package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
)

func (repo *GormRepository) GetRespondent(ctx context.Context, questionnaireID int) ([]*model.Respondent, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	respondents := make([]*model.Respondent, 0)
	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents :%w", err)
	}

	return respondents, nil
}

func (repo *GormRepository) GetMyRespondent(ctx context.Context, questionnaireID int, traPID string) ([]*model.Respondent, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	respondents := make([]*model.Respondent, 0)

	err = db.
		Where("questionnaire_id = ? AND user_traqid = ?", questionnaireID, traPID).
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents :%w", err)
	}

	return respondents, nil
}

func (repo *GormRepository) CreateRespondent(ctx context.Context, respondent *model.Respondent) (int, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get db:%w", err)
	}

	err = db.Create(&respondent).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create repsondent :%w", err)
	}

	return respondent.ResponseID, nil
}

func (repo *GormRepository) DeleteRespondent(ctx context.Context, responseID int) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateRespondent(ctx context.Context, respondent *model.Respondent) error {
	panic("implement me")
}
