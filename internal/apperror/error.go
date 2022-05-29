package apperror

import (
	"encoding/json"
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	Code             string `json:"code"`
}

var (
	ErrNotFound = NewAppError(nil, "Not found", "", "")
)

// TODO - todo note

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developer_message, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developer_message,
		Code:             code,
	}
}


func systemError(err error) *AppError {
	return NewAppError(err, "unknown system error", err.Error(), "500")
}
