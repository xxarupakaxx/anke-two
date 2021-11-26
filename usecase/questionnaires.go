package usecase

import (
	"context"
	"errors"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	myMiddleware "github.com/xxarupkaxx/anke-two/domain/repository/middleware"
	"github.com/xxarupkaxx/anke-two/domain/repository/transaction"
	"github.com/xxarupkaxx/anke-two/domain/repository/traq"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"gorm.io/gorm"
	"net/http"
	"regexp"
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

func (q *questionnaire) ValidationPostQuestionByQuestionnaireID(request input.PostQuestionRequest) error {
	switch request.QuestionType {
	case "Text":
		if _, err := regexp.Compile(request.RegexPattern); err != nil {
			return err
		}
	case "Number":
		if err := q.IValidation.CheckNumberValid(request.MinBound, request.MaxBound); err != nil {
			return err
		}
	}
	return nil
}

func NewQuestionnaire(IQuestionnaire repository.IQuestionnaire, ITarget repository.ITarget, IAdministrator repository.IAdministrator, IQuestion repository.IQuestion, IOption repository.IOption, IScaleLabel repository.IScaleLabel, IValidation repository.IValidation, ITransaction transaction.ITransaction, IWebhook traq.IWebhook) QuestionnaireUsecase {
	return &questionnaire{IQuestionnaire: IQuestionnaire, ITarget: ITarget, IAdministrator: IAdministrator, IQuestion: IQuestion, IOption: IOption, IScaleLabel: IScaleLabel, IValidation: IValidation, ITransaction: ITransaction, IWebhook: IWebhook}
}

type QuestionnaireUsecase interface {
	PostQuestionnaire(ctx context.Context, input input.PostAndEditQuestionnaireRequest) (output.PostQuestionnaireRequest, error)
	GetQuestionnaires(ctx context.Context, param input.GetQuestionnairesQueryParam) (output.GetQuestionnaires, error)
	GetQuestionnaire(ctx context.Context, getQuestionnaire input.GetQuestionnaire) (output.GetQuestionnaire, error)
	PostQuestionByQuestionnaireID(ctx context.Context, request input.PostQuestionRequest) (output.PostQuestionRequest, error)
	ValidationPostQuestionByQuestionnaireID(request input.PostQuestionRequest) error
	EditQuestionnaire(ctx context.Context, request input.PostAndEditQuestionnaireRequest) error
	DeleteQuestionnaire(ctx context.Context, request input.DeleteQuestionnaire) error
	GetQuestions(ctx context.Context, info input.QuestionInfo) ([]output.QuestionInfo, error)
}

func (q *questionnaire) PostQuestionnaire(ctx context.Context, input input.PostAndEditQuestionnaireRequest) (output.PostQuestionnaireRequest, error) {
	if input.ResTimeLimit.Valid {
		isBefore := input.ResTimeLimit.ValueOrZero().Before(time.Now())
		if isBefore {
			return output.PostQuestionnaireRequest{}, model.ErrResTimeBefore
		}
	}

	var questionnaireID int
	var err error

	err = q.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		questionnaireID, err = q.InsertQuestionnaire(ctx, input.Title, input.Description, input.ResTimeLimit, input.ResSharedTo)
		if err != nil {
			return err
		}

		err = q.InsertTargets(ctx, questionnaireID, input.Targets)
		if err != nil {
			return err
		}

		err = q.InsertAdministrator(ctx, questionnaireID, input.Administrators)
		if err != nil {
			return err
		}

		message := q.CreateQuestionnaireMessage(questionnaireID, input.Title, input.Description, input.Administrators, input.ResTimeLimit, input.Targets)

		err = q.PostMessage(message)
		if err != nil {
			return model.ErrFailedPostMessage
		}

		return nil
	})
	if err != nil {
		return output.PostQuestionnaireRequest{}, err
	}

	now := time.Now()
	outputQuestionnaire := output.PostQuestionnaireRequest{
		QuestionnaireID: questionnaireID,
		Title:           input.Title,
		Description:     input.Description,
		ResTimeLimit:    input.ResTimeLimit,
		DeletedAt:       gorm.DeletedAt{},
		CreatedAt:       now.Format(time.RFC3339),
		ModifiedAt:      now.Format(time.RFC3339),
		ResSharedTo:     input.ResSharedTo,
		Targets:         input.Targets,
		Administrators:  input.Administrators,
	}

	return outputQuestionnaire, nil

}

