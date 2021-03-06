package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/repository/middleware"
	"github.com/xxarupkaxx/anke-two/usecase"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"net/http"
	"strconv"
)

type User struct {
	usecase.UsersUsecase
	middleware.IMiddleware
}

func NewUserAPI(usersUsecase usecase.UsersUsecase, IMiddleware middleware.IMiddleware) *User {
	return &User{UsersUsecase: usersUsecase, IMiddleware: IMiddleware}
}

func (u *User) GetUsersMe(c echo.Context) error {
	userID, err := u.IMiddleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	in := input.GetMe{UserID: userID}

	out := u.UsersUsecase.GetUsersMe(c.Request().Context(), in)

	return c.JSON(http.StatusOK, out)
}

func (u *User) GetMyResponse(c echo.Context) error {
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

func (u *User) GetMyResponsesByID(c echo.Context) error {
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

func (u *User) GetTargetedQuestionnaire(c echo.Context) error {
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

func (u *User) GetMyQuestionnaire(c echo.Context) error {
	userID, err := u.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
	}

	in := input.GetMe{
		UserID: userID,
	}

	out, err := u.UsersUsecase.GetMyQuestionnaire(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

func (u *User) GetTargetedQuestionnairesByTraQID(c echo.Context) error {
	traQID := c.Param("traQID")
	sort := c.QueryParam("sort")
	answered := c.QueryParam("answered")

	in := input.GetTargetsByTraQID{
		TraQID:   traQID,
		Sort:     sort,
		Answered: answered,
	}

	out, err := u.UsersUsecase.GetTargetedQuestionnairesByID(c.Request().Context(), in)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, out)
}

type UserAPI interface {
	GetUsersMe(c echo.Context) error
	GetMyResponse(c echo.Context) error
	GetMyResponsesByID(c echo.Context) error
	GetTargetedQuestionnaire(c echo.Context) error
	GetMyQuestionnaire(c echo.Context) error
	GetTargetedQuestionnairesByTraQID(c echo.Context) error
}
