package repository

import (
	"context"
	"github.com/xxarupkaxx/anke-two/src/model"
)

type Questionnaire interface {
	GetQuestionnaire(ctx context.Context,id int) (*model.Questionnaire,error)
	CreateQuestionnaire(ctx context.Context,questionnaire *model.Questionnaire)(int,error)
	UpdateQuestionnaire(ctx context.Context,questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(ctx context.Context,id int) error
}
