package usecase

import (
	"context"
	"errors"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"net/http"
	"regexp"
)

type question struct {
	repository.IValidation
	repository.IOption
	repository.IQuestion
	repository.IScaleLabel
}

func (q *question) ValidationEditQuestion(request input.EditQuestionRequest) error {
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
}

func (q *question) EditQuestion(ctx context.Context, request input.EditQuestionRequest) (output.EditQuestion, error) {
	err := q.IQuestion.UpdateQuestion(ctx, request.QuestionnaireID, request.PageNum, request.QuestionNum, request.QuestionType, request.Body, request.IsRequired, request.QuestionID)
	if err != nil {
		return output.EditQuestion{StatusCode: http.StatusInternalServerError}, err
	}

	switch request.QuestionType {
	case "MultipleChoice", "Checkbox", "Dropdown":
		if err := q.IOption.UpdateOptions(ctx, request.Options, request.QuestionID); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			return output.EditQuestion{StatusCode: http.StatusInternalServerError}, err
		}
	case "LinearScale":
		if err := q.IScaleLabel.UpdateScaleLabel(ctx, request.QuestionID, model.ScaleLabels{
			ScaleLabelLeft:  request.ScaleLabelLeft,
			ScaleLabelRight: request.ScaleLabelRight,
			ScaleMax:        request.ScaleMax,
			ScaleMin:        request.ScaleMin,
		}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			return output.EditQuestion{StatusCode: http.StatusInternalServerError}, err
		}
	case "Text", "Number":
		if err := q.IValidation.UpdateValidation(ctx, request.QuestionID, model.Validations{
			QuestionID:   0,
			RegexPattern: request.RegexPattern,
			MinBound:     request.MinBound,
			MaxBound:     request.MaxBound,
		}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
			return output.EditQuestion{StatusCode: http.StatusInternalServerError}, err
		}
	}

	return output.EditQuestion{StatusCode: http.StatusOK}, nil

}

func (q *question) DeleteQuestion(ctx context.Context, deleteQuestion input.DeleteQuestion) (output.DeleteQuestion, error) {
	if err := q.IQuestion.DeleteQuestion(ctx, deleteQuestion.QuestionID); err != nil {
		return output.DeleteQuestion{StatusCode: http.StatusInternalServerError}, err
	}

	if err := q.IOption.DeleteOptions(ctx, deleteQuestion.QuestionID); err != nil {
		return output.DeleteQuestion{StatusCode: http.StatusInternalServerError}, err
	}

	if err := q.IScaleLabel.DeleteScaleLabel(ctx, deleteQuestion.QuestionID); err != nil {
		return output.DeleteQuestion{StatusCode: http.StatusInternalServerError}, err
	}

	if err := q.IValidation.DeleteValidation(ctx, deleteQuestion.QuestionID); err != nil {
		return output.DeleteQuestion{StatusCode: http.StatusInternalServerError}, err
	}

	return output.DeleteQuestion{StatusCode: http.StatusOK}, nil
}

func NewQuestionUsecase(IValidation repository.IValidation, IOption repository.IOption, IQuestion repository.IQuestion, IScaleLabel repository.IScaleLabel) QuestionUsecase {
	return &question{IValidation: IValidation, IOption: IOption, IQuestion: IQuestion, IScaleLabel: IScaleLabel}
}

type QuestionUsecase interface {
	EditQuestion(ctx context.Context, request input.EditQuestionRequest) (output.EditQuestion, error)
	ValidationEditQuestion(request input.EditQuestionRequest) error
	DeleteQuestion(c context.Context, deleteQuestion input.DeleteQuestion) (output.DeleteQuestion, error)
}
