package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

// Questionnaire アンケートの情報の情報
type Questionnaire struct {
	ID           int            `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title        string         `gorm:"type:char(50);size:50;not null" `
	Description  string         `gorm:"type:text;not null"`
	ResTimeLimit null.Time      `gorm:"type:DATETIME NULL;default:NULL;"`
	ResSharedTo  int            `gorm:"type:int(11);not null;default:0"`
	CreatedAt    time.Time      `gorm:"precision:6"`
	UpdatedAt    time.Time      `gorm:"precision:6"`
	DeletedAt    gorm.DeletedAt `gorm:"precision:6"`

	Administrators []Administrator `gorm:"foreignKey:QuestionnaireID"`
	Targets        []Target        `gorm:"foreignKey:QuestionnaireID"`
	Questions      []Question      `gorm:"foreignKey:QuestionnaireID"`
	Respondents    []Respondent    `gorm:"foreignKey:QuestionnaireID"`
}
