package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"log"
)

type Questionnaire struct {
	//TODO:後で考える
	infrastructure.SqlHandler
}

func NewQuestionnaire(sqlHandler infrastructure.SqlHandler) *Questionnaire {
	err := setUpResSharedTo(sqlHandler.Db)
	if err != nil {
		log.Fatalf("failed to get db:%w", err)
	}
	return &Questionnaire{SqlHandler: sqlHandler}
}

func setUpResSharedTo(db *gorm.DB) error {
	resSharedTypes := []model.ResShareTypes{
		{
			Name: "administrators",
		},
		{
			Name: "respondents",
		},
		{
			Name: "public",
		},
	}
	for _, resSharedType := range resSharedTypes {
		err := db.
			Session(&gorm.Session{}).
			Where("name = ?", resSharedType.Name).
			FirstOrCreate(&resSharedType).Error
		if err != nil {
			return fmt.Errorf("failed to create resSharedType:%w", err)
		}
	}

	return nil
}

func (q *Questionnaire) InsertQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string) (int, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction:%w", err)
	}

	resSharedToType := model.ResShareTypes{}

	err = db.
		Where("name = ?", resSharedTo).
		First(&resSharedToType).
		Select("id").Error
	if err != nil {
		return 0, fmt.Errorf("failed to get resSharedToType :%w", err)
	}
	intResSharedTo := resSharedToType.ID

	var questionnaire model.Questionnaires
	if !resTimeLimit.Valid {
		questionnaire = model.Questionnaires{
			Title:       title,
			Description: description,
			ResSharedTo: intResSharedTo,
		}
	} else {
		questionnaire = model.Questionnaires{
			Title:        title,
			Description:  description,
			ResTimeLimit: resTimeLimit,
			ResSharedTo:  intResSharedTo,
		}
	}

	err = db.Create(&questionnaire).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert questionnaire:%w", err)
	}

	return questionnaire.ID, nil
}

func (q *Questionnaire) UpdateQuestionnaire(ctx context.Context, title string, description string, resTimeLimit null.Time, resSharedTo string, questionnaireID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	resSharedToType := model.ResShareTypes{}

	err = db.
		Where("name = ?", resSharedTo).
		First(&resSharedToType).
		Select("id").Error
	if err != nil {
		return fmt.Errorf("failed to get resSharedToType :%w", err)
	}
	intResSharedTo := resSharedToType.ID

	var questionnaire interface{}
	if resTimeLimit.Valid {
		questionnaire = model.Questionnaires{
			Title:        title,
			Description:  description,
			ResTimeLimit: resTimeLimit,
			ResSharedTo:  intResSharedTo,
		}
	} else {
		questionnaire = map[string]interface{}{
			"title":          title,
			"description":    description,
			"res_time_limit": gorm.Expr("NULL"),
			"res_shared_to":  resSharedTo,
		}
	}

	result := db.
		Model(&model.Questionnaires{}).
		Where("id = ?", questionnaireID).
		Updates(questionnaire)
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update questionnaire %w", err)
	}
	if result.RowsAffected == 0 {
		return model.ErrNoRecordUpdated
	}

	return nil
}

func (q *Questionnaire) DeleteQuestionnaire(ctx context.Context, questionnaireID int) error {
	panic("implement me")
}

func (q *Questionnaire) GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, nonTargeted bool) ([]model.QuestionnaireInfo, error) {
	panic("implement me")
}

func (q *Questionnaire) GetAdminQuestionnaires(ctx context.Context, userID string) ([]model.Questionnaires, error) {
	panic("implement me")
}

func (q *Questionnaire) GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (model.QuestionnaireInfo, error) {
	panic("implement me")
}

func (q *Questionnaire) GetTargetedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]model.TargetedQuestionnaire, error) {
	panic("implement me")
}

func (q *Questionnaire) GetQuestionnaireLimit(ctx context.Context, questionnaireID int) (null.Time, error) {
	panic("implement me")
}

func (q *Questionnaire) GetQuestionnaireLimitByResponseID(ctx context.Context, responseID int) (null.Time, error) {
	panic("implement me")
}

func (q *Questionnaire) GetResponseReadPrivilegeInfoByResponseID(ctx context.Context, userID string, responseID int) (*model.ResponseReadPrivilegeInfo, error) {
	panic("implement me")
}

func (q *Questionnaire) GetResponseReadPrivilegeInfoByQuestionnaireID(ctx context.Context, userID string, questionnaireID int) (*model.ResponseReadPrivilegeInfo, error) {
	panic("implement me")
}
