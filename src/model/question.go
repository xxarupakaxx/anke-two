package model

import (
	"gorm.io/gorm"
	"time"
)

// Question 質問テーブル
type Question struct {
	ID              int            `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `gorm:"type:int(11);not null"`
	QuestionNum     int            `gorm:"type:int(11);not null"`
	Type            int            `gorm:"type:int(11);not null"`
	Body            string         `gorm:"type:text;default:NULL"`
	IsRequired      bool           `gorm:"type:boolean;not null;default:false"`
	DeletedAt       gorm.DeletedAt `gorm:"precision:6"`
	CreatedAt       time.Time      `gorm:"precision:6"`
	UpdatedAt       time.Time      `gorm:"precision:6"`

	Options         []Options     `json:"-" gorm:"foreignKey:QuestionID"`
	Responses       []Response   `json:"-"  gorm:"foreignKey:QuestionID"`
	ScaleLabels     []ScaleLabel `json:"-"  gorm:"foreignKey:QuestionID"`
	Validations     []Validation `json:"-"  gorm:"foreignKey:QuestionID"`
}
