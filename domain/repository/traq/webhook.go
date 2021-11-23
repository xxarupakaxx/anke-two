package traq

import "gopkg.in/guregu/null.v4"

// IWebhook traQのWebhookのinterface
type IWebhook interface {
	PostMessage(message string) error
	CreateQuestionnaireMessage(questionnaireID int, title string, description string, administrators []string, resTimeLimit null.Time, targets []string) string
}
