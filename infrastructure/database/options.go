package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Option struct {
	db *gorm.DB
}

func NewOption(db *gorm.DB) *Option {
	return &Option{db: db}
}

func (o *Option) InsertOption(ctx context.Context, lastID int, num int, body string) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = o.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	option := Options{
		QuestionID: lastID,
		OptionNum:  num,
		Body:       body,
	}
	err = db.Create(&option).Error
	if err != nil {
		return fmt.Errorf("failed to insert a option: %w", err)
	}
	return nil
}

func (o *Option) UpdateOptions(ctx context.Context, options []string, questionID int) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = o.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}
	var previousOptions []Options
	err = db.
		Session(&gorm.Session{}).
		Where("question_id = ?", questionID).
		Select("option_num", "body").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&previousOptions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to get option: %w", err)
	}

	isDelete := false
	optionMap := make(map[int]*Options, len(options))
	for i, option := range previousOptions {
		if option.OptionNum <= len(options) {
			optionMap[option.OptionNum] = &previousOptions[i]
		} else {
			isDelete = true
		}
	}

	createOptions := []Options{}
	for i, optionLabel := range options {
		optionNum := i + 1

		if option, ok := optionMap[optionNum]; ok {
			if option.Body != optionLabel {
				err := db.
					Session(&gorm.Session{}).
					Model(&Options{}).
					Where("question_id = ?", questionID).
					Where("option_num = ?", optionNum).
					Update("body", optionLabel).Error
				if err != nil {
					return fmt.Errorf("failed to update option: %w", err)
				}
			}
		} else {
			createOptions = append(createOptions, Options{
				QuestionID: questionID,
				OptionNum:  optionNum,
				Body:       optionLabel,
			})
		}
	}

	if len(createOptions) > 0 {
		err := db.
			Session(&gorm.Session{}).
			Create(&createOptions).Error
		if err != nil {
			return fmt.Errorf("failed to create option: %w", err)
		}
	}

	if isDelete {
		err = db.
			Where("question_id = ? AND option_num > ?", questionID, len(options)).
			Delete(Options{}).Error
		if err != nil {
			return fmt.Errorf("failed to update option: %w", err)
		}
	}

	return nil
}

func (o *Option) DeleteOptions(ctx context.Context, questionID int) error {
	db, err := GetTx(ctx)
	if db == nil {
		db = o.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	err = db.
		Where("question_id = ?", questionID).
		Delete(Options{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete option: %w", err)
	}
	return nil
}

func (o *Option) GetOptions(ctx context.Context, questionIDs []int) ([]model.Options, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = o.db
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}
	dbOptions := make([]model.Options, len(questionIDs))

	err = db.
		Table("options").
		Where("question_id IN (?)", questionIDs).
		Order("question_id, option_num").
		Scan(&dbOptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get option: %w", err)
	}

	return dbOptions, nil
}
