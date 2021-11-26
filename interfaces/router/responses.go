package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
)

type ResponseAPI interface {
	PostResponse(c echo.Context) error
	GetResponse(c echo.Context) error
	EditResponse(c echo.Context) error
	DeleteResponse(c echo.Context) error
}

type response struct {
	usecase.ResponseUsecase
}

func NewResponseAPI(responseUsecase usecase.ResponseUsecase) ResponseAPI {
	return &response{ResponseUsecase: responseUsecase}
}

func (r *response) PostResponse(c echo.Context) error {
	panic("implement me")
}

func (r *response) GetResponse(c echo.Context) error {
	panic("implement me")
}

func (r *response) EditResponse(c echo.Context) error {
	panic("implement me")
}

func (r *response) DeleteResponse(c echo.Context) error {
	panic("implement me")
}
