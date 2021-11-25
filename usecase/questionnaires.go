package usecase

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/xxarupkaxx/anke-two/domain/model"
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

func NewQuestionnaire(IQuestionnaire repository.IQuestionnaire, ITarget repository.ITarget, IAdministrator repository.IAdministrator, IQuestion repository.IQuestion, IOption repository.IOption, IScaleLabel repository.IScaleLabel, IValidation repository.IValidation, ITransaction transaction.ITransaction, IWebhook traq.IWebhook) QuestionnaireUsecase {
	return &questionnaire{IQuestionnaire: IQuestionnaire, ITarget: ITarget, IAdministrator: IAdministrator, IQuestion: IQuestion, IOption: IOption, IScaleLabel: IScaleLabel, IValidation: IValidation, ITransaction: ITransaction, IWebhook: IWebhook}
}

type QuestionnaireUsecase interface {
	POSTQuestionnaire(ctx context.Context, input input.PostAndEditQuestionnaireRequest) (output.PostAndEditQuestionnaireRequest, error)
	GetQuestionnaires(ctx context.Context, param input.GetQuestionnairesQueryParam) (output.GetQuestionnaire, error)
	GetQuestionnaire(ctx context.Context, getQuestionnaire input.GetQuestionnaire) (output.GetQuestionnaire, error)
	PostQuestionByQuestionnaireID(ctx context.Context, request input.PostQuestionRequest) (output.PostQuestionRequest, error)
	EditQuestionnaire(ctx context.Context, request input.PostAndEditQuestionnaireRequest) error
	DeleteQuestionnaire(ctx context.Context) error
	GetQuestions(ctx context.Context, info input.QuestionInfo) (output.QuestionInfo, error)
}

func (q *questionnaire) POSTQuestionnaire(ctx context.Context, input input.PostAndEditQuestionnaireRequest) (output.PostAndEditQuestionnaireRequest, error) {
	if input.ResTimeLimit.Valid {
		isBefore := input.ResTimeLimit.ValueOrZero().Before(time.Now())
		if isBefore {
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

func (q *questionnaire) GetQuestionnaire(c echo.Context, getQuestionnaire input.GetQuestionnaire) (output.GetQuestionnaire, error) {
	panic("implement me")
}

func (q *questionnaire) PostQuestionByQuestionnaireID(c echo.Context, request input.PostQuestionRequest) (output.PostQuestionRequest, error) {
	panic("implement me")
}

func (q *questionnaire) EditQuestionnaire(c echo.Context, request input.PostAndEditQuestionnaireRequest) error {
	panic("implement me")
}

func (q *questionnaire) DeleteQuestionnaire(c echo.Context) error {
	panic("implement me")
}

func (q *questionnaire) GetQuestions(c echo.Context, info input.QuestionInfo) (output.QuestionInfo, error) {
	allQuestions, err := q.IQuestion.GetQuestions(c.Request().Context(), info.QuestionnaireID)
	if err != nil {
		return output.QuestionInfo{StatusCode: http.StatusInternalServerError}, err
	}
	if len(allQuestions) == 0 {
		return output.QuestionInfo{StatusCode: http.StatusNotFound}, nil
	}

	optionIDs := []int{}
	scaleLabelIDs := []int{}
	validationIDs := []int{}
	var questionsType map[int]model.QuestionType
	for _, question := range allQuestions {
		questionsType,err = q.IQuestion.ChangeStrQuestionType(c.Request().Context(),question.QuestionnaireID)
		if err != nil {
			return output.QuestionInfo{StatusCode: http.StatusInternalServerError},err
		}
	}

	for questionID, questionType := range questionsType {
		switch questionType.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			optionIDs = append(optionIDs, questionID)
		case "LinearScale":
			scaleLabelIDs = append(scaleLabelIDs, questionID)
		case "Text", "Number":
			validationIDs = append(validationIDs, questionID)
		}
	}

	options, err := q.IOption.GetOptions(c.Request().Context(), optionIDs)
	if err != nil {
		return output.QuestionInfo{StatusCode: http.StatusInternalServerError},err
	}
	optionMap := make(map[int][]string, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option.Body)
	}

	scaleLabels, err := q.IScaleLabel.GetScaleLabels(c.Request().Context(), scaleLabelIDs)
	if err != nil {
		return output.QuestionInfo{StatusCode: http.StatusInternalServerError},err
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = label
	}

	validations, err := q.IValidation.GetValidations(c.Request().Context(), validationIDs)
	if err != nil {
		return output.QuestionInfo{StatusCode: http.StatusInternalServerError},err
	}
	validationMap := make(map[int]model.Validations, len(validations))
	for _, validation := range validations {
		validationMap[validation.QuestionID] = validation
	}

}
