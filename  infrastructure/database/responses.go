package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

type Response struct {
	//TODO:使うかどうかはこれから
	infrastructure.SqlHandler
}

func NewResponse(sqlHandler infrastructure.SqlHandler) *Response {
	return &Response{SqlHandler: sqlHandler}
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
