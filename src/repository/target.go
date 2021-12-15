package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Target interface {
	GetTargets(ctx context.Context, questionnaireID int) ([]*model.Target, error)
	CreateTargets(ctx context.Context, questionnaireID int, targets []string) error
	DeleteTargets(ctx context.Context, questionnaireID int) error
}
