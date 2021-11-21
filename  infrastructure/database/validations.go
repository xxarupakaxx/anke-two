package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type Validations struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func NewValidations(sqlHandler infrastructure.SqlHandler) *Validations {
	return &Validations{SqlHandler: sqlHandler}
}

func (v *Validations) InsertValidation(ctx context.Context, lastID int, validation model.Validations) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	validation.QuestionID = lastID
	err = db.Create(&validation).Error
	if err != nil {
		return fmt.Errorf("failed to insert the validation (lastID: %d): %w", lastID, err)
	}
	return nil
}

func (v *Validations) UpdateValidation(ctx context.Context, questionID int, validation model.Validations) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Model(&model.Validations{}).
		Where("question_id = ?", questionID).
		Updates(map[string]interface{}{
			"question_id":   questionID,
			"regex_pattern": validation.RegexPattern,
			"min_bound":     validation.MinBound,
			"max_bound":     validation.MaxBound,
		})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update the validation (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a vaidation record :%w", model.ErrNoRecordUpdated)
	}

	return nil
}

func (v *Validations) DeleteValidation(ctx context.Context, questionID int) error {
	panic("implement me")
}

func (v *Validations) GetValidations(ctx context.Context, qustionIDs []int) ([]model.Validations, error) {
	panic("implement me")
}

func (v *Validations) CheckNumberValidation(validation model.Validations, Body string) error {
	panic("implement me")
}

func (v *Validations) CheckTextValidation(validation model.Validations, Response string) error {
	panic("implement me")
}

func (v *Validations) CheckNumberValid(MinBound, MaxBound string) error {
	panic("implement me")
}
