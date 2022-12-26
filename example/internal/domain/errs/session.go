package errs

func NewSessionNotFound() *Error {
	return &Error{
		Code:    ErrorCodeNotFound,
		Message: "Session not found.",
		Params:  map[string]string{},
	}
}
