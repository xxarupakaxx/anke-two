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

type user struct {
	usecase.UsersUsecase
	middleware.IMiddleware
}

func NewUserAPI(usersUsecase usecase.UsersUsecase, IMiddleware middleware.IMiddleware) UserAPI {
	return &user{UsersUsecase: usersUsecase, IMiddleware: IMiddleware}
}

func (u *user) GetUsersMe(c echo.Context) error {
	userID, err := u.IMiddleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	in := input.GetMe{UserID: userID}

	out := u.UsersUsecase.GetUsersMe(c.Request().Context(), in)

	return c.JSON(http.StatusOK, out)
}

func (u *user) GetMyResponse(c echo.Context) error {
	userID, err := u.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}
	in := input.GetMe{
		UserID: userID,
	}

	out, err := u.UsersUsecase.GetMyResponses(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (u *user) GetMyResponsesByID(c echo.Context) error {
	userID, err := u.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	questionnaireID, err := strconv.Atoi(c.Param("questionnaireID"))
	if err != nil {
		c.Logger().Info(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	in := input.GetMyResponse{
		UserID:          userID,
		QuestionnaireID: questionnaireID,
	}

	out, err := u.UsersUsecase.GetMyResponsesByID(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (u *user) GetTargetedQuestionnaire(c echo.Context) error {
	userID, err := u.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}
	sort := c.QueryParam("sort")

	in := input.GetTargetedQuestionnaire{
		UserID: userID,
		Sort:   sort,
	}

	out, err := u.UsersUsecase.GetTargetedQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (u *user) GetMyQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (u *user) GetTargetedQuestionnairesByTraQID(c echo.Context) error {
	panic("implement me")
}

type UserAPI interface {
	GetUsersMe(c echo.Context) error
	GetMyResponse(c echo.Context) error
	GetMyResponsesByID(c echo.Context) error
	GetTargetedQuestionnaire(c echo.Context) error
	GetMyQuestionnaire(c echo.Context) error
	GetTargetedQuestionnairesByTraQID(c echo.Context) error
}
