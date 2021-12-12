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

	QuestionType QuestionType       `gorm:"foreignKey:ID;references:Type"`
	Options     []Options    		`gorm:"foreignKey:QuestionID"`
	Responses   []Response          `gorm:"foreignKey:QuestionID"`
	ScaleLabels []ScaleLabel        `gorm:"foreignKey:QuestionID"`
	Validations []Validation        `gorm:"foreignKey:QuestionID"`
}

// QuestionType 質問の種類。 'Text'、'TextArea'、'Number'、'MultipleChoice'、'Checkbox', 'Dropdown', 'LinearScale', 'Date', 'Time' のテーブル
type QuestionType struct {
	ID     int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Name   string `gorm:"type:varchar(32);not null"`
	Active bool   `gorm:"type:boolean;not null;default:true"`
}
