package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
)

type Administrator struct {
	infrastructure.SqlHandler
}

func NewAdministrator(sqlHandler infrastructure.SqlHandler) *Administrator {
	return &Administrator{SqlHandler: sqlHandler}
}



func (a *Administrator) InsertAdministrator(ctx context.Context, questionnaireID int, administrators []string) error {
	db,err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	dbAdministrators := make([]model.Administrators,0,len(administrators))

	if len(administrators) == 0 {
		return nil
	}

	for _, administrator := range administrators {
		dbAdministrators = append(dbAdministrators,model.Administrators{
			QuestionnaireID: questionnaireID,
			UserTraqid:      administrator,
		})
	}

	err = db.Create(&dbAdministrators).Error
	if err != nil {
		return fmt.Errorf("failed to insert administrators: %w",err)
	}

	return nil
}

func (a *Administrator) DeleteAdministrators(ctx context.Context, questionnaireID int) error {
	db,err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w",err)
	}
	err = db.
		Where("questionnaire_id = ?",questionnaireID).
		Delete(&model.Administrators{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete administrators: %w",err)
	}
	return nil
}

func (a *Administrator) GetAdministrators(ctx context.Context, questionnaireIDs []int) ([]model.Administrators, error) {
	db,err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("                                   failed to get transaction:%w",err)
	}
	dbAdministrators := make([]model.Administrators,len(questionnaireIDs))
	err = db.
		Where("questionnaire_id IN <?>",questionnaireIDs).
		Find(&dbAdministrators).Error
	if err != nil {
		return nil,fmt.Errorf("failed to get administrators:%w",err)
	}

	return dbAdministrators,nil
}

func (a *Administrator) CheckQuestionnaireAdmin(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	panic("implement me")
}

