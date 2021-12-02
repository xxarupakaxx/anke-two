package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	repository2 "github.com/xxarupkaxx/anke-two/interfaces/repository"
	transaction2 "github.com/xxarupkaxx/anke-two/interfaces/repository/transaction"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Response struct {
	repository2.IRespondent
	repository2.IQuestionnaire
	repository2.IValidation
	repository2.IScaleLabel
	repository2.IResponse
	transaction2.ITransaction
}

func NewResponse(IRespondent repository2.IRespondent, IQuestionnaire repository2.IQuestionnaire, IValidation repository2.IValidation, IScaleLabel repository2.IScaleLabel, IResponse repository2.IResponse, ITransaction transaction2.ITransaction) *Response {
	return &Response{IRespondent: IRespondent, IQuestionnaire: IQuestionnaire, IValidation: IValidation, IScaleLabel: IScaleLabel, IResponse: IResponse, ITransaction: ITransaction}
}

func (r *Response) PostResponse(ctx context.Context, responses input.Responses) (output.PostResponse, error) {
	var submittedAt time.Time
	var responseID int

	err := r.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		limit, err := r.IQuestionnaire.GetQuestionnaireLimit(ctx, responses.ID)
		if err != nil {
			return err
		}

		if limit.Valid && limit.Time.Before(time.Now()) {
			return err
		}

		questionIDs := make([]int, 0, len(responses.Body))
		QuestionTypes := make(map[int]model.ResponseBody, len(responses.Body))

		for _, body := range responses.Body {
			questionIDs = append(questionIDs, body.QuestionID)
			QuestionTypes[body.QuestionID] = body
		}

		validation, err := r.IValidation.GetValidations(ctx, questionIDs)
		if err != nil {
			return err
		}

		for _, v := range validation {
			body := QuestionTypes[v.QuestionID]
			switch body.QuestionType {
			case "Number":
				if err := r.IValidation.CheckNumberValidation(v, body.Body.ValueOrZero()); err != nil {
					return err
				}
			case "Text":
				if err := r.IValidation.CheckTextValidation(v, body.Body.ValueOrZero()); err != nil {
					return err
				}
			}

		}
		scaleLabelIDs := []int{}
		for _, body := range responses.Body {
			switch body.QuestionType {
			case "LinearScale":
				scaleLabelIDs = append(scaleLabelIDs, body.QuestionID)
			}
		}

		scaleLabels, err := r.IScaleLabel.GetScaleLabels(ctx, scaleLabelIDs)
		if err != nil {
			return err
		}

		scaleLabelMap := make(map[int]*model.ScaleLabels, len(scaleLabels))

		for _, label := range scaleLabels {
			scaleLabelMap[label.QuestionID] = &label
		}

		for _, body := range responses.Body {
			switch body.QuestionType {
			case "LinearScale":
				label, ok := scaleLabelMap[body.QuestionID]
				if !ok {
					label = &model.ScaleLabels{}
				}
				if err := r.IScaleLabel.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
					return err
				}
			}
		}

		if responses.Temporarily {
			submittedAt = time.Time{}
		} else {
			submittedAt = time.Now()
		}

		responseID, err = r.IRespondent.InsertRespondent(ctx, responses.UserID, responses.ID, null.NewTime(submittedAt, !responses.Temporarily))
		if err != nil {
			return err
		}

		responseMetas := make([]*model.ResponseMeta, 0, len(responses.Body))
		for _, body := range responses.Body {
			switch body.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				for _, option := range body.OptionResponse {
					responseMetas = append(responseMetas, &model.ResponseMeta{
						QuestionID: body.QuestionID,
						Data:       option,
					})
				}
			default:
				responseMetas = append(responseMetas, &model.ResponseMeta{
					QuestionID: body.QuestionID,
					Data:       body.Body.ValueOrZero(),
				})
			}
		}

		err = r.IResponse.InsertResponses(ctx, responseID, responseMetas)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return output.PostResponse{}, err
	}

	op := output.PostResponse{
		ResponseID:      responseID,
		QuestionnaireID: responses.ID,
		Temporarily:     responses.Temporarily,
		SubmittedAt:     submittedAt,
		Body:            responses.Body,
	}
	return op, nil
}

