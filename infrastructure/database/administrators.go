package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
)

type Administrator struct {
	db *gorm.DB
}

func NewAdministrator(db *gorm.DB) *Administrator {
	return &Administrator{db: db}
}

func (a *Administrator) InsertAdministrator(ctx context.Context, questionnaireID int, administrators []string) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = a.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	dbAdministrators := make([]Administrators, 0, len(administrators))

	if len(administrators) == 0 {
		return nil
	}

	for _, administrator := range administrators {
		dbAdministrators = append(dbAdministrators, Administrators{
			QuestionnaireID: questionnaireID,
			UserTraqid:      administrator,
		})
	}

	err = db.Create(&dbAdministrators).Error
	if err != nil {
		return fmt.Errorf("failed to insert administrators: %w", err)
	}

	return nil
}

func (a *Administrator) DeleteAdministrators(ctx context.Context, questionnaireID int) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = a.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Delete(&Administrators{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete administrators: %w", err)
	}
	return nil
}

func (a *Administrator) GetAdministrators(ctx context.Context, questionnaireIDs []int) ([]model.Administrators, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = a.db
	}
	if err != nil {
		return nil, fmt.Errorf("                                   failed to get transaction:%w", err)
	}
	dbAdministrators := make([]model.Administrators, len(questionnaireIDs))
	err = db.
		Table("administrators").
		Where("questionnaire_id IN <?>", questionnaireIDs).
		Find(&dbAdministrators).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get administrators:%w", err)
	}

	return dbAdministrators, nil
}

func (a *Administrator) CheckQuestionnaireAdmin(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = a.db
	}
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&model.Administrators{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get Administrators:%w", err)
	}
	return true, nil
}
