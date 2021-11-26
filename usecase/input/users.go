package input

import "gopkg.in/guregu/null.v4"

type GetMe struct {
	UserID string `validate:"required"`
}

type GetMyResponse struct {
	UserID          string `validate:"required"`
	QuestionnaireID int    `validate:"required,min=0"`
}

type GetTargetedQuestionnaire struct {
	UserID string
	Sort   string `validate:"omitempty,oneof=created_at -created_at title -title modified_at -modified_at"`
}
type UserQueryParam struct {
	Sort     string `validate:"omitempty,oneof=created_at -created_at title -title modified_at -modified_at"`
	Answered string `validate:"omitempty,oneof=answered unanswered"`
}

type QuestionnaireInfo struct {
	ID             int       `json:"questionnaireID" validate:"required,min=0"`
	Title          string    `json:"title" validate:"required,max=50"`
	Description    string    `json:"description"`
	ResTimeLimit   null.Time `json:"res_time_limit"`
	CreatedAt      string    `json:"created_at"`
	ModifiedAt     string    `json:"modified_at"`
	ResSharedTo    string    `json:"res_shared_to" validate:"required,oneof=administrators respondents public"`
	AllResponded   bool      `json:"all_responded"`
	Targets        []string  `json:"targets" validate:"dive,max=32"`
	Administrators []string  `json:"administrators" validate:"required,min=1,dive,max=32"`
	Respondents    []string  `json:"respondents"`
}

type GetTargetsByTraQID struct {
	TraQID   string
	Sort     string `validate:"omitempty,oneof=created_at -created_at title -title modified_at -modified_at"`
	Answered string `validate:"omitempty,oneof=answered unanswered"`
}
