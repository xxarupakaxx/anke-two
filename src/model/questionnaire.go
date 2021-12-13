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

	Administrators  []Administrator `gorm:"foreignKey:QuestionnaireID"`
	Targets         []Target        `gorm:"foreignKey:QuestionnaireID"`
	Questions       []Question      `gorm:"foreignKey:QuestionnaireID"`
	Respondents     []Respondent    `gorm:"foreignKey:QuestionnaireID"`
	ResSharedToName ResSharedTo     `gorm:"foreignKey:ID;references:ResSharedTo"`
}

func (q *Questionnaire) TableName() string {
	return "questionnaire"
}

// ResSharedTo アンケート結果の公開範囲の種類。 アンケートの結果を、運営は見られる ("administrators")、回答済みの人は見られる ("respondents")、誰でも見られる ("public")。のテーブル
type ResSharedTo struct {
	ID     int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);not null"`
	Active bool   `gorm:"type:boolean;not null;default:true"`
}

func (r *ResSharedTo) TableName() string {
	return "res_shared_to"
}