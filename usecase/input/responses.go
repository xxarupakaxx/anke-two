package input

import "github.com/xxarupkaxx/anke-two/domain/model"

// Responses 質問に対する回答一覧の構造体
type Responses struct {
	ID          int                  `json:"questionnaireID" validate:"min=0"`
	Temporarily bool                 `json:"temporarily"`
	Body        []model.ResponseBody `json:"body" validate:"required,dive"`
}
