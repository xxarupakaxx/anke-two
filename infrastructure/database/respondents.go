package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"math"
	"sort"
	"strconv"
	"time"
)

type Respondent struct {
}

func NewRespondent() *Respondent {
	return &Respondent{}
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
		Find(&questions).Error
	if err != nil {
		return model.RespondentDetail{}, fmt.Errorf("failed to get respondent : %w", err)
	}

	questionsIntType := map[int]int{}
	for _, question := range questions {
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Table("questions").
			Where("id = ?", question.ID).
			Pluck("type", questionsIntType[question.ID]).Error
		if err != nil {
			return model.RespondentDetail{}, fmt.Errorf("failed to get questionType in Questions Table:%w", err)
		}
	}

	questionsType := map[int]model.QuestionType{}

	for questionID, questionType := range questionsIntType {
		err = db.
			Session(&gorm.Session{NewDB: true}).
			Where("id = ?", questionType).
			Find(questionsType[questionID]).Error
		if err != nil {
			return model.RespondentDetail{}, fmt.Errorf("failed to get questionType in Name Table :%w", err)
		}
	}

	respondentDetail := model.RespondentDetail{
		ResponseID:      responseID,
		TraqID:          respondent.UserTraqid,
		QuestionnaireID: respondent.QuestionnaireID,
		SubmittedAt:     respondent.SubmittedAt,
		UpdatedAt:       respondent.UpdatedAt,
	}

	questionsTypeName := []model.QuestionIDAndQuestionType{}

	for _, question := range questions {
		for id, questionType := range questionsType {
			if question.ID == id {
				questionsTypeName = append(questionsTypeName, model.QuestionIDAndQuestionType{
					QuestionID:   question.ID,
					QuestionType: questionType.Name,
					Responses:    question.Responses,
				})
			}
		}
	}

	for _, questionTypeName := range questionsTypeName {
		responseBody := model.ResponseBody{
			QuestionID:   questionTypeName.QuestionID,
			QuestionType: questionTypeName.QuestionType,
		}

		switch questionTypeName.QuestionType {
		case "MultipleChoice", "Checkbox", "Dropdown":
			for _, response := range questionTypeName.Responses {
				responseBody.OptionResponse = append(responseBody.OptionResponse, response.Body.String)
			}
		default:
			if len(questionTypeName.Responses) == 0 {
				responseBody.Body = null.NewString("", false)
			} else {
				responseBody.Body = questionTypeName.Responses[0].Body
			}
		}
		respondentDetail.Responses = append(respondentDetail.Responses, responseBody)
	}

	return respondentDetail, nil
}

func (r *Respondent) GetRespondentDetails(ctx context.Context, questionnaireID int, sort string) ([]model.RespondentDetail, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	respondents := make([]model.Respondents, 0)

	query := db.
		Session(&gorm.Session{}).
		Where("respondents.questionnaire_id = ? AND respondents.submitted_at IS NOT NULL", questionnaireID).
		Select("ResponseID", "UserTraqid", "UpdatedAt", "SubmittedAt")

	query, sortNum, err := setRespondentsOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set order :%w", err)
	}

	err = query.
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents:%w", err)
	}

	if len(respondents) == 0 {
		return []model.RespondentDetail{}, nil
	}

	responseIDs := make([]int, 0, len(respondents))
	for _, respondent := range respondents {
		responseIDs = append(responseIDs, respondent.ResponseID)
	}

	respondentDetails := make([]model.RespondentDetail, 0, len(respondents))
	respondentDetailMap := make(map[int]*model.RespondentDetail, len(respondents))
	for i, respondent := range respondents {
		respondentDetails = append(respondentDetails, model.RespondentDetail{
			ResponseID:      respondent.ResponseID,
			TraqID:          respondent.UserTraqid,
			QuestionnaireID: questionnaireID,
			SubmittedAt:     respondent.SubmittedAt,
			UpdatedAt:       respondent.UpdatedAt,
		})

		respondentDetailMap[respondent.ResponseID] = &respondentDetails[i]
	}

	questions := make([]model.Questions, len(respondents))

	err = db.
		Preload("Responses", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("ResponseID", "QuestionID", "Body").
				Where("response_id IN (?)", responseIDs)
		}).Where("questionnaire_id = ?", questionnaireID).
		Order("question_num").
		Select("ID", "Type").
		Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get questions:%w", err)
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

	questionsTypeName := []model.QuestionIDAndQuestionType{}

	for _, question := range questions {
		for id, questionType := range questionsType {
			if question.ID == id {
				questionsTypeName = append(questionsTypeName, model.QuestionIDAndQuestionType{
					QuestionID:   question.ID,
					QuestionType: questionType.Name,
					Responses:    question.Responses,
				})
			}
		}
	}

	for _, questionTypeName := range questionsTypeName {
		responseBodyMap := make(map[int][]string, len(respondents))
		for _, responses := range questionTypeName.Responses {
			if responses.Body.Valid {
				responseBodyMap[responses.ResponseID] = append(responseBodyMap[responses.ResponseID], responses.Body.String)
			}
		}

		for i := range respondentDetails {
			responseBodies := responseBodyMap[respondentDetails[i].ResponseID]
			responseBody := model.ResponseBody{
				QuestionID:   questionTypeName.QuestionID,
				QuestionType: questionTypeName.QuestionType,
			}

			switch responseBody.QuestionType {
			case "MultipleChoice", "Checkbox", "Dropdown":
				if responseBodies == nil {
					responseBody.OptionResponse = []string{}
				} else {
					responseBody.OptionResponse = responseBodies
				}
			default:
				if len(responseBodies) == 0 {
					responseBody.Body = null.NewString("", false)
				} else {
					responseBody.Body = null.NewString(responseBodies[0], true)
				}
			}

			respondentDetails[i].Responses = append(respondentDetails[i].Responses, responseBody)
		}
	}
	respondentDetails, err = sortRespondentDetail(sortNum, len(questions), respondentDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to sort RespondentDetails: %w", err)
	}
	return respondentDetails, err

}

func (r *Respondent) GetRespondentsUserIDs(ctx context.Context, questionnaireIDs []int) ([]model.Respondents, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}
	respondents := make([]model.Respondents, len(questionnaireIDs))

	err = db.
		Where("questionnaire_id IN (?)", questionnaireIDs).
		Select("questionnaire_id,user_traqid").
		Find(&respondents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get respondents%w", err)
	}

	return respondents, nil
}

func (r *Respondent) CheckRespondent(ctx context.Context, userID string, questionnaireID int) (bool, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction :%w", err)
	}
	err = db.
		Where("user_traqid = ? AND questionnaire_id = ?", userID, questionnaireID).
		First(&model.Respondents{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get response:%w", err)
	}

	return true, nil
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

func setRespondentsOrder(query *gorm.DB, sort string) (*gorm.DB, int, error) {
	var sortNum int
	switch sort {
	case "traqid":
		query = query.Order("user_traqid")
	case "-traqid":
		query = query.Order("user_traqid DESC")
	case "submitted_at":
		query = query.Order("submitted_at")
	case "-submitted_at":
		query = query.Order("submitted_at DESC")
	case "":
	default:
		var err error
		sortNum, err = strconv.Atoi(sort)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert sort param to int: %w", err)
		}
	}

	query = query.Order("response_id")

	return query, sortNum, nil
}
