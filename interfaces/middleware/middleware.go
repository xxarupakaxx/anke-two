package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/repository"
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
