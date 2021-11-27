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

type Questionnaire struct {
	usecase.QuestionnaireUsecase
	middleware.IMiddleware
}

func NewQuestionnaireAPI(questionnaireUsecase usecase.QuestionnaireUsecase, middleware middleware.IMiddleware) *Questionnaire {
	return &Questionnaire{QuestionnaireUsecase: questionnaireUsecase, IMiddleware: middleware}
}

func (q *Questionnaire) GetQuestionnaires(c echo.Context) error {
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

func (q *Questionnaire) PostQuestionnaire(c echo.Context) error {
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

func (q *Questionnaire) GetQuestionnaire(c echo.Context) error {
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

func (q *Questionnaire) PostQuestionByQuestionnaireID(c echo.Context) error {
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

func (q *Questionnaire) EditQuestionnaire(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	in := input.PostAndEditQuestionnaireRequest{}
	in.QuestionnaireID = questionnaireID

	if err = c.Bind(&in); err != nil {
		c.Logger().Info(err)
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

	err = q.QuestionnaireUsecase.EditQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (q *Questionnaire) DeleteQuestionnaire(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.DeleteQuestionnaire{}
	in.QuestionnaireID = questionnaireID

	err = q.QuestionnaireUsecase.DeleteQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (q *Questionnaire) GetQuestions(c echo.Context) error {
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", questionnaireID, err))
	}

	in := input.QuestionInfo{}
	in.QuestionnaireID = questionnaireID

	out, err := q.QuestionnaireUsecase.GetQuestions(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
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
