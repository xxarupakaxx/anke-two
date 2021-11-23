package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xxarupkaxx/anke-two/domain/model"
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

func (m *Middleware) ResponseReadAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		respondent, err := m.GetRespondent(c.Request().Context(), responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}
		if respondent == nil {
			c.Logger().Error("respondent is nil")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if respondent.UserTraqid == userID {
			return next(c)
		}

		if !respondent.SubmittedAt.Valid {
			c.Logger().Info("not submitted")

			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}

		responseReadPrivilegeInfo, err := m.GetResponseReadPrivilegeInfoByResponseID(c.Request().Context(), userID, responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", responseID))
		} else if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilegeInfo)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check response read privilege: %w", err))
		}

		if !haveReadPrivilege {
			return c.String(http.StatusForbidden, "You do not have permission to view this response.")
		}

		return next(c)

	}
}

func (m *Middleware) RespondentsAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strResponseID := c.Param("responseID")
		responseID, err := strconv.Atoi(strResponseID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid responseID:%s(error: %w)", strResponseID, err))
		}

		respondent, err := m.GetRespondent(c.Request().Context(), responseID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("response not found:%d", responseID))
		}
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are a respondent: %w", err))
		}
		if respondent == nil {
			c.Logger().Error("respondent is nil")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if respondent.UserTraqid != userID {
			return c.String(http.StatusForbidden, "You are not a respondent of this response.")
		}

		c.Set(responseIDKey, responseID)

		return next(c)
	}
}

func (m *Middleware) QuestionAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionID := c.Param("questionID")
		questionID, err := strconv.Atoi(strQuestionID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionID:%s(error: %w)", strQuestionID, err))
		}

		for _, adminUserID := range adminUserIDs {
			if userID == adminUserID {
				c.Set(questionIDKey, questionID)

				return next(c)
			}
		}
		isAdmin, err := m.CheckQuestionAdmin(c.Request().Context(), userID, questionID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check if you are administrator: %w", err))
		}

		if !isAdmin {
			return c.String(http.StatusForbidden, "You are not a administrator of this questionnaire.")
		}

		c.Set(questionIDKey, questionID)

		return next(c)
	}
}

func (m *Middleware) ResultAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := getUserID(c)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get userID: %w", err))
		}

		strQuestionnaireID := c.Param("questionnaireID")
		questionnaireID, err := strconv.Atoi(strQuestionnaireID)
		if err != nil {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid questionnaireID:%s(error: %w)", strQuestionnaireID, err))
		}

		responseReadPrivilege, err := m.GetResponseReadPrivilegeInfoByQuestionnaireID(c.Request().Context(), userID, questionnaireID)
		if errors.Is(err, model.ErrRecordNotFound) {
			c.Logger().Info(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid responseID: %d", questionnaireID))
		} else if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to get response read privilege info: %w", err))
		}

		haveReadPrivilege, err := checkResponseReadPrivilege(responseReadPrivilege)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to check response read privilege: %w", err))
		}
		if !haveReadPrivilege {
			return c.String(http.StatusForbidden, "You do not have permission to view this response.")
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

func GetValidator(c echo.Context) (*validator.Validate, error) {
	rowValidate := c.Get(validatorKey)
	validate, ok := rowValidate.(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator")
	}

	return validate, nil
}

func GetQuestionnaireID(c echo.Context) (int, error) {
	rowQuestionnaireID := c.Get(questionnaireIDKey)
	questionnaireID, ok := rowQuestionnaireID.(int)
	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return questionnaireID, nil
}

func GetResponseID(c echo.Context) (int, error) {
	rowResponseID := c.Get(responseIDKey)
	responseID, ok := rowResponseID.(int)

	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return responseID, nil
}

func GetQuestionID(c echo.Context) (int, error) {
	rowQuestionID := c.Get(questionIDKey)
	questionID, ok := rowQuestionID.(int)
	if !ok {
		return 0, errors.New("invalid context userID")
	}

	return questionID, nil
}

func checkResponseReadPrivilege(responseReadPrivilegeInfo *model.ResponseReadPrivilegeInfo) (bool, error) {
	switch responseReadPrivilegeInfo.ResSharedTo {
	case "administrators":
		return responseReadPrivilegeInfo.IsAdministrator, nil
	case "respondents":
		return responseReadPrivilegeInfo.IsAdministrator || responseReadPrivilegeInfo.IsRespondent, nil
	case "public":
		return true, nil
	}

	return false, errors.New("invalid resSharedTo")
}
