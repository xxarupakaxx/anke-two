// +build wireinject
package main

import (
	"github.com/google/wire"
	"github.com/xxarupkaxx/anke-two/infrastructure/database"
	"github.com/xxarupkaxx/anke-two/infrastructure/traq"
	middleware2 "github.com/xxarupkaxx/anke-two/interfaces/middleware"
	repository2 "github.com/xxarupkaxx/anke-two/interfaces/repository"
	middleware3 "github.com/xxarupkaxx/anke-two/interfaces/repository/middleware"
	transaction2 "github.com/xxarupkaxx/anke-two/interfaces/repository/transaction"
	traq3 "github.com/xxarupkaxx/anke-two/interfaces/repository/traq"
	"github.com/xxarupkaxx/anke-two/interfaces/router"
	"github.com/xxarupkaxx/anke-two/usecase"
	"gorm.io/gorm"
)

var superSet = wire.NewSet(
	database.NewAdministrator,
	wire.Bind(new(repository2.IAdministrator), new(*database.Administrator)),
	database.NewOption,
	wire.Bind(new(repository2.IOption), new(*database.Option)),
	database.NewQuestionnaire,
	wire.Bind(new(repository2.IQuestionnaire), new(*database.Questionnaire)),
	database.NewQuestion,
	wire.Bind(new(repository2.IQuestion), new(*database.Question)),
	database.NewRespondent,
	wire.Bind(new(repository2.IRespondent), new(*database.Respondent)),
	database.NewResponse,
	wire.Bind(new(repository2.IResponse), new(*database.Response)),
	database.NewScaleLabel,
	wire.Bind(new(repository2.IScaleLabel), new(*database.ScaleLabel)),
	database.NewTarget,
	wire.Bind(new(repository2.ITarget), new(*database.Target)),
	database.NewTransaction,
	wire.Bind(new(transaction2.ITransaction), new(*database.Tx)),
	database.NewValidation,
	wire.Bind(new(repository2.IValidation), new(*database.Validation)),
	traq.NewWebhook,
	wire.Bind(new(traq3.IWebhook), new(*traq.Webhook)),

	middleware2.NewMiddleware,
	wire.Bind(new(middleware3.IMiddleware), new(*middleware2.Mv)),

	usecase.NewQuestion,
	wire.Bind(new(usecase.QuestionUsecase), new(*usecase.Question)),
	usecase.NewQuestionnaire,
	wire.Bind(new(usecase.QuestionnaireUsecase), new(*usecase.Questionnaire)),
	usecase.NewResponse,
	wire.Bind(new(usecase.ResponseUsecase), new(*usecase.Response)),
	usecase.NewResult,
	wire.Bind(new(usecase.ResultUsecase), new(*usecase.Result)),
	usecase.NewUser,
	wire.Bind(new(usecase.UsersUsecase), new(*usecase.User)),

	router.NewAPI,
	router.NewQuestionnaireAPI,
	wire.Bind(new(router.QuestionnaireAPI), new(*router.Questionnaire)),
	router.NewQuestionAPI,
	wire.Bind(new(router.QuestionAPI), new(*router.Question)),
	router.NewResponseAPI,
	wire.Bind(new(router.ResponseAPI), new(*router.Response)),
	router.NewResultAPI,
	wire.Bind(new(router.ResultAPI), new(*router.Result)),
	router.NewUserAPI,
	wire.Bind(new(router.UserAPI), new(*router.User)),

)

func injectAPIServer(db *gorm.DB) *router.API {
	wire.Build(superSet)

	return nil
}
