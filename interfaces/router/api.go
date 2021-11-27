package router

import (
	"github.com/xxarupkaxx/anke-two/interfaces/middleware"
)

type API struct {
	*middleware.Mv
	*Questionnaire
	*Question
	*Response
	*Result
	*User
}

func NewAPI(mv *middleware.Mv, questionnaire *Questionnaire, question *Question, response *Response, result *Result, user *User) *API {
	return &API{Mv: mv, Questionnaire: questionnaire, Question: question, Response: response, Result: result, User: user}
}
