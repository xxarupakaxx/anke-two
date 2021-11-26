package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
)

type response struct {
	repository.IRespondent
	repository.IQuestionnaire
	repository.IValidation
	repository.IScaleLabel
	repository.IResponse
}

func NewResponse(IRespondent repository.IRespondent, IQuestionnaire repository.IQuestionnaire, IValidation repository.IValidation, IScaleLabel repository.IScaleLabel, IResponse repository.IResponse) ResponseUsecase {
	return &response{IRespondent: IRespondent, IQuestionnaire: IQuestionnaire, IValidation: IValidation, IScaleLabel: IScaleLabel, IResponse: IResponse}
}

func (r *response) PostResponse(ctx context.Context, responses input.Responses) (output.PostResponse, error) {

	panic("implement me")
}

func (r *response) GetResponse(ctx context.Context, getResponse input.GetResponse) (model.RespondentDetail, error) {
	panic("implement me")
}

func (r *response) EditResponse(ctx context.Context, editResponse input.EditResponse) error {
	panic("implement me")
}

func (r *response) DeleteResponse(ctx context.Context, deleteResponse input.DeleteResponse) error {
	panic("implement me")
}

type ResponseUsecase interface {
	PostResponse(ctx context.Context, responses input.Responses) (output.PostResponse, error)
	GetResponse(ctx context.Context, getResponse input.GetResponse) (model.RespondentDetail, error)
	EditResponse(ctx context.Context, editResponse input.EditResponse) error
	DeleteResponse(ctx context.Context, deleteResponse input.DeleteResponse) error
}
