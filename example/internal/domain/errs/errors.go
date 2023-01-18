package errs

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"text/template"

	"github.com/lib/pq"
	"go.uber.org/zap/zapcore"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorCode uint

type Params map[string]string

func (p Params) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for key, value := range p {
		encoder.AddString(key, value)
	}
	return nil
}

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
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Params  Params    `json:"params"`
}

func (e *Error) WithParam(key, value string) *Error {
	e.AddParam(key, value)
	return e
}

func (e Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *Error) Is(tgt error) bool {
	target, ok := tgt.(*Error)
	if !ok {
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

func NewInvalidFormError() *Error {
	return NewError(ErrorCodeInvalidArgument, "The form sent is not valid, please correct the errors below.")
}

func NewInvalidParameter(message string) *Error {
	e := NewError(ErrorCodeInvalidArgument, message)
	return e
}

func FromValidationError(err error) *Error {
	var validationErrors validation.Errors
	var validationErrorObject validation.ErrorObject
	if errors.As(err, &validationErrors) {
		e := NewError(ErrorCodeInvalidArgument, "The form sent is not valid, please correct the errors below.")
		for key, value := range validationErrors {
			switch t := value.(type) {
			case validation.ErrorObject:
				e.AddParam(key, renderErrorMessage(t))
			case *Error:
				e.AddParam(key, t.Message)
			default:
				e.AddParam(key, value.Error())
			}
		}
		return e
	}
	if errors.As(err, &validationErrorObject) {
		return NewInvalidParameter(renderErrorMessage(validationErrorObject))
	}
	return nil
}

func renderErrorMessage(object validation.ErrorObject) string {
	parse, err := template.New("message").Parse(object.Message())
	if err != nil {
		return ""
	}
	var tpl bytes.Buffer
	_ = parse.Execute(&tpl, object.Params())
	return tpl.String()
}

func FromPostgresError(err error) *Error {
	e := &Error{
		Code:    ErrorCodeInternal,
		Message: "Unexpected behavior.",
		Params: map[string]string{
			"error": err.Error(),
		},
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		e.AddParam("details", pqErr.Detail)
		e.AddParam("message", pqErr.Message)
		e.AddParam("postgres_code", fmt.Sprint(pqErr.Code))
	}
	if errors.Is(err, sql.ErrNoRows) {
		e = NewEntityNotFound()
	}
	return e
}

func NewEntityNotFound() *Error {
	return &Error{
		Code:    ErrorCodeNotFound,
		Message: "Entity not found.",
		Params:  map[string]string{},
	}
}

func NewPermissionDeniedError() *Error {
	return &Error{
		Code:    ErrorCodePermissionDenied,
		Message: "Permission denied.",
		Params:  map[string]string{},
	}
}

func NewBadToken() *Error {
	return &Error{
		Code:    ErrorCodeUnauthenticated,
		Message: "Bad token.",
		Params:  map[string]string{},
	}
}

func NewPermissionDenied() *Error {
	return &Error{
		Code:    ErrorCodePermissionDenied,
		Message: "Permission denied.",
		Params:  map[string]string{},
	}
}
