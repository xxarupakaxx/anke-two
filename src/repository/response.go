package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Response interface {
	CreateResponse(ctx context.Context, response *model.Response) error
	DeleteResponse(ctx context.Context, responseID int) error
}
