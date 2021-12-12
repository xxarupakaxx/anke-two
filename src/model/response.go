package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

// Response 回答者のテーブル
type Response struct {
	ResponseID int            `gorm:"type:int(11);not null;primaryKey"`
	QuestionID int            `gorm:"type:int(11);not null;primaryKey"`
	Body       null.String    `gorm:"type:text;default:NULL"`
	UpdatedAt  time.Time      `gorm:"precision:6"`
	DeletedAt  gorm.DeletedAt `gorm:"precision:6"`
}