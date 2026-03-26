package datamark

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// статусы позиции документа
	ItemStatusNotFound      int = -1
	ItemStatusNotRegistered int = 0
	ItemStatusRegistered    int = 1
	ItemStatusReceivedCodes int = 10
)

var (
	ErrIncorrectURL           = errors.New("некорректный URL")
	ErrConnectionTimeout      = errors.New("таймаут соединения/запроса")
	ErrIncorrectRequestMethod = errors.New("неподдерживаемый метод запроса")
)

type ErrorResponse struct {
	ErrorCode int      `json:"error"`
	Message   string   `json:"message"`
	Details   []string `json:"details"`
}

// ErrorDescription возвращает описание ошибки из тела ответа.
func (e *ErrorResponse) ErrorDescription() string {
	description := e.Message

	if len(e.Details) > 0 {
		description = fmt.Sprintf("%s, %s", e.Message, strings.Join(e.Details, ", "))
	}

	return description
}
