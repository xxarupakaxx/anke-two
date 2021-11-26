package usecase

import (
	"context"
	"errors"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/domain/repository"
	"github.com/xxarupkaxx/anke-two/domain/repository/transaction"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"regexp"
)

type question struct {
	repository.IValidation
	repository.IOption
	repository.IQuestion
	repository.IScaleLabel
	transaction.ITransaction
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

	return nil
}

func (q *question) EditQuestion(ctx context.Context, request input.EditQuestionRequest) error {
	err := q.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		err := q.IQuestion.UpdateQuestion(ctx, request.QuestionnaireID, request.PageNum, request.QuestionNum, request.QuestionType, request.Body, request.IsRequired, request.QuestionID)
		if err != nil {
			return err
		}

		switch request.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			if err := q.IOption.UpdateOptions(ctx, request.Options, request.QuestionID); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
				return err
			}
		case "LinearScale":
			if err := q.IScaleLabel.UpdateScaleLabel(ctx, request.QuestionID, model.ScaleLabels{
				ScaleLabelLeft:  request.ScaleLabelLeft,
				ScaleLabelRight: request.ScaleLabelRight,
				ScaleMax:        request.ScaleMax,
				ScaleMin:        request.ScaleMin,
			}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
				return err
			}
		case "Text", "Number":
			if err := q.IValidation.UpdateValidation(ctx, request.QuestionID, model.Validations{
				QuestionID:   0,
				RegexPattern: request.RegexPattern,
				MinBound:     request.MinBound,
				MaxBound:     request.MaxBound,
			}); err != nil && !errors.Is(err, model.ErrNoRecordUpdated) {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (q *question) DeleteQuestion(ctx context.Context, deleteQuestion input.DeleteQuestion) error {
	err := q.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		if err := q.IQuestion.DeleteQuestion(ctx, deleteQuestion.QuestionID); err != nil {
			return err
		}

		if err := q.IOption.DeleteOptions(ctx, deleteQuestion.QuestionID); err != nil {
			return err
		}

		if err := q.IScaleLabel.DeleteScaleLabel(ctx, deleteQuestion.QuestionID); err != nil {
			return err
		}

		if err := q.IValidation.DeleteValidation(ctx, deleteQuestion.QuestionID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func NewQuestionUsecase(IValidation repository.IValidation, IOption repository.IOption, IQuestion repository.IQuestion, IScaleLabel repository.IScaleLabel, ITransaction transaction.ITransaction) QuestionUsecase {
	return &question{IValidation: IValidation, IOption: IOption, IQuestion: IQuestion, IScaleLabel: IScaleLabel, ITransaction: ITransaction}
}

type QuestionUsecase interface {
	EditQuestion(ctx context.Context, request input.EditQuestionRequest) error
	ValidationEditQuestion(request input.EditQuestionRequest) error
	DeleteQuestion(c context.Context, deleteQuestion input.DeleteQuestion) error
}
