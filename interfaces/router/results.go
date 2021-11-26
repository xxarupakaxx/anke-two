package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
)

type result struct {
	usecase.ResultUsecase
}

func NewResultAPI(resultUsecase usecase.ResultUsecase) ResultAPI {
	return &result{ResultUsecase: resultUsecase}
}

func (r *result) GetResults(c echo.Context) error {
	panic("implement me")
}

type ResultAPI interface {
	GetResults(c echo.Context) error
}
