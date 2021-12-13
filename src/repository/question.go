package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Question interface {
	// GetQuestions 各アンケートごとの質問を取得
	GetQuestions(ctx context.Context,questionnaireID int) ([]*model.Question,error)
	// CreateQuestion 質問を作成
	CreateQuestion(ctx context.Context,question *model.Question) (int, error)
	// DeleteQuestion 質問の削除
	DeleteQuestion(ctx context.Context,id int) error
	// UpdateQuestion 質問を更新
	UpdateQuestion(ctx context.Context,question *model.Question) error
}