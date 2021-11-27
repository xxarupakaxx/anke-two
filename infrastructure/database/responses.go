package database

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) InsertResponses(ctx context.Context, responseID int, responseMetas []*model.ResponseMeta) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	responses := make([]model.Responses, 0, len(responseMetas))
	for _, responseMeta := range responseMetas {
		responses = append(responses, model.Responses{
			ResponseID: responseID,
			QuestionID: responseMeta.QuestionID,
			Body:       null.NewString(responseMeta.Data, true),
		})
	}

	err = db.Create(&responses).Error
	if err != nil {
		return fmt.Errorf("failed to insert reponses: %w", err)
	}

	return nil
}

func (r *Response) DeleteResponse(ctx context.Context, responseID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	result := db.
		Where("response_id = ?", responseID).
		Delete(&model.Responses{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete response :%w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete response:%w", model.ErrNoRecordDeleted)
	}

	return nil
}
