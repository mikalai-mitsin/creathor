package errs

func NewApproachNotFound() *Error {
	return &Error{
		Code:    ErrorCodeNotFound,
		Message: "Approach not found.",
		Params:  map[string]string{},
	}
}
