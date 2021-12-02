package model

import (
	"gorm.io/gorm"
	"time"
)

type Questions struct {
	ID              int            `json:"id"`
	QuestionnaireID int            `json:"questionnaireID"`
	PageNum         int            `json:"page_num"`
	QuestionNum     int            `json:"question_num"`
	Type            int            `json:"type"`
	Body            string         `json:"body"`
	IsRequired      bool           `json:"is_required"`
	DeletedAt       gorm.DeletedAt `json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	Options         []Options      `json:"-"`
	Responses       []Responses    `json:"-"`
	ScaleLabels     []ScaleLabels  `json:"-"`
	Validations     []Validations  `json:"-"`
}

type QuestionType struct {
	ID     int    `json:"id"`
	Name   string `json:"question_type"`
	Active bool   `json:"active"`
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
