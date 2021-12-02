package model

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"time"
)

type Questionnaires struct {
	ID             int              `json:"questionnaireID" `
	Title          string           `json:"title" `
	Description    string           `json:"description" `
	ResTimeLimit   null.Time        `json:"res_time_limit,omitempty" `
	DeletedAt      gorm.DeletedAt   `json:"-" `
	ResSharedTo    int              `json:"res_shared_to" `
	CreatedAt      time.Time        `json:"created_at" `
	ModifiedAt     time.Time        `json:"modified_at" `
	Administrators []Administrators `json:"-"  `
	Targets        []Targets        `json:"-"`
	Questions      []Questions      `json:"-" `
	Respondents    []Respondents    `json:"-"`
}

type ResSharedTo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

//QuestionnaireInfo Questionnaireにtargetかの情報追加
type QuestionnaireInfo struct {
	Questionnaires Questionnaires
	Targets        []string
	Administrators []string
	Respondents    []string
	PageMax        int
}

//TargetedQuestionnaire targetになっているアンケートの情報
type TargetedQuestionnaire struct {
	Questionnaires Questionnaires
	RespondedAt    null.Time `json:"responded_at"`
	HasResponse    bool      `json:"has_response"`
}

type ResponseReadPrivilegeInfo struct {
	ResSharedTo     string
	IsAdministrator bool
	IsRespondent    bool
}

type ReturnQuestionnaires struct {
	ID             int              `json:"questionnaireID"`
	Title          string           `json:"title" `
	Description    string           `json:"description"`
	ResTimeLimit   null.Time        `json:"res_time_limit,omitempty"`
	DeletedAt      gorm.DeletedAt   `json:"-"`
	ResSharedTo    string           `json:"res_shared_to"`
	CreatedAt      time.Time        `json:"created_at" `
	ModifiedAt     time.Time        `json:"modified_at"`
	Administrators []Administrators `json:"-" `
	Targets        []Targets        `json:"-" `
	Questions      []Questions      `json:"-"`
	Respondents    []Respondents    `json:"-"`
}

//BeforeCreate create時に自動でmodified_atを現在時刻に
func (questionnaire *Questionnaires) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	questionnaire.ModifiedAt = now
	questionnaire.CreatedAt = now

	return nil
}

//BeforeUpdate Update時に自動でmodified_atを現在時刻に
func (questionnaire *Questionnaires) BeforeUpdate(tx *gorm.DB) error {
	questionnaire.ModifiedAt = time.Now()

	return nil
}
