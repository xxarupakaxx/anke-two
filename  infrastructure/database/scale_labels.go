package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type ScaleLabel struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func NewScaleLabel(sqlHandler infrastructure.SqlHandler) *ScaleLabel {
	return &ScaleLabel{SqlHandler: sqlHandler}
}

func (s *ScaleLabel) InsertScaleLabel(ctx context.Context, lastID int, label model.ScaleLabels) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	label.QuestionID = lastID

	err = db.Create(&label).Error
	if err != nil {
		fmt.Errorf("failed to insert the scale label(lastID:%d): %w", lastID, err)
	}

	return nil
}

func (s *ScaleLabel) UpdateScaleLabel(ctx context.Context, questionID int, label model.ScaleLabels) error {
	db, err := GetTx(ctx)
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
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	result := db.
		Where("question_id = ?", questionID).
		Delete(&model.ScaleLabels{})
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
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}

	labels := make([]model.ScaleLabels, len(questionIDs))

	err = db.
		Where("question_id IN (?)", questionIDs).
		Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get scaleLable :%w", err)
	}

	return labels, nil
}

func (s *ScaleLabel) CheckScaleLabel(label model.ScaleLabels, response string) error {
	panic("implement me")
}
