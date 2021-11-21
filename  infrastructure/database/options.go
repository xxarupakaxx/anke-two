package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type Option struct {
	//TODO:使うかどうかはこれから
	infrastructure.SqlHandler
}

func (o *Option) InsertOption(ctx context.Context, lastID int, num int, body string) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	option := model.Options{
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
	panic("implement me")
}

func (o *Option) DeleteOptions(ctx context.Context, questionID int) error {
	panic("implement me")
}

func (o *Option) GetOptions(ctx context.Context, questionIDs []int) ([]model.Options, error) {
	panic("implement me")
}

func NewOption(sqlHandler infrastructure.SqlHandler) *Option {
	return &Option{SqlHandler: sqlHandler}
}
