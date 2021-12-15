package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Validation interface {
	GetValidation(ctx context.Context,questionID int) (*model.Validation,error)
	CreateValidation(ctx context.Context,validation *model.Validation) error
	DeleteValidation(ctx context.Context,questionID int) error
	UpdateValidation(ctx context.Context,validation *model.Validation) error
}
