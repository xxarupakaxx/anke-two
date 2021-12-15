package repository

type Repository interface {
	Transaction
	Question
	Admin
	Option
	ScaleLabel
	Target
	Response
	Respondent
	Questionnaire
	Validation
}
