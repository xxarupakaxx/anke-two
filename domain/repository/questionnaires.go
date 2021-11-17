package repository

import (
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"gopkg.in/guregu/null.v4"
)

type IQuestionnaire interface {
	InsertQuestionnaire(questionnaires input.Questionnaires) (int, error)
	UpdateQuestionnaire(questionnaires input.Questionnaires) error
	DeleteQuestionnaire(questionnaires input.Questionnaires) error
	GetQuestionnaires(questionnaires input.Questionnaires) ([]output.QuestionnaireInfo, error)
	GetAdminQuestionnaires(questionnaires input.Questionnaires) ([]model.Questionnaires, error)
	GetQuestionnaireInfo(questionnaires input.Questionnaires) (output.QuestionnaireInfo, error)
	GetTargetedQuestionnaires(questionnaires input.Questionnaires) ([]output.TargetedQuestionnaire, error)
	GetQuestionnaireLimit(questionnaires input.Questionnaires) (null.Time, error)
	GetQuestionnaireLimitByResponseID(questionnaires input.Questionnaires)(null.Time,error)
	GetResponseReadPrivilegeInfoByResponseID(questionnaires input.Questionnaires) (*output.ResponseReadPrivilegeInfo, error)
	GetResponseReadPrivilegeInfoByQuestionnaireID(questionnaires input.Questionnaires) (*output.ResponseReadPrivilegeInfo, error)
}
