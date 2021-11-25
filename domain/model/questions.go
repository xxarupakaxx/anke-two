package model

import (
	"gorm.io/gorm"
	"time"
)

type Questions struct {
	ID              int            `json:"id" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID" gorm:"type:int(11);not null"`
	PageNum         int            `json:"page_num" gorm:"type:int(11);not null"`
	QuestionNum     int            `json:"question_num" gorm:"type:int(11);not null"`
	Type            int            `json:"type" gorm:"type:int(11);not null"`
	Body            string         `json:"body" gorm:"type:text;default:NULL"`
	IsRequired      bool           `json:"is_required" gorm:"type:tinyint(4);size:4;not null;default:0"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"type:DATETIME NULL;default:NULL"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	Options         []Options      `json:"-" gorm:"foreignKey:QuestionID"`
	Responses       []Responses    `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []ScaleLabels  `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []Validations  `json:"-"  gorm:"foreignKey:QuestionID"`
}

type QuestionType struct {
	ID           int    `json:"id" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionType string `json:"question_type" gorm:"type:varchar(30);not null"`
	Active       bool   `json:"active" gorm:"type:boolean;not null;default:true"`
}

type ReturnQuestions struct {
	ID              int            `json:"id" `
	QuestionnaireID int            `json:"questionnaireID"`
	PageNum         int            `json:"page_num" `
	QuestionNum     int            `json:"question_num"`
	Type            string         `json:"type"`
	Body            string         `json:"body"`
	IsRequired      bool           `json:"is_required"`
	DeletedAt       gorm.DeletedAt `json:"-" `
	CreatedAt       time.Time      `json:"created_at"`
	Options         []Options      `json:"-"`
	Responses       []Responses    `json:"-"`
	ScaleLabels     []ScaleLabels  `json:"-"`
	Validations     []Validations  `json:"-"`
}

// BeforeCreate Update時に自動でmodified_atを現在時刻に
func (question *Questions) BeforeCreate(tx *gorm.DB) error {
	question.CreatedAt = time.Now()

	return nil
}

//TableName テーブル名が単数形なのでその対応
func (*Questions) TableName() string {
	return "question"
}

type QuestionIDAndQuestionType struct {
	QuestionID   int
	QuestionType string
	Responses    []Responses
}
