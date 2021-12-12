package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type IQuestion interface {
	InsertQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool) (int, error)
	UpdateQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool, questionID int) error
	DeleteQuestion(ctx context.Context, questionID int) error
	GetQuestions(ctx context.Context, questionnaireID int) ([]model.ReturnQuestions, error)
	CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error)
}
