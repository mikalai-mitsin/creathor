// nolint: deadcode
package errs

import (
	"encoding/json"
	"errors"
	"path"
	"reflect"
)

type ErrorCode uint

const (
	ErrorCodeOK ErrorCode = iota
	ErrorCodeCanceled
	ErrorCodeUnknown
	ErrorCodeInvalidArgument
	ErrorCodeDeadlineExceeded
	ErrorCodeNotFound
	ErrorCodeAlreadyExists
	ErrorCodePermissionDenied
	ErrorCodeResourceExhausted
	ErrorCodeFailedPrecondition
	ErrorCodeAborted
	ErrorCodeOutOfRange
	ErrorCodeUnimplemented
	ErrorCodeInternal
	ErrorCodeUnavailable
	ErrorCodeDataLoss
	ErrorCodeUnauthenticated
)

type Error struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Params  map[string]string `json:"params"`
}

func (e Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *Error) Is(tgt error) bool {
	var target *Error
	if !errors.As(tgt, &target) {
		return false
	}
	return reflect.DeepEqual(e, target)
}

func (e *Error) SetCode(code ErrorCode) {
	e.Code = code
}

func (e *Error) SetMessage(message string) {
	e.Message = message
}

func (e *Error) SetParams(params map[string]string) {
	e.Params = params
}

func (e *Error) AddParam(key string, value string) {
	e.Params[key] = value
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Params:  map[string]string{},
	}
}

func NewUnexpectedBehaviorError(details string) *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "Unexpected behavior.",
		Params: map[string]string{
			"details": details,
		},
	}
}

func NewDirectoryNotExistsError(filePath string) *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "The directory does not exist and it would be better if you create it. Kurwa.",
		Params: map[string]string{
			"path": path.Dir(filePath),
		},
	}
}

func NewUserModelNotExistError() *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "User model not exist. Cant enable auth.",
		Params:  map[string]string{},
	}
}

func NewPermissionError(filePath string) *Error {
	return &Error{
		Code:    ErrorCodeInternal,
		Message: "No permissions. You haven't forgotten anything there?",
		Params: map[string]string{
			"path": path.Dir(filePath),
		},
	}
}

func NewBadTemplateError(details string) *Error {
	return &Error{
		Code:    ErrorCodeFailedPrecondition,
		Message: "Bad template.",
		Params: map[string]string{
			"details": details,
		},
	}
}
