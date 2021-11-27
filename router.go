package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func SetRouting(port string, db *gorm.DB) {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	api := InjectAPIServer(db)

	e.Static("/", "client/dist")
	e.Static("/js", "client/dist/js")
	e.Static("/img", "client/dist/img")
	e.Static("/fonts", "client/dist/fonts")
	e.Static("/css", "client/dist/css")

	e.File("/app.js", "client/dist/app.js")
	e.File("/favicon.ico", "client/dist/favicon.ico")
	e.File("*", "client/dist/index.html")

	echoAPI := e.Group("/api", api.SetValidatorMiddleware, api.SetUserIDMiddleware, api.TraPMemberAuthenticate)
	{
		apiQuestionnaires := echoAPI.Group("/questionnaires")
		{
			apiQuestionnaires.GET("", api.GetQuestionnaires, api.TrapReteLimitMiddlewareFunc())
			apiQuestionnaires.POST("", api.PostQuestionnaire)
			apiQuestionnaires.GET("/:questionnaireID", api.GetQuestionnaire)
			apiQuestionnaires.PATCH("/:questionnaireID", api.EditQuestionnaire, api.QuestionnaireAdministratorAuthenticate)
			apiQuestionnaires.DELETE("/:questionnaireID", api.DeleteQuestionnaire, api.QuestionnaireAdministratorAuthenticate)
			apiQuestionnaires.GET("/:questionnaireID/questions", api.GetQuestions)
			apiQuestionnaires.POST("/:questionnaireID/questions", api.PostQuestionByQuestionnaireID)
		}

		apiQuestions := echoAPI.Group("/questions")
		{
			apiQuestions.PATCH("/:questionID", api.EditQuestion, api.QuestionAdministratorAuthenticate)
			apiQuestions.DELETE("/:questionID", api.DeleteQuestion, api.QuestionAdministratorAuthenticate)
		}

		apiResponses := echoAPI.Group("/responses")
		{
			apiResponses.POST("", api.PostResponse)
			apiResponses.GET("/:responseID", api.GetResponse, api.ResponseReadAuthenticate)
			apiResponses.PATCH("/:responseID", api.EditResponse, api.ResponseReadAuthenticate)
			apiResponses.DELETE("/:responseID", api.DeleteResponse, api.ResponseReadAuthenticate)
		}

		apiUsers := echoAPI.Group("/users")
		{
			apiUsersMe := apiUsers.Group("/me")
			{
				apiUsersMe.GET("", api.GetUsersMe)
				apiUsersMe.GET("/responses", api.GetMyResponse)
				apiUsersMe.GET("/responses/:questionnaireID", api.GetMyResponsesByID)
				apiUsersMe.GET("/targeted", api.GetTargetedQuestionnaire)
				apiUsersMe.GET("/administrates", api.GetMyQuestionnaire)
			}
			apiUsers.GET("/:traQID/targeted", api.GetTargetedQuestionnairesByTraQID)
		}

		apiResults := echoAPI.Group("/results")
		{
			apiResults.GET("/:questionnaireID", api.GetResults, api.ResultAuthenticate)
		}
	}

	e.Logger.Fatal(e.Start(port))
}
