package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
)

type questionnaire struct {
	usecase.QuestionnaireUsecase
}

func NewQuestionnaireAPI(questionnaireUsecase usecase.QuestionnaireUsecase) QuestionnaireAPI {
	return &questionnaire{QuestionnaireUsecase: questionnaireUsecase}
}

func (q *questionnaire) GetQuestionnaires(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) PostQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) GetQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) PostQuestionByQuestionnaireID(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) EditQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) DeleteQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) GetQuestions(c echo.Context) error {
	panic("implement me")
}

type QuestionnaireAPI interface {
	GetQuestionnaires(c echo.Context) error
	PostQuestionnaire(c echo.Context) error
	GetQuestionnaire(c echo.Context) error
	PostQuestionByQuestionnaireID(c echo.Context) error
	EditQuestionnaire(c echo.Context) error
	DeleteQuestionnaire(c echo.Context) error
	GetQuestions(c echo.Context) error
}
