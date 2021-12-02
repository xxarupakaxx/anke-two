package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

// Responses 解答情報
type Responses struct {
	ResponseID int            `json:"-"`
	QuestionID int            `json:"-"`
	Body       null.String    `json:"response"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

// ResponseBody 質問に対する回答の構造体
type ResponseBody struct {
	QuestionID     int         `json:"questionID"`
	QuestionType   string      `json:"question_type"`
	Body           null.String `json:"response"`
	OptionResponse []string    `json:"option_response"`
}

// ResponseMeta 質問に対する回答の構造体
type ResponseMeta struct {
	QuestionID int
	Data       string
}

//BeforeCreate insert時に自動でmodifiedAt更新
func (r *Responses) BeforeCreate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (r *Responses) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()

	return nil
}

//TableName テーブル名が単数形なのでその対応
func (r *Responses) TableName() string {
	return "response"
}
