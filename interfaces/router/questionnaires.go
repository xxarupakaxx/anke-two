package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/repository/middleware"
	"github.com/xxarupkaxx/anke-two/usecase"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"net/http"
	"strconv"
)

type questionnaire struct {
	usecase.QuestionnaireUsecase
	middleware.IMiddleware
}

func NewQuestionnaireAPI(questionnaireUsecase usecase.QuestionnaireUsecase, middleware middleware.IMiddleware) QuestionnaireAPI {
	return &questionnaire{QuestionnaireUsecase: questionnaireUsecase, IMiddleware: middleware}
}

func (q *questionnaire) GetQuestionnaires(c echo.Context) error {
	userID, err := q.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	sort := c.QueryParam("sort")
	search := c.QueryParam("search")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	nontargeted, err := strconv.ParseBool(c.QueryParam("nontargeted"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.GetQuestionnairesQueryParam{
		UserID:      userID,
		Sort:        sort,
		Search:      search,
		Page:        page,
		Nontargeted: nontargeted,
	}

	out, err := q.QuestionnaireUsecase.GetQuestionnaires(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (q *questionnaire) PostQuestionnaire(c echo.Context) error {
	in := input.PostAndEditQuestionnaireRequest{}

	if err := c.Bind(&in); err != nil {
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

	out, err := q.QuestionnaireUsecase.PostQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (q *questionnaire) GetQuestionnaire(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.GetQuestionnaire{QuestionnaireID: questionnaireID}

	out, err := q.QuestionnaireUsecase.GetQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (q *questionnaire) PostQuestionByQuestionnaireID(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	in := input.PostQuestionRequest{}
	in.QuestionnaireID = questionnaireID

	if err = c.Bind(&in); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
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

	err = q.QuestionnaireUsecase.ValidationPostQuestionByQuestionnaireID(in)
	if err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	out, err := q.QuestionnaireUsecase.PostQuestionByQuestionnaireID(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
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
