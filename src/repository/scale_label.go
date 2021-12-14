package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type ScaleLabel interface {
	GetScaleLabel(ctx context.Context, questionID int) (*model.ScaleLabel, error)
	CreateScaleLabel(ctx context.Context, label *model.ScaleLabel) error
	UpdateScaleLabel(ctx context.Context, label *model.ScaleLabel) error
	DeleteScaleLabel(ctx context.Context, questionID int) error
}
