package logger

import "github.com/xxarupkaxx/anke-two/domain/logger"

type errorManager struct {

}

func (e *errorManager) Wrap(err error, code int) error {
	panic("implement me")
}

func (e *errorManager) LogMessage(err error) string {
	panic("implement me")
}

func (e *errorManager) Code(err error) int {
	panic("implement me")
}

func NewErrorManager() logger.ErrorManager {
	return &errorManager{}
}
