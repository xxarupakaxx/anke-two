package usecase

import (
	"context"
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

func (q *question) EditQuestion(ctx context.Context, request input.EditQuestionRequest) (output.EditQuestion, error) {
	switch request.QuestionType {
	case "Text":
		if _, err := regexp.Compile(request.RegexPattern); err != nil {
			return output.EditQuestion{StatusCode: http.StatusBadRequest}, err
		}
	case "Number":
		if err := q.IValidation.CheckNumberValid(request.MinBound, request.MaxBound); err != nil {
			return output.EditQuestion{StatusCode: http.StatusBadRequest}, err
		}
	}
}

func (q *question) DeleteQuestion(ctx context.Context, deleteQuestion input.DeleteQuestion) (output.DeleteQuestion, error) {
	panic("implement me")
}

func NewQuestionUsecase(IValidation repository.IValidation, IOption repository.IOption, IQuestion repository.IQuestion, IScaleLabel repository.IScaleLabel) QuestionUsecase {
	return &question{IValidation: IValidation, IOption: IOption, IQuestion: IQuestion, IScaleLabel: IScaleLabel}
}

type QuestionUsecase interface {
	EditQuestion(ctx context.Context, request input.EditQuestionRequest) (output.EditQuestion, error)
	DeleteQuestion(c context.Context, deleteQuestion input.DeleteQuestion) (output.DeleteQuestion, error)
}
