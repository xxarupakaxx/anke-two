package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Questionnaires struct {
	ID             int              `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title          string           `json:"title" gorm:"type:char(50);size:50;not null" `
	Description    string           `json:"description" gorm:"type:text;not null"`
	ResTimeLimit   null.Int         `json:"res_time_limit,omitempty" gorm:"type:DATETIME NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"type:DATETIME NULL;default:NULL;"`
	ResSharedTo    int              `json:"res_shared_to" gorm:"type:int(11);not null;default:administrators"`
	CreatedAt      time.Time        `json:"created_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time        `json:"modified_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	Administrators []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets        []Targets        `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Questions      []Questions      `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Respondents    []Respondents    `json:"-" gorm:"foreignKey:QuestionnaireID"`
}

type ResShareTypes struct {
	ID     int    `json:"id" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Name   string `json:"name" gorm:"type:varchar(30);not null"`
	Active bool   `json:"active" gorm:"type:boolean;not null;default:true"`
}
