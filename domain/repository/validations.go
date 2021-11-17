package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

// IValidation ValidationのRepository
type IValidation interface {
	InsertValidation(ctx context.Context, lastID int, validation model.Validations) error
	UpdateValidation(ctx context.Context, questionID int, validation model.Validations) error
	DeleteValidation(ctx context.Context, questionID int) error
	GetValidations(ctx context.Context, qustionIDs []int) ([]model.Validations, error)
	CheckNumberValidation(validation model.Validations, Body string) error
	CheckTextValidation(validation model.Validations, Response string) error
	CheckNumberValid(MinBound, MaxBound string) error
}
