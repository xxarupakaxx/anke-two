package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
)

type question struct {
	usecase.QuestionUsecase
}

func NewQuestionAPI(questionUsecase usecase.QuestionUsecase) QuestionAPI {
	return &question{QuestionUsecase: questionUsecase}
}

func (q *question) EditQuestion(c echo.Context) error {
	panic("implement me")
}

func (q *question) DeleteQuestion(c echo.Context) error {
	panic("implement me")
}

type QuestionAPI interface {
	EditQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
}