func (r *Response) GetResponse(ctx context.Context, getResponse input.GetResponse) (model.RespondentDetail, error) {
	respondentDetail, err := r.IRespondent.GetRespondentDetail(ctx, getResponse.ResponseID)
	if err != nil {
		return model.RespondentDetail{}, err
	}

	return respondentDetail, nil
}

func (r *Response) EditResponse(ctx context.Context, editResponse input.EditResponse) error {
	err := r.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		limit, err := r.IQuestionnaire.GetQuestionnaireLimit(ctx, editResponse.ID)
		if err != nil {
			return err
		}

		if limit.Valid && limit.Time.Before(time.Now()) {
			return err
		}

		questionIDs := make([]int, 0, len(editResponse.Body))
		questionTypes := make(map[int]model.ResponseBody, len(editResponse.Body))

		for _, body := range editResponse.Body {
			questionIDs = append(questionIDs, body.QuestionID)
			questionTypes[body.QuestionID] = body
		}

		validations, err := r.IValidation.GetValidations(ctx, questionIDs)
		if err != nil {
			return err
		}

		for _, validation := range validations {
			body := questionTypes[validation.QuestionID]
			switch body.QuestionType {
			case "Number":
				if err := r.IValidation.CheckNumberValidation(validation, body.Body.ValueOrZero()); err != nil {
					return err
				}
			case "Text":
				if err := r.IValidation.CheckTextValidation(validation, body.Body.ValueOrZero()); err != nil {
					return err
				}
			}
		}

		scaleLabelIDs := []int{}
		for _, body := range editResponse.Body {
			switch body.QuestionType {
			case "LinearScale":
				scaleLabelIDs = append(scaleLabelIDs, body.QuestionID)
			}
		}

		scaleLabels, err := r.IScaleLabel.GetScaleLabels(ctx, questionIDs)
		if err != nil {
			return err
		}

		scaleLabelMap := make(map[int]*model.ScaleLabels, len(scaleLabels))
		for _, label := range scaleLabels {
			scaleLabelMap[label.QuestionID] = &label
		}

		for _, body := range editResponse.Body {
			switch body.QuestionType {
			case "LinearScale":
				label, ok := scaleLabelMap[body.QuestionID]
				if !ok {
					label = &model.ScaleLabels{}
				}
				if err := r.CheckScaleLabel(*label, body.Body.ValueOrZero()); err != nil {
					return err
				}
			}
		}

		if !editResponse.Temporarily {
			err := r.IRespondent.UpdateSubmittedAt(ctx, editResponse.ResponseID)
			if err != nil {
				return err
			}
		}

		if err := r.IResponse.DeleteResponse(ctx, editResponse.ResponseID); err != nil {
			return err
		}

		responseMetas := make([]*model.ResponseMeta, 0, len(editResponse.Body))
		for _, body := range editResponse.Body {
			switch body.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				for _, option := range body.OptionResponse {
					responseMetas = append(responseMetas, &model.ResponseMeta{
						QuestionID: body.QuestionID,
						Data:       option,
					})
				}
			default:
				responseMetas = append(responseMetas, &model.ResponseMeta{
					QuestionID: body.QuestionID,
					Data:       body.Body.ValueOrZero(),
				})
			}
		}

		err = r.IResponse.InsertResponses(ctx, editResponse.ResponseID, responseMetas)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) DeleteResponse(ctx context.Context, deleteResponse input.DeleteResponse) error {
	err := r.ITransaction.Do(ctx, nil, func(ctx context.Context) error {
		limit, err := r.IQuestionnaire.GetQuestionnaireLimitByResponseID(ctx, deleteResponse.ResponseID)
		if err != nil {
			return err
		}

		if limit.Valid && limit.Time.Before(time.Now()) {
			return err
		}

		err = r.IRespondent.DeleteRespondent(ctx, deleteResponse.ResponseID)
		if err != nil {
			return err
		}

		err = r.IResponse.DeleteResponse(ctx, deleteResponse.ResponseID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type ResponseUsecase interface {
	PostResponse(ctx context.Context, responses input.Responses) (output.PostResponse, error)
	GetResponse(ctx context.Context, getResponse input.GetResponse) (model.RespondentDetail, error)
	EditResponse(ctx context.Context, editResponse input.EditResponse) error
	DeleteResponse(ctx context.Context, deleteResponse input.DeleteResponse) error
}
