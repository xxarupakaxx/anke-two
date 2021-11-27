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
	Conn *gorm.DB
}

func NewQuestion(conn *gorm.DB) *Question {
	err := setUpQuestionTypes(conn)
	if err != nil {
		log.Fatalf("failed to get db: %w", err)
	}

	return &Question{Conn: conn}
}

func setUpQuestionTypes(db *gorm.DB) error {
	questionTypes := []model.QuestionType{
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

	var qType model.QuestionType
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

	questions := make([]model.ReturnQuestions, 0)

	err = db.
		Joins("INNER JOIN question_type ON questions.type = question_type.id").
		Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Select("questions.id, questions.questionnaireID,questions.page_num , questions.question_num, question_type.name, questions.body,questions.is_required, questions.deleted_at , questions.created_at, questions.options,questions.responses, questions.scale_labels, questions.validations").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
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
			return nil, fmt.Errorf("failed to get questionType in Name Table :%w", err)
		}
	}

	return questionsType, nil
}
