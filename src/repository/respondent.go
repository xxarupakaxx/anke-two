package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Respondent interface {
	GetRespondent(ctx context.Context, questionnaireID int) ([]*model.Respondent, error)
	GetMyRespondent(ctx context.Context, questionnaireID int, traPID string) ([]*model.Respondent, error)
	CreateRespondent(ctx context.Context, respondent *model.Respondent) (int, error)
	DeleteRespondent(ctx context.Context, responseID int) error
	UpdateRespondent(ctx context.Context, respondent *model.Respondent) error
}
