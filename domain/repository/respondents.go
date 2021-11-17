package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

type IRespondent interface {
	InsertRespondent(ctx context.Context, userID string, questionnaireID int, submittedAt null.Time) (int, error)
	UpdateSubmittedAt(ctx context.Context, responseID int) error
	DeleteRespondent(ctx context.Context, responseID int) error
	GetRespondent(ctx context.Context, responseID int) (*model.Respondents, error)
	GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]model.RespondentInfo, error)
	GetRespondentDetail(ctx context.Context, responseID int) (model.RespondentDetail, error)
	GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]model.RespondentDetail, error)
	GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]model.Respondents, error)
	CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error)
}