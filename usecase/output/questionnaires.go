package output

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type PostAndEditQuestionnaireRequest struct {
	QuestionnaireID int            `json:"questionnaireID" validate:"required,min=0"`
	Title           string         `json:"title" validate:"required,max=50"`
	Description     string         `json:"description"`
	ResTimeLimit    null.Time      `json:"res_time_limit"`
	DeletedAt       gorm.DeletedAt `json:"-"`
	ResSharedTo     string         `json:"res_shared_to" validate:"required,oneof=administrators respondents public"`
	Targets         []string       `json:"targets" validate:"dive,max=32"`
	Administrators  []string       `json:"administrators" validate:"required,min=1,dive,max=32"`
}

type PostAndEditQuestionRequest struct {
	QuestionnaireID int      `json:"questionnaireID" validate:"min=0"`
	QuestionType    string   `json:"question_type" validate:"required,oneof=Text TextArea Number MultipleChoice Checkbox LinearScale"`
	QuestionNum     int      `json:"question_num" validate:"min=0"`
	PageNum         int      `json:"page_num" validate:"min=0"`
	Body            string   `json:"body" validate:"required"`
	IsRequired      bool     `json:"is_required"`
	Options         []string `json:"options" validate:"required_if=QuestionType Checkbox,required_if=QuestionType MultipleChoice,dive,max=50"`
	ScaleLabelRight string   `json:"scale_label_right" validate:"max=50"`
	ScaleLabelLeft  string   `json:"scale_label_left" validate:"max=50"`
	ScaleMin        int      `json:"scale_min"`
	ScaleMax        int      `json:"scale_max" validate:"gtecsfield=ScaleMin"`
	RegexPattern    string   `json:"regex_pattern"`
	MinBound        string   `json:"min_bound" validate:"omitempty,number"`
	MaxBound        string   `json:"max_bound" validate:"omitempty,number"`
}
