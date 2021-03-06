package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gorm.io/gorm"
	"log"
)

type Question struct {
	db *gorm.DB
}

func NewQuestion(db *gorm.DB) *Question {
	err := setUpQuestionTypes(db)
	if err != nil {
		log.Fatalf("failed to get db: %w", err)
	}

	return &Question{db: db}
}

func setUpQuestionTypes(db *gorm.DB) error {
	questionTypes := []QuestionType{
		{
			Name: "Text",
		},
		{
			Name: "TextArea",
		},
		{
			Name: "Number",
		},
		{
			Name: "MultipleChoice",
		},
		{
			Name: "Checkbox",
		},
		{
			Name: "Dropdown",
		},
		{
			Name: "LinearScale",
		},
		{
			Name: "Date",
		},
		{
			Name: "Time",
		},
	}

	for _, questionType := range questionTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", questionType.Name).
			FirstOrCreate(&questionType).Error
		if err != nil {
			return fmt.Errorf("failed to create Name:%w", err)
		}
	}

	return nil
}
func (q *Question) InsertQuestion(ctx context.Context, questionnaireID int, pageNum int, questionNum int, questionType string, body string, isRequired bool) (int, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = q.db
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction:%w", err)
	}

	var questionsType QuestionType
	err = db.
		Where("question_type = ?", questionType).
		Select("id").
		First(&questionsType).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get questionType :%w", err)
	}

	intQuestionType := questionsType.ID

	question := Questions{
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
	if db == nil {
		db = q.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}

	var qType QuestionType
	err = db.
		Where("name = ?", questionType).
		Select("id").
		First(&questionType).Error
	if err != nil {
		return fmt.Errorf("failed to get questionType :%w", err)
	}

	intQuestionType := qType.ID

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
	if db == nil {
		db = q.db
	}
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Where("id = ?", questionID).
		Delete(&Questions{})
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
	if db == nil {
		db = q.db
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	questions := make([]model.ReturnQuestions, 0)

	err = db.
		Joins("INNER JOIN question_types ON question.type = question_types.id").
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Select("question.id, question.questionnaireID,question.page_num , question.question_num, question_types.name AS type, question.body,question.is_required, question.deleted_at , questions.created_at, question.options,question.responses, question.scale_labels, question.validations").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
}

func (q *Question) CheckQuestionAdmin(ctx context.Context, userID string, questionID int) (bool, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = q.db
	}
	if err != nil {
		return false, fmt.Errorf("failed to get transaction :%w", err)
	}

	err = db.
		Joins("INNER JOIN administrators ON question.questionnaire_id = administrators.questionnaire_id").
		Where("question.id = ? AND administrators.user_traqid = ?", questionID, userID).
		Select("question.id").
		First(&Questions{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get question_id: %w", err)
	}

	return true, nil
}

func (q *Question) ChangeStrQuestionType(ctx context.Context, questions []Questions) (map[int]model.QuestionType, error) {
	db, err := GetTx(ctx)
	if db == nil {
		db = q.db
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}

	questionsIntType := map[int]int{}
	for _, question := range questions {
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Table("question").
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
			return nil, fmt.Errorf("failed to get questionType in Name Table :%w", err)
		}
	}

	return questionsType, nil
}
