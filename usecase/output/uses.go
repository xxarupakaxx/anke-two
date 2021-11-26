package output

import "gopkg.in/guregu/null.v4"

type GetMe struct {
	TraqID string `json:"traqID" validate:"required"`
}

type QuestionnaireInfo struct {
	ID             int       `json:"questionnaireID"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ResTimeLimit   null.Time `json:"res_time_limit"`
	CreatedAt      string    `json:"created_at"`
	ModifiedAt     string    `json:"modified_at"`
	ResSharedTo    string    `json:"res_shared_to"`
	AllResponded   bool      `json:"all_responded"`
	Targets        []string  `json:"targets"`
	Administrators []string  `json:"administrators"`
	Respondents    []string  `json:"respondents"`
}