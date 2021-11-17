package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Respondents struct {
	ResponseID      int            `json:"responseID" gorm:"column:response_id;type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID" gorm:"type:int(11);not null"`
	UserTraqid      string         `json:"user_traq_id,omitempty" gorm:"type:varchar(32);size:32;default:NULL"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty" gorm:"type:DATETIME NULL;default:NULL"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:DATETIME NULL;default:NULL"`
	Responses       []Responses    `json:"-" gorm:"foreignKey:ResponseID;references:ResponseID"`
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