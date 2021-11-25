package middleware

import (
	"github.com/labstack/echo/v4"
)

type IMiddleware interface {
	SetValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	TraPMemberAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	TrapReteLimitMiddlewareFunc() echo.MiddlewareFunc
	QuestionnaireAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	ResponseReadAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	RespondentsAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	QuestionAdministratorAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	ResultAuthenticate(next echo.HandlerFunc) echo.HandlerFunc
	GetUserID(c echo.Context) (string, error)
}
