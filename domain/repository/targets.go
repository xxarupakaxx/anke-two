package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

// ITarget TargetのRepository
type ITarget interface {
	InsertTargets(ctx context.Context, questionnaireID int, targets []string) error
	DeleteTargets(ctx context.Context, questionnaireID int) error
	GetTargets(ctx context.Context, questionnaireIDs []int) ([]model.Targets, error)
}