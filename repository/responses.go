package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

// IResponse ResponseのRepository
type IResponse interface {
	InsertResponses(ctx context.Context, responseID int, responseMetas []*model.ResponseMeta) error
	DeleteResponse(ctx context.Context, responseID int) error
}