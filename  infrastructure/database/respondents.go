package database

import (
	"context"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
)

type Respondent struct {
	//TODO:使うかどうかはあとから考える
	infrastructure.SqlHandler
}

func NewRespondent(sqlHandler infrastructure.SqlHandler) *Respondent {
	return &Respondent{SqlHandler: sqlHandler}
}

func (r *Respondent) InsertRespondent(ctx context.Context, userID string, questionnaireID int, submittedAt null.Time) (int, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction:%w", err)
	}
	var respondent model.Respondents
	if submittedAt.Valid {
		respondent = model.Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
			SubmittedAt:     submittedAt,
		}
	} else {
		respondent = model.Respondents{
			QuestionnaireID: questionnaireID,
			UserTraqid:      userID,
		}
	}
	err = db.Create(&respondent).Error
	if err != nil {
		return 0, fmt.Errorf("failed to insert respondent:%w", err)
	}

	return respondent.ResponseID, err
}

func (r *Respondent) DeleteRespondent(ctx context.Context, responseID int) error {
	panic("implement me")
}

func (r *Respondent) GetRespondent(ctx context.Context, responseID int) (*model.Respondents, error) {
	panic("implement me")
}

func (r *Respondent) GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]model.RespondentInfo, error) {
	panic("implement me")
}

func (r *Respondent) GetRespondentDetail(ctx context.Context, responseID int) (model.RespondentDetail, error) {
	panic("implement me")
}

func (r *Respondent) GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]model.RespondentDetail, error) {
	panic("implement me")
}

func (r *Respondent) GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]model.Respondents, error) {
	panic("implement me")
}

func (r *Respondent) CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	panic("implement me")
}
