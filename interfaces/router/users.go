package router

import (
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/usecase"
)

type user struct {
	usecase.UsersUsecase
}

func NewUserAPI(usersUsecase usecase.UsersUsecase) UserAPI {
	return &user{UsersUsecase: usersUsecase}
}

func (u *user) GetUsesMe(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetMyResponse(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetMyResponsesByID(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetTargetedQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetMyQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetTargetedQuestionnairesByTraQID(c echo.Context) error {
	panic("implement me")
}

type UserAPI interface {
	GetUsesMe(c echo.Context) error
	GetMyResponse(c echo.Context) error
	GetMyResponsesByID(c echo.Context) error
	GetTargetedQuestionnaire(c echo.Context) error
	GetMyQuestionnaire(c echo.Context) error
	GetTargetedQuestionnairesByTraQID(c echo.Context) error
}
