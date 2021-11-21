package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type Target struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func NewTarget(sqlHandler infrastructure.SqlHandler) *Target {
	return &Target{SqlHandler: sqlHandler}
}

func (t *Target) InsertTargets(ctx context.Context, questionnaireID int, targets []string) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	if len(targets) == 0 {
		return nil
	}

	dbTargets := make([]model.Targets, 0, len(targets))
	for _, target := range targets {
		dbTargets = append(dbTargets, model.Targets{
			QuestionnaireID: questionnaireID,
			UserTraqid:      target,
		})
	}

	err = db.Create(&dbTargets).Error
	if err != nil {
		return fmt.Errorf("failed to insert targets :%w", err)
	}

	return nil
}

func (t *Target) DeleteTargets(ctx context.Context, questionnaireID int) error {
	panic("implement me")
}

func (t *Target) GetTargets(ctx context.Context, questionnaireIDs []int) ([]model.Targets, error) {
	panic("implement me")
}
