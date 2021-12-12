package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

// Respondent アンケートごとの回答者
type Respondent struct {
	ResponseID      int            `gorm:"column:response_id;type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `gorm:"type:int(11);not null"`
	UserTraqid      string         `gorm:"type:varchar(32);size:32;default:NULL"`
	UpdatedAt       time.Time      `gorm:"precision:6"`
	SubmittedAt     null.Time      `gorm:"precision:6"`
	DeletedAt       gorm.DeletedAt `gorm:"precision:6"`

	Responses       []Response     `gorm:"foreignKey:ResponseID;references:ResponseID"`
}
