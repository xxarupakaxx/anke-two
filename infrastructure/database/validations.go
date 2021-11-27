package database

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"regexp"
	"strconv"
)

type Validation struct{}

func NewValidations() *Validation {
	return &Validation{}
}

func (v *Validation) InsertValidation(ctx context.Context, lastID int, validation model.Validations) error {
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

func (v *Validation) UpdateValidation(ctx context.Context, questionID int, validation model.Validations) error {
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

func (v *Validation) DeleteValidation(ctx context.Context, questionID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Where("question_id =?", questionID).
		Delete(&model.Validations{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete validation(questionID  :%d), : %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete a validation :%w", model.ErrNoRecordDeleted)
	}

	return nil
}

func (v *Validation) GetValidations(ctx context.Context, questionIDs []int) ([]model.Validations, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	validations := make([]model.Validations, len(questionIDs))

	err = db.
		Where("question_id IN (?)", questionIDs).
		Find(&validations).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get the validations : %w", err)
	}

	return validations, nil
}

func (v *Validation) CheckNumberValidation(validation model.Validations, Body string) error {
	if err := v.CheckNumberValid(validation.MinBound, validation.MaxBound); err != nil {
		return err
	}

	if Body == "" {
		return nil
	}

	number, err := strconv.ParseFloat(Body, 64)
	if err != nil {
		return model.ErrInvalidNumber
	}

	if validation.MaxBound != "" {
		maxBoundNum, _ := strconv.ParseFloat(validation.MaxBound, 64)
		if maxBoundNum < number {
			return fmt.Errorf("failed to meet the boundary value. the number must be greater than MinBound (number: %g, MinBound: %g): %w", number, maxBoundNum, model.ErrNumberBoundary)
		}
	}
	if validation.MinBound != "" {
		minBoundNum, _ := strconv.ParseFloat(validation.MinBound, 64)
		if minBoundNum > number {
			return fmt.Errorf("failed to meet the boundary value. the number must be greater than MinBound (number: %g, MinBound: %g): %w", number, minBoundNum, model.ErrNumberBoundary)
		}
	}
	return nil
}

func (v *Validation) CheckTextValidation(validation model.Validations, Response string) error {
	r, err := regexp.Compile(validation.RegexPattern)
	if err != nil {
		return fmt.Errorf("failed to compile the pattern (RegexPattern: %s): %w", r, model.ErrInvalidRegex)
	}

	if !r.MatchString(Response) && Response != "" {
		return fmt.Errorf("failed to match the pattern (Response: %s),RegexPattern: %s : %w", Response, r, model.ErrTextMatching)
	}

	return nil
}

func (v *Validation) CheckNumberValid(MinBound, MaxBound string) error {
	var minBoundNum, maxBoundNum float64
	if MinBound != "" {
		min, err := strconv.ParseFloat(MinBound, 64)
		minBoundNum = min
		if err != nil {
			return fmt.Errorf("failed to check the boundary value. MinBound is not a numerical value: %w", model.ErrInvalidNumber)
		}
	}
	if MaxBound != "" {
		max, err := strconv.ParseFloat(MaxBound, 64)
		maxBoundNum = max
		if err != nil {
			return fmt.Errorf("failed to check the boundary value. MaxBound is not a numerical value: %w", model.ErrInvalidNumber)
		}
	}

	if MinBound != "" && MaxBound != "" {
		if minBoundNum > maxBoundNum {
			return fmt.Errorf("failed to check the boundary value. MinBound must be less than MaxBound (MinBound: %g, MaxBound: %g): %w", minBoundNum, maxBoundNum, model.ErrInvalidNumber)
		}
	}

	return nil
}
