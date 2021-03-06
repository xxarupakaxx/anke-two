package database

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
	"strconv"
)

type ScaleLabel struct {
	db *gorm.DB
}

func NewScaleLabel(db *gorm.DB) *ScaleLabel {
	return &ScaleLabel{db: db}
}

func (s *ScaleLabel) InsertScaleLabel(ctx context.Context, lastID int, label model.ScaleLabels) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = s.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}

	l := ScaleLabels{
		QuestionID:      lastID,
		ScaleLabelRight: label.ScaleLabelRight,
		ScaleLabelLeft:  label.ScaleLabelLeft,
		ScaleMin:        label.ScaleMin,
		ScaleMax:        label.ScaleMax,
	}

	err = db.Create(&l).Error
	if err != nil {
		fmt.Errorf("failed to insert the scale label(lastID:%d): %w", lastID, err)
	}

	return nil
}

func (s *ScaleLabel) UpdateScaleLabel(ctx context.Context, questionID int, label model.ScaleLabels) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = s.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	result := db.
		Model(&ScaleLabel{}).
		Where("question_id", questionID).
		Updates(map[string]interface{}{
			"question_id":       questionID,
			"scale_label_right": label.ScaleLabelRight,
			"scale_label_left":  label.ScaleLabelLeft,
			"scale_min":         label.ScaleMin,
			"scale_max":         label.ScaleMax,
		})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update the scale label (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update a scale label record:%w", model.ErrNoRecordUpdated)
	}

	return nil
}

func (s *ScaleLabel) DeleteScaleLabel(ctx context.Context, questionID int) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = s.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	result := db.
		Where("question_id = ?", questionID).
		Delete(&ScaleLabels{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete the scale label (questionID: %d): %w", questionID, err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete a scale label : %w", model.ErrNoRecordDeleted)
	}

	return nil
}

func (s *ScaleLabel) GetScaleLabels(ctx context.Context, questionIDs []int) ([]model.ScaleLabels, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = s.db
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}

	labels := make([]model.ScaleLabels, len(questionIDs))

	err = db.
		Table("scale_labels").
		Where("question_id IN (?)", questionIDs).
		Scan(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get scaleLable :%w", err)
	}

	return labels, nil
}

func (s *ScaleLabel) CheckScaleLabel(label model.ScaleLabels, body string) error {
	if body == "" {
		return nil
	}

	r, err := strconv.Atoi(body)
	if err != nil {
		return err
	}
	if r < label.ScaleMin {
		return fmt.Errorf("failed to meet the scale. the response must be greater than ScaleMin (number: %d, ScaleMin: %d)", r, label.ScaleMin)
	} else if r > label.ScaleMax {
		return fmt.Errorf("failed to meet the scale. the response must be less than ScaleMax (number: %d, ScaleMax: %d)", r, label.ScaleMax)
	}

	return nil
}
