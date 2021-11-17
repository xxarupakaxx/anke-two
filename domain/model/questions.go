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
	Type            string         `json:"type" gorm:"type:char(20);size:20;not null"`
	Body            string         `json:"body" gorm:"type:text;default:NULL"`
	IsRequired      bool           `json:"is_required" gorm:"type:tinyint(4);size:4;not null;default:0"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"type:DATETIME NULL;default:NULL"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	Options         []Options      `json:"-" gorm:"foreignKey:QuestionID"`
	Responses       []Responses    `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []ScaleLabels  `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []Validations  `json:"-"  gorm:"foreignKey:QuestionID"`
}
