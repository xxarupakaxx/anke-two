package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"net/http"
)

type Middleware struct {
	repository.IAdministrator
	repository.IRespondent
	repository.IQuestion
	repository.IQuestionnaire
}

func NewMiddleware(IAdministrator repository.IAdministrator, IRespondent repository.IRespondent, IQuestion repository.IQuestion, IQuestionnaire repository.IQuestionnaire) *Middleware {
	return &Middleware{IAdministrator: IAdministrator, IRespondent: IRespondent, IQuestion: IQuestion, IQuestionnaire: IQuestionnaire}
}

const (
	validatorKey       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

func (m *Middleware) SetValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		c.Set(validatorKey, validate)

		return next(c)
	}
}

func (m *Middleware) SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Showcase-User")
		if userID == "" {
			userID = "xxarupakaxx"
		}

		c.Set(userIDKey, userID)

		return next(c)
	}
}

func (m *Middleware) TraPMemberAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID:%w", err))
		}

		if userID == "-" {
			return echo.NewHTTPError(http.StatusUnauthorized, "you are not log in")
		}

		return next(c)
	}
}

func getUserID(c echo.Context) (string, error) {
	rowUserID := c.Get(userIDKey)
	userID, ok := rowUserID.(string)
	if !ok {
		return "", errors.New("invalid context userID")
	}

	return userID, nil
}