func (q *questionnaire) GetQuestionnaires(ctx context.Context, param input.GetQuestionnairesQueryParam) (output.GetQuestionnaires, error) {
	questionnaires, pageMax, err := q.IQuestionnaire.GetQuestionnaires(ctx, param.UserID, param.Sort, param.Search, param.Page, param.Nontargeted)
	if err != nil {
		return output.GetQuestionnaires{}, err
	}

	outputGetQuestionnaire := output.GetQuestionnaires{
		PageMax:        pageMax,
		Questionnaires: questionnaires,
	}
	return outputGetQuestionnaire, nil
}

func (q *questionnaire) GetQuestionnaire(ctx context.Context, getQuestionnaire input.GetQuestionnaire) (output.GetQuestionnaire, error) {
	qe, targets, administrators, respondents, err := q.IQuestionnaire.GetQuestionnaireInfo(ctx, getQuestionnaire.QuestionnaireID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return output.GetQuestionnaire{StatusCode: http.StatusNotFound}, err
		}
		return output.GetQuestionnaire{StatusCode: http.StatusInternalServerError}, err
	}

	outputQ := output.GetQuestionnaire{
		QuestionnaireID: qe.ID,
		Title:           qe.Title,
		Description:     qe.Description,
		ResTimeLimit:    qe.ResTimeLimit,
		CreatedAt:       qe.CreatedAt.Format(time.RFC3339),
		ModifiedAt:      qe.ModifiedAt.Format(time.RFC3339),
		ResSharedTo:     qe.ResSharedTo,
		Targets:         targets,
		Administrators:  administrators,
		Respondents:     respondents,
	}

	return outputQ, nil
}

func (q *questionnaire) PostQuestionByQuestionnaireID(ctx context.Context, request input.PostQuestionRequest) (output.PostQuestionRequest, error) {
	var opQuestion output.PostQuestionRequest
	err := q.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		lastID, err := q.IQuestion.InsertQuestion(ctx, request.QuestionnaireID, request.PageNum, request.QuestionNum, request.QuestionType, request.Body, request.IsRequired)
		if err != nil {
			return err
		}

		switch request.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for i, option := range request.Options {
				if err := q.IOption.InsertOption(ctx, lastID, i+1, option); err != nil {
					return err
				}
			}
		case "LinearScale":
			if err := q.IScaleLabel.InsertScaleLabel(ctx, lastID, model.ScaleLabels{
				ScaleLabelRight: request.ScaleLabelRight,
				ScaleLabelLeft:  request.ScaleLabelLeft,
				ScaleMin:        request.ScaleMin,
				ScaleMax:        request.ScaleMax,
			}); err != nil {
				return err
			}
		case "Text", "Number":
			if err := q.IValidation.InsertValidation(ctx, lastID, model.Validations{
				QuestionID:   0,
				RegexPattern: request.RegexPattern,
				MinBound:     request.MinBound,
				MaxBound:     request.MaxBound,
			}); err != nil {
				return err
			}
		}

		opQuestion = output.PostQuestionRequest{
			QuestionID:      lastID,
			QuestionType:    request.QuestionType,
			QuestionNum:     request.QuestionNum,
			PageNum:         request.PageNum,
			Body:            request.Body,
			IsRequired:      request.IsRequired,
			Options:         request.Options,
			ScaleLabelRight: request.ScaleLabelRight,
			ScaleLabelLeft:  request.ScaleLabelLeft,
			ScaleMin:        request.ScaleMin,
			ScaleMax:        request.ScaleMax,
			RegexPattern:    request.RegexPattern,
			MinBound:        request.MinBound,
			MaxBound:        request.MaxBound,
		}

		return nil
	})
	if err != nil {
		return output.PostQuestionRequest{}, err
	}

	return opQuestion, nil
}

