package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
)

type result struct {
	repository.IRespondent
	repository.IAdministrator
	repository.IQuestionnaire
}

func NewResult(IRespondent repository.IRespondent, IAdministrator repository.IAdministrator, IQuestionnaire repository.IQuestionnaire) ResultUsecase {
	return &result{IRespondent: IRespondent, IAdministrator: IAdministrator, IQuestionnaire: IQuestionnaire}
}

func (r *result) GetResults(ctx context.Context, results input.GetResults) ([]model.RespondentDetail, error) {
	outputResults, err := r.GetRespondentDetails(ctx, results.QuestionnaireID, results.Sort)
	if err != nil {
		return nil, err
	}
	return outputResults, err
}

type ResultUsecase interface {
	GetResults(ctx context.Context, results input.GetResults) ([]model.RespondentDetail, error)
}
