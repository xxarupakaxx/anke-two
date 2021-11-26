package output

import (
	"github.com/xxarupkaxx/anke-two/domain/model"
	"time"
)

type PostResponse struct {
	ResponseID      int                  `json:"responseID"`
	QuestionnaireID int                  `json:"questionnaireID"`
	Temporarily     bool                 `json:"temporarily"`
	SubmittedAt     time.Time            `json:"submitted_at"`
	Body            []model.ResponseBody `json:"body"`
}
