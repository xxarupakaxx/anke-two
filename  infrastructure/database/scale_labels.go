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
	panic("implement me")
}

func (s *ScaleLabel) DeleteScaleLabel(ctx context.Context, questionID int) error {
	panic("implement me")
}

func (s *ScaleLabel) GetScaleLabels(ctx context.Context, questionIDs []int) ([]model.ScaleLabels, error) {
	panic("implement me")
}

func (s *ScaleLabel) CheckScaleLabel(label model.ScaleLabels, response string) error {
	panic("implement me")
}
