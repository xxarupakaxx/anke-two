package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
)

type Question struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func setUpQuestionTypes(db *gorm.DB) error {
	questionTypes := []model.QuestionType{
		{
			QuestionType: "Text",
		},
		{
			QuestionType: "TextArea",
		},
		{
			QuestionType: "Number",
		},
		{
			QuestionType: "MultipleChoice",
		},
		{
			QuestionType: "Checkbox",
		},
		{
			QuestionType: "Dropdown",
		},
		{
			QuestionType: "LinearScale",
		},
		{
			QuestionType: "Date",
		},
		{
			QuestionType: "Time",
		},
	}

	for _, questionType := range questionTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("question_type = ?", questionType.QuestionType).
			FirstOrCreate(&questionType).Error
		if err != nil {
			return fmt.Errorf("failed to create QuestionType:%w", err)
		}
	}

	return nil
}
func (q *Question) InsertQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool) (int, error) {

}

func (q *Question) UpdateQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool, questionID int) error {
	panic("implement me")
}

func (q *Question) DeleteQuestion(ctx context.Context, questionID int) error {
	panic("implement me")
}

func (q *Question) GetQuestions(ctx context.Context, questionnaireID int) ([]model.Questions, error) {
	panic("implement me")
}

func (q *Question) CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error) {
	panic("implement me")
}
