package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Respondents struct {
	ResponseID      int            `json:"responseID"`
	QuestionnaireID int            `json:"questionnaireID"`
	UserTraqid      string         `json:"user_traq_id,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// RespondentInfo 回答とその周辺情報の構造体
type RespondentInfo struct {
	Title        string    `json:"questionnaire_title"`
	ResTimeLimit null.Time `json:"res_time_limit"`
	Respondents
}

// RespondentDetail 回答の詳細情報の構造体
type RespondentDetail struct {
	ResponseID      int            `json:"responseID,omitempty"`
	TraqID          string         `json:"traqID,omitempty"`
	QuestionnaireID int            `json:"questionnaireID,omitempty"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty"`
	UpdatedAt      time.Time      `json:"modified_at,omitempty"`
	Responses       []ResponseBody `json:"body"`
}