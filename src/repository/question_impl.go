package repository

import (
	"context"
	"fmt"
	"github.com/xxarupkaxx/anke-two/src/model"
	"gorm.io/gorm"
)

func (repo *GormRepository) GetQuestions(ctx context.Context, questionnaireID int) ([]*model.Question, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get db:%w", err)
	}

	questions := make([]*model.Question, 0)

	err = db.
		Where("questionnaire_id = ?", questionnaireID).
		Preload("QuestionType").
		Order("question_num").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions :%w", err)
	}

	return questions, nil
}

func (repo *GormRepository) CreateQuestion(ctx context.Context, question *model.Question) (int, error) {
	db, err := repo.getDB(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get db:%w", err)
	}

	err = db.Create(&question).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create question:%w", err)
	}

	return question.ID, err
}

func (repo *GormRepository) DeleteQuestion(ctx context.Context, id int) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Where("id = ?", id).
		Delete(&model.Question{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete question:%w", err)
	}
	if result.RowsAffected == 0 {
		return ErrNoRecordDeleted
	}

	return nil
}

func (repo *GormRepository) UpdateQuestion(ctx context.Context, question *model.Question) error {
	db, err := repo.getDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db:%w", err)
	}

	result := db.
		Model(&model.Question{}).
		Where("id =?", question.ID).
		Updates(&question)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update question:%w", err)
	}

	return nil
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