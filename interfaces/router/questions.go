package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"net/http"
	"strconv"
)

type Question struct {
	usecase.QuestionUsecase
}

func NewQuestionAPI(questionUsecase usecase.QuestionUsecase) *Question {
	return &Question{QuestionUsecase: questionUsecase}
}

func (q *Question) EditQuestion(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.EditQuestionRequest{}
	in.QuestionID = questionID

	if err := c.Bind(&in); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	statusCode, err := usecase.ValidateRequest(c, in)
	if err != nil {
		switch statusCode {
		case http.StatusBadRequest:
			c.Logger().Info(err)
			return echo.NewHTTPError(statusCode)
		case http.StatusInternalServerError:
			c.Logger().Error(err)
			return echo.NewHTTPError(statusCode)
		}
	}

	err = q.QuestionUsecase.EditQuestion(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (q *Question) DeleteQuestion(c echo.Context) error {
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.DeleteQuestion{}
	in.QuestionID = questionID

	err = q.QuestionUsecase.DeleteQuestion(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

type QuestionAPI interface {
	EditQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
}
