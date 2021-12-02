package model

type Validations struct {
	QuestionID   int    `json:"questionID"`
	RegexPattern string `json:"regex_pattern"`
	MinBound     string `json:"min_bound"`
	MaxBound     string `json:"max_bound"`
}
