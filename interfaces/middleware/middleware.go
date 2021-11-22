package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"net/http"
	"strconv"
)

type Middleware struct {
	repository.IAdministrator
	repository.IRespondent
	repository.IQuestion
	repository.IQuestionnaire
}

var adminUserIDs = []string{"temma", "sappi_red", "ryoha", "mazrean", "xxarupakaxx", "asari"}

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

func (m *Middleware) TrapReteLimitMiddlewareFunc() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		IdentifierExtractor: func(c echo.Context) (string, error) {
			userID, err := getUserID(c)
			if err != nil {
				c.Logger().Error(err)
				return "", echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
			}

			return userID, err
		},
		Store: middleware.NewRateLimiterMemoryStore(5),
	}

	return middleware.RateLimiterWithConfig(config)
}

func (m *Middleware) QuestionnaireAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Info()
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s (error :%w", strQuestionnaireID, err))
		}

		for _, adminUserID := range adminUserIDs {
			if userID == adminUserID {
				c.Set(questionnaireIDKey, questionnaireID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionnaireAdmin(c.Request().Context(), userID, questionnaireID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}
		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionnaireIDKey, questionnaireID)

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
