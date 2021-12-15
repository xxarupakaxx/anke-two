package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *GormRepository) GetScaleLabels(ctx context.Context, questionIDs []int) ([]*model.ScaleLabel, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil,fmt.Errorf("failed to get db:%w", err)
	}

	labels := make([]*model.ScaleLabel,len(questionIDs))

	err = db.
		Where("question_id IN ?",questionIDs).
		Find(&labels).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get scalelabels :%w",err)
	}
	return labels,err
}

func (repo *GormRepository) CreateScaleLabel(ctx context.Context, label *model.ScaleLabel) error {
	panic("implement me")
}

func (repo *GormRepository) UpdateScaleLabel(ctx context.Context, label *model.ScaleLabel) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}
	var preScale *model.ScaleLabel

	err = db.
		Where("question_id", label.QuestionID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&preScale).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed get to preScaleLable :%w", err)
	}

	mapScale := compareChangeField(preScale, label)
	result := db.
		Session(&gorm.Session{}).
		Where("question_id = ?", label.QuestionID).
		Updates(&mapScale)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update scaleLabel :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordUpdated
	}

	return nil
}

func (repo *GormRepository) DeleteScaleLabel(ctx context.Context, questionID int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("question_id = ?", questionID).
		Delete(&model.ScaleLabel{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete scalelLabel :%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

func compareChangeField(preScale, scale *model.ScaleLabel) map[string]interface{} {
	mapScale := map[string]interface{}{}
	if preScale.ScaleLabelLeft != scale.ScaleLabelLeft {
		mapScale["scale_label_left"] = scale.ScaleLabelLeft
	}
	if preScale.ScaleLabelRight != scale.ScaleLabelRight {
		mapScale["scale_label_right"] = scale.ScaleLabelRight
	}
	if preScale.ScaleMax != scale.ScaleMax {
		mapScale["scale_max"] = scale.ScaleMin
	}
	if preScale.ScaleMin != scale.ScaleMin {
		mapScale["scale_min"] = scale.ScaleMin
	}

	return mapScale
}
