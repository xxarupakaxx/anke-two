package input

import (
	"context"
	"gopkg.in/guregu/null.v4"
)

type Questionnaires struct {
	context         context.Context
	title           string
	description     string
	resTimeLimit    null.Time
	resSharedTo     int
	questionnaireID int
	userID          string
	sort            string
	search          string
	pageNum         int
	nonTargeted     bool
	answered        string
	responseID      int
}
