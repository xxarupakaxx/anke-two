package input

import "github.com/xxarupkaxx/anke-two/domain/model"

// Responses 質問に対する回答一覧の構造体
type Responses struct {
	UserID      string
	ID          int                  `json:"questionnaireID" validate:"min=0"`
	Temporarily bool                 `json:"temporarily"`
	Body        []model.ResponseBody `json:"body" validate:"required,dive"`
}

type GetResponse struct {
	ResponseID int
}

type EditResponse struct {
	ResponseID  int
	ID          int                  `json:"questionnaireID" validate:"min=0"`
	Temporarily bool                 `json:"temporarily"`
	Body        []model.ResponseBody `json:"body" validate:"required,dive"`
}

type DeleteResponse struct {
	ResponseID int
}
