package database

import (
	"context"
	"errors"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
	"time"
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

func (r *Respondent) UpdateSubmittedAt(ctx context.Context, responseID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	result := db.
		Model(&model.Respondents{}).
		Where("response_id = ?", responseID).
		Update("submitted_at", time.Now())
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to update submittedAt :%w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update no data :%w", model.ErrNoRecordUpdated)
	}

	return nil
}

func (r *Respondent) DeleteRespondent(ctx context.Context, responseID int) error {
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction :%w", err)
	}

	result := db.
		Where("response_id = ?", responseID).
		Delete(&model.Respondents{})
	err = result.Error
	if err != nil {
		return fmt.Errorf("failed to delete respondent :%w", err)
	}
	if result.RowsAffected == 0 {
		return model.ErrNoRecordDeleted
	}

	return nil
}

func (r *Respondent) GetRespondent(ctx context.Context, responseID int) (*model.Respondents, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}
	var respondent model.Respondents

	err = db.
		Where("response_id = ?", responseID).
		First(&respondent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, model.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	return &respondent, nil
}

func (r *Respondent) GetRespondentInfos(ctx context.Context, userID string, questionnaireIDs ...int) ([]model.RespondentInfo, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	respondents := make([]model.RespondentInfo, len(questionnaireIDs))

	query := db.
		Table("respondents").
		Joins("LEFT OUTER JOIN questionnaires ON respondents.questionnaire_id = questionnaires.id").
		Order("respondents.submitted_at DESC").
		Where("user_traqid = ? AND respondents.deleted_at IS NULL AND questionnaires.deleted_at IS NULL", userID)

	if len(questionnaireIDs) != 0 {
		questionnaireID := questionnaireIDs[0]
		query = query.Where("questionnaire_id = ?", questionnaireID)
	} else if len(questionnaireIDs) > 1 {
		return nil, errors.New("illegal function usage")
	}

	err = query.
		Select("respondents.questionnaire_id, respondents.response_id, respondents.modified_at, respondents.submitted_at, questionnaires.title, questionnaires.res_time_limit").
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get my responses: %w", err)
	}

	return respondents, nil
}

func (r *Respondent) GetRespondentDetail(ctx context.Context, responseID int) (model.RespondentDetail, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return model.RespondentDetail{}, fmt.Errorf("failed to get transaction:%w", err)
	}

	respondent := model.Respondents{}

	err = db.
		Session(&gorm.Session{}).
		Where("respondents.response_id = ?", responseID).
		Select("QuestionnaireID", "UserTraqid", "ModifiedAt", "SubmittedAt").
		Take(&respondent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.RespondentDetail{}, model.ErrRecordNotFound
	}
	if err != nil {
		return model.RespondentDetail{}, fmt.Errorf("failed to get respondent: %w", err)
	}

	questions := make([]model.Questions, 0)

	err = db.
		Where("questionnaire_id = ?", respondent.QuestionnaireID).
		Preload("Responses", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("QuestionID", "Body").
				Where("response_id = ?", responseID)
		}).Select("ID", "Type").
		Order("ID").
		Find(&questions).Error
	if err != nil {
		return model.RespondentDetail{}, fmt.Errorf("failed to get respondent : %w", err)
	}
	questionsID := []int{}
	for _, question := range questions {
		questionsID = append(questionsID, question.ID)
	}
	questionsType := make([]model.QuestionType, 0)

	err = db.
		Where("ID IN (?)", questionsID).
		Order("ID").
		Find(&questionsType).Error
	if err != nil {
		return model.RespondentDetail{}, fmt.Errorf("failed to get questionsType:%w", err)
	}

	respondentDetail := model.RespondentDetail{
		ResponseID:      responseID,
		TraqID:          respondent.UserTraqid,
		QuestionnaireID: respondent.QuestionnaireID,
		SubmittedAt:     respondent.SubmittedAt,
		UpdatedAt:       respondent.UpdatedAt,
	}

	type questionIDAndQuestionType struct {
		QuestionID   int
		QuestionType string
		Responses    []model.Responses
	}
	questionsTypeName := []questionIDAndQuestionType{}

	for _, question := range questions {
		for _, questionType := range questionsType {
			if questionType.ID == question.Type {
				questionsTypeName = append(questionsTypeName, questionIDAndQuestionType{
					QuestionID:   question.ID,
					QuestionType: questionType.QuestionType,
				})
			}
		}

	}

	for _, question := range questionsTypeName {
		responseBody := model.ResponseBody{
			QuestionID:   question.QuestionID,
			QuestionType: question.QuestionType,
		}

		switch question.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, response := range question.Responses {
				responseBody.OptionResponse = append(responseBody.OptionResponse, response.Body.String)
			}
		default:
			if len(question.Responses) == 0 {
				responseBody.Body = null.NewString("", false)
			} else {
				responseBody.Body = question.Responses[0].Body
			}
		}
		respondentDetail.Responses = append(respondentDetail.Responses, responseBody)
	}

	return respondentDetail, nil
}

func (r *Respondent) GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]model.RespondentDetail, error) {

}

func (r *Respondent) GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]model.Respondents, error) {
	panic("implement me")
}

func (r *Respondent) CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	panic("implement me")
}

func sortRespondentDetail(sortNum, questionNum int, respondentDetails []model.RespondentDetail) ([]model.RespondentDetail, error) {
	if sortNum == 0 {
		return respondentDetails, nil
	}

	sortNumAbs := int(math.Abs(float64(sortNum)))
	if sortNumAbs > questionNum {
		return nil, fmt.Errorf("sort param is too large:%d", sortNum)
	}

	sort.Slice(respondentDetails, func(i, j int) bool {
		bodyI := respondentDetails[i].Responses[sortNumAbs-1]
		bodyJ := respondentDetails[j].Responses[sortNumAbs-1]
		if bodyI.QuestionType == "Number" {
			numi, err := strconv.ParseFloat(bodyI.Body.String, 64)
			if err != nil {
				return true
			}
			numj, err := strconv.ParseFloat(bodyJ.Body.String, 64)
			if err != nil {
				return true
			}
			if sortNum < 0 {
				return numi > numj
			}
			return numi < numj
		}
		if sortNum < 0 {
			return bodyI.Body.String > bodyJ.Body.String
		}
		return bodyI.Body.String < bodyJ.Body.String
	})
	return respondentDetails, nil
}
