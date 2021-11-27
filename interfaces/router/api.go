package router

import (
	"github.com/xxarupkaxx/anke-two/interfaces/middleware"
)

type API struct {
	*middleware.Mv
	*questionnaire
	*question
	*response
	*result
	*user
}

func NewAPI(mv *middleware.Mv, questionnaire *questionnaire, question *question, response *response, result *result, user *user) *API {
	return &API{Mv: mv, questionnaire: questionnaire, question: question, response: response, result: result, user: user}
}