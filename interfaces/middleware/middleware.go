package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	myMiddleware "github.com/xxarupkaxx/anke-two/domain/repository/middleware"
	"net/http"
	"strconv"
)

type mw struct {
	repository.IAdministrator
	repository.IRespondent
	repository.IQuestion
	repository.IQuestionnaire
}

func NewMiddleware(IAdministrator repository.IAdministrator, IRespondent repository.IRespondent, IQuestion repository.IQuestion, IQuestionnaire repository.IQuestionnaire) myMiddleware.IMiddleware {
	return &mw{IAdministrator: IAdministrator, IRespondent: IRespondent, IQuestion: IQuestion, IQuestionnaire: IQuestionnaire}
}

var adminUserIDs = []string{"temma", "sappi_red", "ryoha", "mazrean", "xxarupakaxx", "asari"}


const (
	validatorKey       = "validator"
	userIDKey          = "userID"
	questionnaireIDKey = "questionnaireID"
	responseIDKey      = "responseID"
	questionIDKey      = "questionID"
)

func (m *mw) SetValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		c.Set(validatorKey, validate)

		return next(c)
	}
}

func (m *mw) SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Showcase-User")
		if userID == "" {
			userID = "xxarupakaxx"
		}

		c.Set(userIDKey, userID)

		return next(c)
	}
}

func (m *mw) TraPMemberAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) TrapReteLimitMiddlewareFunc() echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		IdentifierExtractor: func(c echo.Context) (string, error) {
			userID, err := m.GetUserID(c)
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

func (m *mw) QuestionnaireAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) ResponseReadAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) RespondentsAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) QuestionAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) ResultAuthenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := m.GetUserID(c)
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

func (m *mw) GetUserID(c echo.Context) (string, error) {
	rowUserID := c.Get(userIDKey)
	userID, ok := rowUserID.(string)
	if !ok {
		return "", errors.New("invalid context userID")
	}

	return userID, nil
}

func (m *mw) GetValidator(c echo.Context) (*validator.Validate, error) {
	rowValidate := c.Get(validatorKey)
	validate, ok := rowValidate.(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("failed to get validator")
	}

	return validate, nil
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