func (q *questionnaire) EditQuestionnaire(ctx context.Context, request input.PostAndEditQuestionnaireRequest) error {
	err := q.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		err := q.IQuestionnaire.UpdateQuestionnaire(ctx, request.Title, request.Description, request.ResTimeLimit, request.ResSharedTo, request.QuestionnaireID)
		if err != nil {
			return err
		}
		err = q.ITarget.DeleteTargets(ctx, request.QuestionnaireID)
		if err != nil {
			return err
		}

		err = q.ITarget.InsertTargets(ctx, request.QuestionnaireID, request.Targets)
		if err != nil {
			return err
		}

		err = q.IAdministrator.DeleteAdministrators(ctx, request.QuestionnaireID)
		if err != nil {
			return err
		}

		err = q.IAdministrator.InsertAdministrator(ctx, request.QuestionnaireID, request.Administrators)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.ErrTransaction
	}

	return nil
}

func (q *questionnaire) DeleteQuestionnaire(ctx context.Context, request input.DeleteQuestionnaire) error {
	err := q.ITransaction.Do(ctx, nil, func(c context.Context) error {
		err := q.IQuestionnaire.DeleteQuestionnaire(ctx, request.QuestionnaireID)
		if err != nil {
			return err
		}

		err = q.ITarget.DeleteTargets(ctx, request.QuestionnaireID)
		if err != nil {
			return err
		}

		err = q.IAdministrator.DeleteAdministrators(ctx, request.QuestionnaireID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.ErrTransaction
	}

	return nil
}

func (q *questionnaire) GetQuestions(ctx context.Context, info input.QuestionInfo) ([]output.QuestionInfo, error) {
	allQuestions, err := q.IQuestion.GetQuestions(ctx, info.QuestionnaireID)
	if err != nil {
		return nil, err
	}
	if len(allQuestions) == 0 {
		return nil, model.ErrRecordNotFound
	}

	optionIDs := []int{}
	scaleLabelIDs := []int{}
	validationIDs := []int{}

	for _, question := range allQuestions {
		switch question.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			optionIDs = append(optionIDs, question.ID)
		case "LinearScale":
			scaleLabelIDs = append(scaleLabelIDs, question.ID)
		case "Text", "Number":
			validationIDs = append(validationIDs, question.ID)
		}
	}

	options, err := q.IOption.GetOptions(ctx, optionIDs)
	if err != nil {
		return nil, err
	}
	optionMap := make(map[int][]string, len(options))
	for _, option := range options {
		optionMap[option.QuestionID] = append(optionMap[option.QuestionID], option.Body)
	}

	scaleLabels, err := q.IScaleLabel.GetScaleLabels(ctx, scaleLabelIDs)
	if err != nil {
		return nil, err
	}
	scaleLabelMap := make(map[int]model.ScaleLabels, len(scaleLabels))
	for _, label := range scaleLabels {
		scaleLabelMap[label.QuestionID] = label
	}

	validations, err := q.IValidation.GetValidations(ctx, validationIDs)
	if err != nil {
		return nil, err
	}
	validationMap := make(map[int]model.Validations, len(validations))
	for _, validation := range validations {
		validationMap[validation.QuestionID] = validation
	}

	var outputs []output.QuestionInfo

	for _, question := range allQuestions {
		options := []string{}
		scalelabel := model.ScaleLabels{}
		validation := model.Validations{}
		switch question.Type {
		case "MultipleChoice", "Checkbox", "Dropdown":
			var ok bool
			options, ok = optionMap[question.ID]
			if !ok {
				options = []string{}
			}
		case "LinearScale":
			var ok bool
			scalelabel, ok = scaleLabelMap[question.ID]
			if !ok {
				scalelabel = model.ScaleLabels{}
			}
		case "Text", "Number":
			var ok bool
			validation, ok = validationMap[question.ID]
			if !ok {
				validation = model.Validations{}
			}
		}

		outputs = append(outputs, output.QuestionInfo{
			QuestionID:      question.ID,
			QuestionType:    question.Type,
			QuestionNum:     question.QuestionNum,
			PageNum:         question.PageNum,
			Body:            question.Body,
			IsRequired:      question.IsRequired,
			CreatedAt:       question.CreatedAt.Format(time.RFC3339),
			Options:         options,
			ScaleLabelRight: scalelabel.ScaleLabelRight,
			ScaleLabelLeft:  scalelabel.ScaleLabelLeft,
			ScaleMin:        scalelabel.ScaleMin,
			ScaleMax:        scalelabel.ScaleMax,
			RegexPattern:    validation.RegexPattern,
			MinBound:        validation.MinBound,
			MaxBound:        validation.MaxBound,
		})
	}

	return outputs, nil

}
