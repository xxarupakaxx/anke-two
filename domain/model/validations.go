package model

type Validations struct {
	QuestionID   int    `json:"questionID"    gorm:"type:int(11);not null;primaryKey"`
	RegexPattern string `json:"regex_pattern" gorm:"type:text;default:NULL"`
	MinBound     string `json:"min_bound"     gorm:"type:text;default:NULL"`
	MaxBound     string `json:"max_bound"     gorm:"type:text;default:NULL"`
}
