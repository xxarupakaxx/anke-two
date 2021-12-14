package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Admin interface {
	CreateAdmins(ctx context.Context, questionnaireID int, administrator []string) error
	GetMyAdmins(ctx context.Context, traqID string) ([]*model.Administrator, error)
	DeleteAdmin(ctx context.Context, questionnaireID int) error
}
