package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"net/http"
	"strconv"
)

type result struct {
	usecase.ResultUsecase
}

func NewResultAPI(resultUsecase usecase.ResultUsecase) ResultAPI {
	return &result{ResultUsecase: resultUsecase}
}

func (r *result) GetResults(c echo.Context) error {
	sort := c.QueryParam("sort")
	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.GetResults{
		Sort:            sort,
		QuestionnaireID: questionnaireID,
	}
	out, err := r.ResultUsecase.GetResults(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

type ResultAPI interface {
	GetResults(c echo.Context) error
}
