package usecase

import (
	"context"
	"github.com/xxarupkaxx/anke-two/domain/model"
	"github.com/xxarupkaxx/anke-two/repository"
	"github.com/xxarupkaxx/anke-two/usecase/input"
	"github.com/xxarupkaxx/anke-two/usecase/output"
	"time"
)

type User struct {
	repository.IRespondent
	repository.IQuestionnaire
	repository.ITarget
	repository.IAdministrator
}

func NewUser(IRespondent repository.IRespondent, IQuestionnaire repository.IQuestionnaire, ITarget repository.ITarget, IAdministrator repository.IAdministrator) *User {
	return &User{IRespondent: IRespondent, IQuestionnaire: IQuestionnaire, ITarget: ITarget, IAdministrator: IAdministrator}
}

func (u *User) GetUsersMe(ctx context.Context, me input.GetMe) output.GetMe {
	opUser := output.GetMe{TraqID: me.UserID}
	return opUser
}

func (u *User) GetMyResponses(ctx context.Context, me input.GetMe) ([]model.RespondentInfo, error) {
	myResponses, err := u.IRespondent.GetRespondentInfos(ctx, me.UserID)
	if err != nil {
		return nil, err
	}

	return myResponses, nil
}

func (u *User) GetMyResponsesByID(ctx context.Context, response input.GetMyResponse) ([]model.RespondentInfo, error) {
	myResponses, err := u.GetRespondentInfos(ctx, response.UserID, response.QuestionnaireID)
	if err != nil {
		return nil, err
	}

	return myResponses, nil
}

func (u *User) GetTargetedQuestionnaire(ctx context.Context, request input.GetTargetedQuestionnaire) ([]model.TargetedQuestionnaire, error) {
	op, err := u.IQuestionnaire.GetTargetedQuestionnaires(ctx, request.UserID, "", request.Sort)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (u *User) GetMyQuestionnaire(ctx context.Context, me input.GetMe) ([]output.QuestionnaireInfo, error) {
	questionnaires, err := u.IQuestionnaire.GetAdminQuestionnaires(ctx, me.UserID)
	if err != nil {
		return nil, err
	}

	questionnaireIDs := make([]int, 0, len(questionnaires))
	for _, q := range questionnaires {
		questionnaireIDs = append(questionnaireIDs, q.ID)
	}

	targets, err := u.ITarget.GetTargets(ctx, questionnaireIDs)
	if err != nil {
		return nil, err
	}

	targetMap := map[int][]string{}
	for _, t := range targets {
		target, ok := targetMap[t.QuestionnaireID]
		if !ok {
			targetMap[t.QuestionnaireID] = []string{t.UserTraqid}
		} else {
			targetMap[t.QuestionnaireID] = append(target, t.UserTraqid)
		}
	}

	respondents, err := u.IRespondent.GetRespondentsUserIDs(ctx, questionnaireIDs)
	if err != nil {
		return nil, err
	}

	respondentMap := map[int][]string{}
	for _, respondent := range respondents {
		rspdts, ok := respondentMap[respondent.QuestionnaireID]
		if !ok {
			respondentMap[respondent.QuestionnaireID] = []string{respondent.UserTraqid}
		} else {
			respondentMap[respondent.QuestionnaireID] = append(rspdts, respondent.UserTraqid)
		}
	}

	administrators, err := u.IAdministrator.GetAdministrators(ctx, questionnaireIDs)
	if err != nil {
		return nil, err
	}

	administratorMap := map[int][]string{}
	for _, administrator := range administrators {
		admins, ok := administratorMap[administrator.QuestionnaireID]
		if !ok {
			administratorMap[administrator.QuestionnaireID] = []string{administrator.UserTraqid}
		} else {
			administratorMap[administrator.QuestionnaireID] = append(admins, administrator.UserTraqid)
		}
	}

	op := []output.QuestionnaireInfo{}

	for _, q := range questionnaires {
		targets, ok := targetMap[q.ID]
		if !ok {
			targets = []string{}
		}

		administrators, ok := administratorMap[q.ID]
		if !ok {
			administrators = []string{}
		}

		respondents, ok := respondentMap[q.ID]
		if !ok {
			respondents = []string{}
		}

		allresponded := true
		for _, t := range targets {
			found := false
			for _, r := range respondents {
				if t == r {
					found = true
					break
				}
			}
			if !found {
				allresponded = false
				break
			}
		}

		op = append(op, output.QuestionnaireInfo{
			ID:             q.ID,
			Title:          q.Title,
			Description:    q.Description,
			ResTimeLimit:   q.ResTimeLimit,
			CreatedAt:      q.CreatedAt.Format(time.RFC3339),
			ModifiedAt:     q.ModifiedAt.Format(time.RFC3339),
			ResSharedTo:    q.ResSharedTo,
			AllResponded:   allresponded,
			Targets:        targets,
			Administrators: administrators,
			Respondents:    respondents,
		})
	}

	return op, nil
}

func (u *User) GetTargetedQuestionnairesByID(ctx context.Context, qid input.GetTargetsByTraQID) ([]model.TargetedQuestionnaire, error) {
	op, err := u.IQuestionnaire.GetTargetedQuestionnaires(ctx, qid.TraQID, qid.Answered, qid.Sort)
	if err != nil {
		return nil, err
	}

	return op, nil
}

type UsersUsecase interface {
	GetUsersMe(ctx context.Context, me input.GetMe) output.GetMe
	GetMyResponses(ctx context.Context, me input.GetMe) ([]model.RespondentInfo, error)
	GetMyResponsesByID(ctx context.Context, response input.GetMyResponse) ([]model.RespondentInfo, error)
	GetTargetedQuestionnaire(ctx context.Context, request input.GetTargetedQuestionnaire) ([]model.TargetedQuestionnaire, error)
	GetMyQuestionnaire(ctx context.Context, me input.GetMe) ([]output.QuestionnaireInfo, error)
	GetTargetedQuestionnairesByID(ctx context.Context, qid input.GetTargetsByTraQID) ([]model.TargetedQuestionnaire, error)
}
