package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Responses struct {
	ResponseID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	QuestionID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	Body       null.String    `json:"response" gorm:"type:text;default:NULL"`
	UpdatedAt  time.Time      `json:"-" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"type:DATETIME NULL;default:NULL"`
}
