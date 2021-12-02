package database

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
)

type Target struct{
	db *gorm.DB
}

func NewTarget(db *gorm.DB) *Target {
	return &Target{db: db}
}

func (t *Target) InsertTargets(ctx context.Context, questionnaireID int, targets []string) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = t.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	if len(targets) == 0 {
		return nil
	}

	dbTargets := make([]Targets, 0, len(targets))
	for _, target := range targets {
		dbTargets = append(dbTargets, Targets{
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
	db, err := GetTx(ctx)
	if db == nil {
		db = t.db
	}

	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Targets{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete targets: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delet response:%w", model.ErrNoRecordDeleted)
	}

	return nil
}

func (t *Target) GetTargets(ctx context.Context, questionnaireIDs []int) ([]model.Targets, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = t.db
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	targets := make([]model.Targets, 0, len(questionnaireIDs))
	
	err = db.
		Table("targets").
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Scan(&targets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	return targets, nil
}
