package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	repository2 "github.com/xxarupkaxx/anke-two/interfaces/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
)

type Result struct {
	repository2.IRespondent
	repository2.IAdministrator
	repository2.IQuestionnaire
}

func NewResult(IRespondent repository2.IRespondent, IAdministrator repository2.IAdministrator, IQuestionnaire repository2.IQuestionnaire) *Result {
	return &Result{IRespondent: IRespondent, IAdministrator: IAdministrator, IQuestionnaire: IQuestionnaire}
}

func (r *Result) GetResults(ctx context.Context, results input.GetResults) ([]model.RespondentDetail, error) {
	outputResults, err := r.GetRespondentDetails(ctx, results.QuestionnaireID, results.Sort)
	if err != nil {
		return nil, err
	}
	return outputResults, err
}

type ResultUsecase interface {
	GetResults(ctx context.Context, results input.GetResults) ([]model.RespondentDetail, error)
}
