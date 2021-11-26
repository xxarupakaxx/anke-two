package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/domain/repository/transaction"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
)

type response struct {
	repository.IRespondent
	repository.IQuestionnaire
	repository.IValidation
	repository.IScaleLabel
	repository.IResponse
	transaction.ITransaction
}

func NewResponse(IRespondent repository.IRespondent, IQuestionnaire repository.IQuestionnaire, IValidation repository.IValidation, IScaleLabel repository.IScaleLabel, IResponse repository.IResponse, ITransaction transaction.ITransaction) ResponseUsecase {
	return &response{IRespondent: IRespondent, IQuestionnaire: IQuestionnaire, IValidation: IValidation, IScaleLabel: IScaleLabel, IResponse: IResponse, ITransaction: ITransaction}
}

func (r *response) PostResponse(ctx context.Context, responses input.Responses) (output.PostResponse, error) {
	err
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
