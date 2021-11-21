package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

// Responses 解答情報
type Responses struct {
	ResponseID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	QuestionID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	Body       null.String    `json:"response" gorm:"type:text;default:NULL"`
	UpdatedAt  time.Time      `json:"-" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"type:DATETIME NULL;default:NULL"`
}

// ResponseBody 質問に対する回答の構造体
type ResponseBody struct {
	QuestionID     int         `json:"questionID" gorm:"column:id" validate:"min=0"`
	QuestionType   string      `json:"question_type" gorm:"column:type" validate:"required,oneof=Text TextArea Number MultipleChoice Checkbox LinearScale"`
	Body           null.String `json:"response" validate:"required"`
	OptionResponse []string    `json:"option_response" validate:"required_if=QuestionType Checkbox,required_if=QuestionType MultipleChoice,dive,max=50"`
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
