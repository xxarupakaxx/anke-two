package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

type IQuestionnaire interface {
	InsertQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error)
	UpdateQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error
	DeleteQuestionnaire(ctx context.Context, questionnaireID int) error
	GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, nonTargeted bool) ([]model.QuestionnaireInfo, int, error)
	GetAdminQuestionnaires(ctx context.Context, userID string) ([]model.Questionnaires, error)
	GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (*model.ReturnQuestionnaires, []string, []string, []string, error)
	GetTargetedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]model.TargetedQuestionnaire, error)
	GetQuestionnaireLimit(ctx context.Context, questionnaireID int) (null.Time, error)
	GetQuestionnaireLimitByResponseID(ctx context.Context, responseID int) (null.Time, error)
	GetResponseReadPrivilegeInfoByResponseID(ctx context.Context, userID string, responseID int) (*model.ResponseReadPrivilegeInfo, error)
	GetResponseReadPrivilegeInfoByQuestionnaireID(ctx context.Context, userID string, questionnaireID int) (*model.ResponseReadPrivilegeInfo, error)
}
