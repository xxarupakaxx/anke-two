package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type IAdministrator interface {
	InsertAdministrator(ctx context.Context, questionnaireID int, administrators []string) error
	DeleteAdministrators(ctx context.Context, questionnaireID int) error
	GetAdministrators(ctx context.Context, questionnaireIDs []int) ([]model.Administrators, error)
	CheckQuestionnaireAdmin(ctx context.Context, userID string, questionnaireID int) (bool, error)
}
