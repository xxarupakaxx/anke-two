package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

// IScaleLabel ScaleLabelのRepository
type IScaleLabel interface {
	InsertScaleLabel(ctx context.Context, lastID int, label model.ScaleLabels) error
	UpdateScaleLabel(ctx context.Context, questionID int, label model.ScaleLabels) error
	DeleteScaleLabel(ctx context.Context, questionID int) error
	GetScaleLabels(ctx context.Context, questionIDs []int) ([]model.ScaleLabels, error)
	CheckScaleLabel(label model.ScaleLabels, response string) error
}
