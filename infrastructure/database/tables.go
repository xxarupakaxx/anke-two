package database

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

var (
	allTables = []interface{}{
		Questionnaires{},
		Questions{},
		Respondents{},
		Responses{},
		Administrators{},
		Options{},
		ScaleLabels{},
		ResSharedTo{},
		QuestionType{},
		Targets{},
		Validations{},
	}
)

type Administrators struct {
	QuestionnaireID int    `gorm:"type:int(11);not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}


type Options struct {
	ID         int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionID int    `gorm:"type:int(11);not null"`
	OptionNum  int    `gorm:"type:int(11);not null"`
	Body       string `gorm:"type:text;default:NULL;"`
}

type Questionnaires struct {
	ID             int              `json:"questionnaireID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Title          string           `json:"title" gorm:"type:char(50);size:50;not null" `
	Description    string           `json:"description" gorm:"type:text;not null"`
	ResTimeLimit   null.Time        `json:"res_time_limit,omitempty" gorm:"type:DATETIME NULL;default:NULL;"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"type:DATETIME NULL;default:NULL;"`
	ResSharedTo    int              `json:"res_shared_to" gorm:"type:int(11);not null;default:0"`
	CreatedAt      time.Time        `json:"created_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	ModifiedAt     time.Time        `json:"modified_at" gorm:"DATETIME;not null;default:CURRENT_TIMESTAMP"`
	Administrators []Administrators `json:"-"  gorm:"foreignKey:QuestionnaireID"`
	Targets        []Targets        `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Questions      []Questions      `json:"-" gorm:"foreignKey:QuestionnaireID"`
	Respondents    []Respondents    `json:"-" gorm:"foreignKey:QuestionnaireID"`
}

type ResSharedTo struct {
	ID     int    `json:"id" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Name   string `json:"name" gorm:"type:varchar(30);not null"`
	Active bool   `json:"active" gorm:"type:boolean;not null;default:true"`
}

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
	ID     int    `json:"id" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	Name   string `json:"question_type" gorm:"type:varchar(30);not null"`
	Active bool   `json:"active" gorm:"type:boolean;not null;default:true"`
}

type Respondents struct {
	ResponseID      int            `json:"responseID" gorm:"column:response_id;type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	QuestionnaireID int            `json:"questionnaireID" gorm:"type:int(11);not null"`
	UserTraqid      string         `json:"user_traq_id,omitempty" gorm:"type:varchar(32);size:32;default:NULL"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	SubmittedAt     null.Time      `json:"submitted_at,omitempty" gorm:"type:DATETIME NULL;default:NULL"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:DATETIME NULL;default:NULL"`
	Responses       []Responses    `json:"-" gorm:"foreignKey:ResponseID;references:ResponseID"`
}

type Responses struct {
	ResponseID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	QuestionID int            `json:"-" gorm:"type:int(11);not null;primaryKey"`
	Body       null.String    `json:"response" gorm:"type:text;default:NULL"`
	UpdatedAt  time.Time      `json:"-" gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"type:DATETIME NULL;default:NULL"`
}

type ScaleLabels struct {
	QuestionID      int    `json:"questionID" gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	ScaleLabelRight string `json:"scale_label_right" gorm:"type:varchar(50);default:NULL;"`
	ScaleLabelLeft  string `json:"scale_label_left" gorm:"type:varchar(50);default:NULL;"`
	ScaleMin        int    `json:"scale_min" gorm:"type:int(11);default:NULL;"`
	ScaleMax        int    `json:"scale_max" gorm:"type:int(11);default:NULL;"`
}

type Targets struct {
	QuestionnaireID int    `gorm:"type:int(11) AUTO_INCREMENT;not null;primaryKey"`
	UserTraqid      string `gorm:"type:varchar(32);size:32;not null;primaryKey"`
}

type Validations struct {
	QuestionID   int    `json:"questionID"    gorm:"type:int(11);not null;primaryKey"`
	RegexPattern string `json:"regex_pattern" gorm:"type:text;default:NULL"`
	MinBound     string `json:"min_bound"     gorm:"type:text;default:NULL"`
	MaxBound     string `json:"max_bound"     gorm:"type:text;default:NULL"`
}
