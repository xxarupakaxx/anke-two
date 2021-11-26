package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
)

type user struct {
	repository.IRespondent
	repository.IQuestionnaire
	repository.ITarget
	repository.IAdministrator
}

func NewUser(IRespondent repository.IRespondent, IQuestionnaire repository.IQuestionnaire, ITarget repository.ITarget, IAdministrator repository.IAdministrator) UsersUsecase {
	return &user{IRespondent: IRespondent, IQuestionnaire: IQuestionnaire, ITarget: ITarget, IAdministrator: IAdministrator}
}

func (u *user) GetUsersMe(ctx context.Context, me input.GetMe) output.GetMe {
	opUser := output.GetMe{TraqID: me.UserID}
	return opUser
}

func (u *user) GetMyResponses(ctx context.Context, me input.GetMe) ([]model.RespondentInfo, error) {
	myResponses, err := u.IRespondent.GetRespondentInfos(ctx, me.UserID)
	if err != nil {
		return nil, err
	}

	return myResponses, nil
}

func (u *user) GetMyResponsesByID(ctx context.Context, response input.GetMyResponse) ([]model.RespondentInfo, error) {
	panic("implement me")
}

func (u *user) GetTargetedQuestionnaire(ctx context.Context, request input.GetTargetedQuestionnaire) ([]model.TargetedQuestionnaire, error) {
	panic("implement me")
}

func (u *user) GetMyQuestionnaire(ctx context.Context, me input.GetMe) ([]output.QuestionnaireInfo, error) {
	panic("implement me")
}

func (u *user) GetTargetedQuestionnairesByID(ctx context.Context, qid input.GetTargetsByTraQID) ([]model.TargetedQuestionnaire, error) {
	op, err := u.IQuestionnaire.GetTargetedQuestionnaires(ctx, qid.TraQID, qid.Answered, qid.Sort)
	if err != nil {
		return nil, err
	}

	return op, nil
}

type UsersUsecase interface {
	GetUsersMe(ctx context.Context, me input.GetMe) output.GetMe
	GetMyResponses(ctx context.Context, me input.GetMe) ([]model.RespondentInfo, error)
	GetMyResponsesByID(ctx context.Context, response input.GetMyResponse) ([]model.RespondentInfo, error)
	GetTargetedQuestionnaire(ctx context.Context, request input.GetTargetedQuestionnaire) ([]model.TargetedQuestionnaire, error)
	GetMyQuestionnaire(ctx context.Context, me input.GetMe) ([]output.QuestionnaireInfo, error)
	GetTargetedQuestionnairesByID(ctx context.Context, qid input.GetTargetsByTraQID) ([]model.TargetedQuestionnaire, error)
}
