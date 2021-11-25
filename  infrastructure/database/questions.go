package database

import (
	"context"
	"errors"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
	"log"
)

type Question struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func NewQuestion(sqlHandler infrastructure.SqlHandler) *Question {
	err := setUpQuestionTypes(sqlHandler.Db)
	if err != nil {
		log.Fatalf("failed to get db: %w", err)
	}
	return &Question{SqlHandler: sqlHandler}
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
	db, err := GetTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction:%w", err)
	}

	var questionsType model.QuestionType
	err = db.
		Where("question_type = ?", questionType).
		Select("id").
		First(&questionsType).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get questionType :%w", err)
	}

	intQuestionType := questionsType.ID

	question := model.Questions{
		QuestionnaireID: questionnaireID,
		PageNum:         pageNum,
		QuestionNum:     questionNum,
		Type:            intQuestionType,
		Body:            body,
		IsRequired:      isRequired,
	}
	err = db.Create(&question).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert a question record: %w", err)
	}

	return question.ID, nil
}

func (q *Question) UpdateQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool, questionID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}

	var questionsType model.QuestionType
	err = db.
		Where("question_type = ?", questionType).
		Select("id").
		First(&questionsType).Error
	if err != nil {
		return fmt.Errorf("failed to get questionType :%w", err)
	}

	intQuestionType := questionsType.ID

	question := map[string]interface{}{
		"questionnaire_id": questionnaireID,
		"page_num":         pageNum,
		"question_num":     questionNum,
		"type":             intQuestionType,
		"body":             body,
		"is_required":      isRequired,
	}

	result := db.
		Where("question_id = ?", questionID).
		Updates(question)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update question:%w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("NO update question :%w", model.ErrNoRecordUpdated)
	}

	return nil
}

func (q *Question) DeleteQuestion(ctx context.Context, questionID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Where("id = ?", questionID).
		Delete(&model.Questions{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete question :%w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no delete question :%w", model.ErrNoRecordDeleted)
	}
	return nil
}

func (q *Question) GetQuestions(ctx context.Context, questionnaireID int) ([]model.ReturnQuestions, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	questions := make([]model.Questions, 0)

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	questionType, err := q.ChangeStrQuestionType(ctx, questions)
	if err != nil {
		return nil, fmt.Errorf("failed to change QuestionType :%w", err)
	}
	returnQuestions := []model.ReturnQuestions{}
	for _, question := range questions {
		for i, m := range questionType {
			if question.ID == i {
				returnQuestions = append(returnQuestions, model.ReturnQuestions{
					ID:              question.ID,
					QuestionnaireID: question.QuestionnaireID,
					PageNum:         question.PageNum,
					QuestionNum:     question.QuestionNum,
					Type:            m.QuestionType,
					Body:            question.Body,
					IsRequired:      question.IsRequired,
					DeletedAt:       question.DeletedAt,
					CreatedAt:       question.CreatedAt,
					Options:         question.Options,
					Responses:       question.Responses,
					ScaleLabels:     question.ScaleLabels,
					Validations:     question.Validations,
				})
			}
		}
	}
	return returnQuestions, nil
}

func (q *Question) CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction :%w", err)
	}

	err = db.
		Joins("INNER JOIN administrators ON question.questionnaire_id = administrators.questionnaire_id").
		Where("question.id = ? AND administrators.user_traqid = ?", questionID, userID).
		Select("question.id").
		First(&model.Questions{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get question_id: %w", err)
	}

	return true, nil
}

func (q *Question) ChangeStrQuestionType(ctx context.Context, questions []model.Questions) (map[int]model.QuestionType, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}

	questionsIntType := map[int]int{}
	for _, question := range questions {
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Table("questions").
			Where("id = ?", question.ID).
			Pluck("type", questionsIntType[question.ID]).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get questionType in Questions Table:%w", err)
		}
	}

	questionsType := map[int]model.QuestionType{}

	for questionID, questionType := range questionsIntType {
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("id = ?", questionType).
			Find(questionsType[questionID]).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get questionType in QuestionType Table :%w", err)
		}
	}

	return questionsType, nil
}
