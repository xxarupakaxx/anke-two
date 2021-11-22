package database

import (
	"context"
	"errors"
	"fmt"
	infrastructure "github.com/xxarupkaxx/anke-two/ infrastructure"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"log"
	"regexp"
	"time"
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
	db, err := GetTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction:%w", err)
	}
	result := db.Delete(&model.Questionnaires{ID: questionnaireID})
	err = result.Error

	if err != nil {
		return fmt.Errorf("failed to delete questionnaire: %w", err)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete questionnaire: %w", model.ErrNoRecordDeleted)
	}

	return nil
}

func (q *Questionnaire) GetQuestionnaires(ctx context.Context, userID string, sort string, search string, pageNum int, nonTargeted bool) ([]model.QuestionnaireInfo, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	db, err := GetTx(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transaction :%w", err)
	}

	questionnaireInfoes := make([]model.QuestionnaireInfo, 0, 20)

	query := db.
		Table("questionnaires").
		Joins("LEFT OUTER JOIN targets ON questionnaires.id = targets.questionnaire_id")

	query, err = setQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to set the order of the questionnaire table :%w", err)
	}

	if nonTargeted {
		query = query.Where("targets.questionnaire_id IS NULL OR (targets.user_traqid != ? AND targets.user_traqid != 'traP')", userID)
	}
	if len(search) != 0 {
		_, err := regexp.Compile(search)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid search param: %w", model.ErrInvalidRegex)
		}

		query = query.Where("questionnaires.title REGEXP ?", search)
	}

	var count int64
	err = query.
		Session(&gorm.Session{}).
		Group("questionnaires.id").
		Count(&count).Error
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, 0, model.ErrDeadlineExceeded
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve the number of questionnaires: %w", err)
	}

	if count == 0 {
		return []model.QuestionnaireInfo{}, 0, nil
	}
	pageMax := (int(count) + 19) / 20

	if pageNum > pageMax {
		return nil, 0, fmt.Errorf("failed to set page offset :%w", model.ErrTooLargePageNum)
	}

	offset := (pageNum - 1) * 20

	err = query.
		Limit(20).
		Offset(offset).
		Group("questionnaires.id").
		Select("questionnaires.*, (targets.user_traqid = ? OR targets.user_traqid = 'traP') AS is_targeted", userID).
		Find(&questionnaireInfoes).Error

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, 0, model.ErrDeadlineExceeded
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaireInfoes, pageMax, nil
}

func (q *Questionnaire) GetAdminQuestionnaires(ctx context.Context, userID string) ([]model.Questionnaires, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction :%w", err)
	}

	questionnaires := make([]model.Questionnaires, 0)
	err = db.
		Table("questionnaires").
		Joins("INNER JOIN administrators ON questionnaires.id = administrators.questionnaire_id").
		Where("administrators.user_traqid = ?", userID).
		Order("questionnaires.modified_at DESC").
		Find(&questionnaires).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get a questionnaire: %w", err)
	}

	return questionnaires, nil
}

func (q *Questionnaire) GetQuestionnaireInfo(ctx context.Context, questionnaireID int) (*model.Questionnaires, []string, []string, []string, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get transaction:%w", err)
	}
	questionnaire := model.Questionnaires{}
	targets := []string{}
	administrators := []string{}
	respondents := []string{}

	err = db.
		Where("questionnaires.id = ?", questionnaireID).
		First(&questionnaire).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get questionnaire:%w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("targets").
		Where("questionnaire_id =?", questionnaire.ID).
		Pluck("user_traqid", &targets).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get targets : %w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("administrators").
		Where("questionnaire_id = ?", questionnaireID).
		Pluck("user_traqid", &administrators).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get administrators:%w", err)
	}

	err = db.
		Session(&gorm.Session{NewDB: true}).
		Table("respondents").
		Where("questionnaire_id = ? AND deleted_at IS NULL AND submitted_at IS NOT NULL", questionnaire.ID).
		Pluck("user_traqid", &respondents).Error
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get respondent :%w", err)
	}

	return &questionnaire, targets, administrators, respondents, nil
}

func (q *Questionnaire) GetTargetedQuestionnaires(ctx context.Context, userID string, answered string, sort string) ([]model.TargetedQuestionnaire, error) {
	db, err := GetTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction:%w", err)
	}

	query := db.
		Table("questionnaires").
		Where("questionnaires.res_time_limit > ? OR questionnaires.res_time_limit IS NULL", time.Now()).
		Joins("INNER JOIN targets ON questionnaires.id = targets.questionnaire_id").
		Where("targets.user_traqid = ? OR targets.user_traqid = 'traP'", userID).
		Joins("LEFT OUTER JOIN respondents ON questionnaires.id = respondents.questionnaire_id AND respondents.user_traqid = ? AND respondents.deleted_at IS NULL", userID).
		Group("questionnaires.id,respondents.user_traqid").
		Select("questionnaires.*, MAX(respondents.submitted_at) AS responded_at, COUNT(respondents.response_id) != 0 AS has_response")

	query, err = setQuestionnairesOrder(query, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to set the order of the questionnaire table: %w", err)
	}

	query = query.
		Order("questionnaires.res_time_limit").
		Order("questionnaires.modified_at desc")

	switch answered {
	case "answered":
		query = query.Where("respondents.questionnaire_id IS NOT NULL")
	case "unanswered":
		query = query.Where("respondents.questionnaire_id IS NULL")
	case "":
	default:
		return nil, fmt.Errorf("invalid answered parameter value(%s): %w", answered, model.ErrInvalidAnsweredParam)
	}

	questionnaires := []model.TargetedQuestionnaire{}
	err = query.Find(&questionnaires).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get the targeted questionnaires: %w", err)
	}

	return questionnaires, nil
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

func setQuestionnairesOrder(query *gorm.DB, sort string) (*gorm.DB, error) {
	switch sort {
	case "created_at":
		query = query.Order("questionnaires.created_at")
	case "-created_at":
		query = query.Order("questionnaires.created_at desc")
	case "title":
		query = query.Order("questionnaires.title")
	case "-title":
		query = query.Order("questionnaires.title desc")
	case "modified_at":
		query = query.Order("questionnaires.modified_at")
	case "-modified_at":
		query = query.Order("questionnaires.modified_at desc")
	case "":
	default:
		return nil, model.ErrInvalidSortParam
	}
	query = query.Order("questionnaires.id desc")

	return query, nil
}
