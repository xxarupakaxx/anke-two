package output

import (
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

//QuestionnaireInfo Questionnaireにtargetかの情報追加
type QuestionnaireInfo struct {
	Questionnaires model.Questionnaires
	targets        []string
	administrators []string
	respondents    []string
	pageMax        int
	IsTargeted     bool `json:"is_targeted" gorm:"type:boolean"`
}

//TargetedQuestionnaire targetになっているアンケートの情報
type TargetedQuestionnaire struct {
	Questionnaires model.Questionnaires
	RespondedAt    null.Time `json:"responded_at"`
	HasResponse    bool      `json:"has_response"`
}

type ResponseReadPrivilegeInfo struct {
	ResSharedTo     int
	IsAdministrator bool
	IsRespondent    bool
}
