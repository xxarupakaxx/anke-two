package usecase

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	myMiddleware "github.com/xxarupkaxx/anke-two/domain/repository/middleware"
	"github.com/xxarupkaxx/anke-two/domain/repository/transaction"
	"github.com/xxarupkaxx/anke-two/domain/repository/traq"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type questionnaire struct {
	repository.IQuestionnaire
	repository.ITarget
	repository.IAdministrator
	repository.IQuestion
	repository.IOption
	repository.IScaleLabel
	repository.IValidation
	transaction.ITransaction
	myMiddleware.IMiddleware
	traq.IWebhook
}

func (q *questionnaire) GetQuestionnaires(c echo.Context, param input.GetQuestionnairesQueryParam) (output.GetQuestionnaire, error) {
	questionnaires, pageMax, err := q.IQuestionnaire.GetQuestionnaires(c.Request().Context(), param.UserID, param.Sort, param.Search, param.Page, param.Nontargeted)
	if err != nil {
		return output.GetQuestionnaire{}, err
	}

	outputGetQuestionnaire := output.GetQuestionnaire{
		PageMax:        pageMax,
		Questionnaires: questionnaires,
	}
	return outputGetQuestionnaire, nil
}

func (q *questionnaire) ValidateGetQuestionnaires(c echo.Context, param input.GetQuestionnairesQueryParam) (int, error) {
	validate, err := q.GetValidator(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = validate.StructCtx(c.Request().Context(), param)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func NewQuestionnaire(IQuestionnaire repository.IQuestionnaire, ITarget repository.ITarget, IAdministrator repository.IAdministrator, IQuestion repository.IQuestion, IOption repository.IOption, IScaleLabel repository.IScaleLabel, IValidation repository.IValidation, ITransaction transaction.ITransaction, IWebhook traq.IWebhook) Questionnaire {
	return &questionnaire{IQuestionnaire: IQuestionnaire, ITarget: ITarget, IAdministrator: IAdministrator, IQuestion: IQuestion, IOption: IOption, IScaleLabel: IScaleLabel, IValidation: IValidation, ITransaction: ITransaction, IWebhook: IWebhook}
}

type Questionnaire interface {
	POSTQuestionnaire(c echo.Context, input input.PostAndEditQuestionnaireRequest) (output.PostAndEditQuestionnaireRequest, error)
	ValidatePostQuestionnaire(c echo.Context, input input.PostAndEditQuestionnaireRequest) (int, error)
	GetQuestionnaires(c echo.Context, param input.GetQuestionnairesQueryParam) (output.GetQuestionnaire, error)
	ValidateGetQuestionnaires(c echo.Context, param input.GetQuestionnairesQueryParam) (int, error)
}

func (q *questionnaire) POSTQuestionnaire(c echo.Context, input input.PostAndEditQuestionnaireRequest) (output.PostAndEditQuestionnaireRequest, error) {
	if input.ResTimeLimit.Valid {
		isBefore := input.ResTimeLimit.ValueOrZero().Before(time.Now())
		if isBefore {
			c.Logger().Infof("invalid resTimeLimit: %+v", input.ResTimeLimit)
			return output.PostAndEditQuestionnaireRequest{}, echo.NewHTTPError(http.StatusBadRequest, "res time limit is before now")
		}
	}

	var questionnaireID int
	var err error

	err = q.ITransaction.Do(c.Request().Context(), nil, func(ctx context.Context) error {
		questionnaireID, err = q.InsertQuestionnaire(ctx, input.Title, input.Description, input.ResTimeLimit, input.ResSharedTo)
		if err != nil {
			c.Logger().Errorf("failed to insert a questionnaire:%w", err)
			return err
		}

		err = q.InsertTargets(ctx, questionnaireID, input.Targets)
		if err != nil {
			c.Logger().Errorf("failed to insert targets:%w", err)
			return err
		}

		err = q.InsertAdministrator(ctx, questionnaireID, input.Administrators)
		if err != nil {
			c.Logger().Errorf("failed to insert administrators:%w", err)
			return err
		}

		message := q.CreateQuestionnaireMessage(questionnaireID, input.Title, input.Description, input.Administrators, input.ResTimeLimit, input.Targets)

		err = q.PostMessage(message)
		if err != nil {
			c.Logger().Errorf("failed to post message: %w", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to post message to traQ")
		}

		return nil
	})
	if err != nil {
		return output.PostAndEditQuestionnaireRequest{}, err
	}

	outputQuestionnaire := output.PostAndEditQuestionnaireRequest{
		QuestionnaireID: questionnaireID,
		Title:           input.Title,
		Description:     input.Description,
		ResTimeLimit:    input.ResTimeLimit,
		DeletedAt:       gorm.DeletedAt{},
		ResSharedTo:     input.ResSharedTo,
		Targets:         input.Targets,
		Administrators:  input.Administrators,
	}

	return outputQuestionnaire, nil

}

func (q *questionnaire) ValidatePostQuestionnaire(c echo.Context, input input.PostAndEditQuestionnaireRequest) (int, error) {
	validate, err := q.GetValidator(c)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = validate.StructCtx(c.Request().Context(), input)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
