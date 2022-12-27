package errs

func NewUserNotFound() *Error {
	return &Error{
		Code:    ErrorCodeNotFound,
		Message: "User not found.",
		Params:  map[string]string{},
	}
}
