package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Questionnaires struct {
	ID             int              `json:"questionnaireID" gorm:"type:int(1)) AUTO_INCREMENT;not null;primaryKey"`
	Title          string           `json:"title" gorm:"type:char(50);size:50;not null" `
	Description    string           `json:"description" gorm:"type:text;not null"`
	ResTimeLimit   null.Int         `json:"res_time_limit,omitempty" gorm:"type:DATETIME NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"type:DATETIME NULL;default:NULL;"`
	ResSharedTo    string           `json:"res_shared_to" gorm:"type:char(30);size:30;not null;default:administrators"`
	CreatedAt      time.Time        `json:"created_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time        `json:"modified_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	Administrators []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets []Targets `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Questions []Questions `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Respondents []Respondents `json:"-" gorm:"foreignKey:QuestionnaireID"`
}
